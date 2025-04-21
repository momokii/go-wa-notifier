package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/momokii/go-llmbridge/pkg/openai"
	"github.com/momokii/go-wa-notifier/internal/models"
	"github.com/momokii/go-wa-notifier/pkg/newsapi"
	"github.com/momokii/go-wa-notifier/pkg/openweatherapi"
	"github.com/momokii/go-wa-notifier/pkg/utils"
	"github.com/momokii/go-wa-notifier/pkg/whatsapp"
)

// ================ WHATSAPP HANDLER APPENDIX FUNCTION/DATA TYPE/CONST

// for swagger docs
type WAStatusResponse struct {
	Error   bool   `json:"error" example:"false"`
	Message string `json:"message"`
	Data    struct {
		IsConnected bool   `json:"is_connected"`
		IsReady     bool   `json:"is_ready"`
		QRCode      string `json:"qr_code"`
	} `json:"data"`
}

func whatsappSendMessages(messages string, numbers []string) error {
	// setup whatsapp client and not disconnect it when done
	waClient, err := whatsapp.NewWhatsApp()
	if err != nil {
		return fmt.Errorf("failed to initiate WhatsApp: %w", err)
	}

	if !waClient.IsConnected() {
		return fmt.Errorf("WhatsApp client is not connected")
	}

	// send messages to all numbers
	for _, number := range numbers {
		if err := waClient.SendMessage(number, messages, false); err != nil {
			log.Println("Error sending message on number " + number + " error: " + err.Error())
			continue
		}
		log.Println("Message sent successfully")
	}

	return nil
}

// ================ MAIN HANDLER

type whatsappHandler struct {
	newsapi_api_key     string
	openweather_api_key string
	openaiClient        openai.OpenAI
}

func NewWhatsappHandler(
	newsapi_api_key string,
	openweather_api_key string,
	openaiClient openai.OpenAI,
) (*whatsappHandler, error) {

	if newsapi_api_key == "" {
		return nil, fmt.Errorf("newsapi_api_key is required")
	}

	return &whatsappHandler{
		newsapi_api_key:     newsapi_api_key,
		openweather_api_key: openweather_api_key,
		openaiClient:        openaiClient,
	}, nil
}

// SendMessages godoc
//
//	@Summary		Send messages custom to whatsapp
//	@Description	Send messages custom to whatsapp
//	@Tags			News
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.WhatsappMessagesReq	true	"body request detail"
//	@Success		200		{object}	utils.MessageResponseSuccess
//	@Failure		400		{object}	utils.MessageResponseError
//	@Failure		500		{object}	utils.MessageResponseError
//	@Router			/wa/messages [post]
func (h *whatsappHandler) SendMessages(c *fiber.Ctx) error {

	// for now, set to max whatsapp send number to 100
	req_body := new(models.WhatsappMessagesReq)
	if err := c.BodyParser(req_body); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req_body.Messages == "" {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Messages is required")
	}

	if len(req_body.WhatsappNumbers) == 0 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Whatsapp numbers is required")
	}

	// check that max numbers send is 100
	if len(req_body.WhatsappNumbers) > 100 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Max Whatsapp numbers is 100")
	}

	// send messages to all numbers
	if err := whatsappSendMessages(req_body.Messages, req_body.WhatsappNumbers); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to send messages: "+err.Error())
	}

	return utils.ResponseMessage(c, fiber.StatusOK, "Send Messages to Whatsapp")
}

// SendNewsAPIWhatsapp godoc
//
//	@Summary		Send news to whatsapp
//	@Description	Send news to whatsapp
//	@Tags			News
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.NewsSendWhatsappReq	true	"body request detail"
//	@Success		200		{object}	utils.MessageResponseSuccess
//	@Failure		400		{object}	utils.MessageResponseError
//	@Failure		500		{object}	utils.MessageResponseError
//	@Router			/wa/news [post]
func (h *whatsappHandler) SendNewsAPIWhatsapp(c *fiber.Ctx) error {

	req_body := new(models.NewsSendWhatsappReq)
	if err := c.BodyParser(req_body); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if len(req_body.WhatsappNumbers) == 0 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Whatsapp numbers is required")
	}

	if req_body.Category == "" {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Category is required")
	}

	// create req struct for newsapi
	// this will be adjust to my need for whastapp notifier that one req will be max get 10 top headlines for better expererience
	query_newsapi := newsapi.NewsAPITopHeadlinesReq{
		PageSize: 10,
		Page:     1,
		Category: req_body.Category,
	}

	// call newsapi to get the news
	news_resp, err := newsapi.NewsAPITopHeadlines(h.newsapi_api_key, query_newsapi)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Error Get News Data: "+err.Error())
	}

	// check if newsapi response is error from the api or not
	if news_resp.Status != "ok" {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to get news from newsapi: "+news_resp.Message)
	}

	// api call success, process the articles
	category_uppercase := strings.ToUpper(req_body.Category)
	message_whatsapp := fmt.Sprintf("📰 *TOP %s NEWS TODAY* 📰\n\n", category_uppercase)

	var news_data string
	for i, article := range news_resp.Articles {
		message_whatsapp += fmt.Sprintf("*%d. %s*\n", i+1, article.Title)
		message_whatsapp += fmt.Sprintf("📄 *Source:* %s\n", article.Source.Name)

		if article.Author != "" {
			message_whatsapp += fmt.Sprintf("✍️ *Author:* %s\n", article.Author)
		}

		if article.PublishedAt != "" {
			// Parse the ISO 8601 date format
			// example of article.PublishedAt: "2025-04-04T14:19:00Z"
			t, err := time.Parse(time.RFC3339, article.PublishedAt)
			if err == nil {
				// Format as a more readable date: e.g., "04 Apr 2025, 14:19"
				formattedDate := t.Format("02 Jan 2006, 15:04")
				message_whatsapp += fmt.Sprintf("📅 *Published:* %s\n", formattedDate)
			} else {
				// Fallback to original format if parsing fails
				message_whatsapp += fmt.Sprintf("📅 *Published:* %s\n", article.PublishedAt)
			}
		}

		if article.Description != "" {
			message_whatsapp += fmt.Sprintf("📝 *Summary:* %s\n", article.Description)
		}

		message_whatsapp += fmt.Sprintf("🔗 *Read more:* %s\n\n", article.Url)
	}

	// continue using llm if using_llm is true
	if req_body.UsingLLM {
		news_data = message_whatsapp

		news_type, err := utils.GetNewsType(req_body.Category)
		if err != nil {
			return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid category: "+err.Error())
		}

		prompt_news_summaries, err := utils.GenerateNewsSummariesPrompt(news_data, news_type)
		if err != nil {
			return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to generate news summaries: "+err.Error())
		}

		// create message and send to openai for summarization
		message_summaries := []openai.OAMessageReq{
			{
				Role:    "user",
				Content: prompt_news_summaries,
			},
		}

		summaries_news_resp, err := h.openaiClient.OpenAIGetFirstContentDataResp(&message_summaries, false, nil, false, nil)
		if err != nil {
			return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to get summaries from OpenAI: "+err.Error())
		}

		// add the summaries to the message
		message_whatsapp += fmt.Sprintf("🤖 *AI Summaries:*\n%s\n\n", summaries_news_resp.Content)
	}

	// add footer to the message
	message_whatsapp += "Powered by NewsAPI | Kelana Chandra Helyandika | kelanach.xyz"

	// send messages to all numbers
	if err := whatsappSendMessages(message_whatsapp, req_body.WhatsappNumbers); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to send messages: "+err.Error())
	}

	return utils.ResponseMessage(c, fiber.StatusOK, "News sent to WhatsApp successfully")
}

// SendWeatherAPIWhatsapp godoc
//
//	@Summary		Send weather daily forecast to whatsapp
//	@Description	Send weather daily forecast
//	@Tags			News
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.WeatherSendWhatsappReq	true	"body request detail"
//	@Success		200		{object}	utils.MessageResponseSuccess
//	@Failure		400		{object}	utils.MessageResponseError
//	@Failure		500		{object}	utils.MessageResponseError
//	@Router			/wa/weathers [post]
func (h *whatsappHandler) SendWeatherAPIWhatsapp(c *fiber.Ctx) error {

	// parser body request and check the validity
	req_body := new(models.WeatherSendWhatsappReq)
	if err := c.BodyParser(req_body); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if len(req_body.WhatsappNumbers) == 0 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Whatsapp numbers is required")
	}

	if req_body.Type == "" && req_body.Type != "today" && req_body.Type != "tomorrow" {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Type is required and must be 'today' or 'tomorrow'")
	}

	if req_body.Lat < -90 || req_body.Lat > 90 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Latitude must be between -90 and 90")
	}

	if req_body.Lon < -180 || req_body.Lon > 180 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Longitude must be between -180 and 180")
	}

	// process the weather data here

	// get date, based on today or tomorrow
	var date, reportType, localTime string
	if req_body.Type == "today" {
		date = time.Now().Format("2006-01-02")
		localTime = time.Now().Format("15:04:05")
		reportType = "today"
	} else {
		date = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		localTime = time.Now().AddDate(0, 0, 1).Format("15:04:05")
		reportType = "tomorrow"
	}

	weather_base_req := openweatherapi.OpenWeatherAPIV3OneCallBaseReq{
		Lat:   req_body.Lat,
		Lon:   req_body.Lon,
		AppID: h.openweather_api_key,
		Units: "metric", // using celcius as default for this endpoint
	}

	// first, get OVERVIEW DATA
	weather_overview_req := openweatherapi.OpenWeatherAPIV3OneCallOverviewReq{
		Date:                           date,
		OpenWeatherAPIV3OneCallBaseReq: weather_base_req,
	}

	weather_overview, err := openweatherapi.OpenWeatherV3OneCallOverviewAPI(weather_overview_req)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to get weather overview: "+err.Error())
	}

	// second, get DAILY AGGREATE DATA
	weather_daily_req := openweatherapi.OpenWeatherAPIV3OneCallDailySummaryReq{
		Date:                           date,
		OpenWeatherAPIV3OneCallBaseReq: weather_base_req,
	}

	weather_daily_aggregate, err := openweatherapi.OpenWeatherV3OneCallDailySummaryAPI(weather_daily_req)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to get weather daily aggregate: "+err.Error())
	}

	// third, get the weather data for hourly to get 24 hours data with onecall basic api
	weather_onecall_24hr := openweatherapi.OpenWeatherAPIV3OneCallReq{
		OpenWeatherAPIV3OneCallBaseReq: weather_base_req,
		Exclude:                        []string{"current", "minutely", "daily", "alerts"}, // just get hourly data
	}

	weather_onecall_24hr_resp, err := openweatherapi.OpenWeatherV3OneCallAPI(weather_onecall_24hr)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to get weather onecall 24 hours: "+err.Error())
	}

	// get 24 data hourly from the response
	weather_24hr := weather_onecall_24hr_resp.Hourly[0:24]

	weatherData := openweatherapi.WeatherDataAggregate{
		Date:             date,
		ReportType:       reportType,
		Latitude:         req_body.Lat,
		Longitude:        req_body.Lon,
		WeatherOverview:  weather_overview.WeatherOverview,
		Timezone:         weather_overview.TZ,
		DailyAggregate:   weather_daily_aggregate,
		HourlyForecast:   weather_24hr,
		CurrentTimeLocal: localTime,
	}

	var messages_wa string
	// check if using llm or not, if not just send the weather data to whatsapp
	if req_body.UsingLLM {
		// generate prompt
		prompt := utils.GenerateWeatherPrompt(&weatherData)

		// send to openai for summarization
		messages := []openai.OAMessageReq{
			{
				Role:    "user",
				Content: prompt,
			},
		}

		weather_ai, err := h.openaiClient.OpenAIGetFirstContentDataResp(&messages, false, nil, false, nil)
		if err != nil {
			return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to get weather summary: "+err.Error())
		}

		// format response to message whatsapp
		messages_wa = utils.FormatWeatherMessage(weather_ai.Content, &weatherData)
	} else {
		// If not using LLM, format the weather data manually
		messages_wa = utils.FormatWeatherMessageManual(&weatherData)
	}

	// send messages
	if err := whatsappSendMessages(messages_wa, req_body.WhatsappNumbers); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to send messages: "+err.Error())
	}

	return utils.ResponseMessage(c, fiber.StatusOK, "Send WeatherAPI to Whatsapp")
}

// WhatsAppLogout godoc
//
//	@Summary		Logout Whatsapp Account
//	@Description	Logout Whatsapp Account
//	@Tags			Whatsapp
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	utils.MessageResponseSuccess
//	@Failure		400		{object}	utils.MessageResponseError
//	@Failure		500		{object}	utils.MessageResponseError
//	@Router			/wa/logout [post]
func (h *whatsappHandler) WhatsAppLogout(c *fiber.Ctx) error {

	waClient, err := whatsapp.NewWhatsApp()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed Initiate Whatsapp: "+err.Error())
	}

	if err := waClient.Logout(); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed Logout Whatsapp: "+err.Error())
	}

	return utils.ResponseMessage(c, fiber.StatusOK, "Logout success")
}

// WAStatus godoc
//
//	@Summary		Check Whatsapp Status
//	@Description	Check Whatsapp Status
//	@Tags			Whatsapp
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	handlers.WAStatusResponse
//	@Failure		400		{object}	utils.MessageResponseError
//	@Failure		500		{object}	utils.MessageResponseError
//	@Router			/wa/status [get]
func (h *whatsappHandler) WAStatus(c *fiber.Ctx) error {
	waClient, err := whatsapp.NewWhatsApp()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed Initiate Whatsapp: "+err.Error())
	}

	qrCode, qrReady := waClient.GetQRCode()
	isConnected := waClient.IsConnected()

	return utils.ResponseWitData(c, fiber.StatusOK, "WhatsApp Status", fiber.Map{
		"IsConnected": isConnected,
		"IsReady":     qrReady,
		"QRCode":      qrCode,
	})
}
