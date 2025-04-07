package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/momokii/go-wa-notifier/internal/models"
	"github.com/momokii/go-wa-notifier/pkg/newsapi"
	"github.com/momokii/go-wa-notifier/pkg/utils"
	"github.com/momokii/go-wa-notifier/pkg/whatsapp"
)

type newsHandler struct {
	newsapi_api_key string
}

func NewNewsHandler(newsapi_api_key string) (*newsHandler, error) {

	if newsapi_api_key == "" {
		return nil, fmt.Errorf("newsapi_api_key is required")
	}

	return &newsHandler{
		newsapi_api_key: newsapi_api_key,
	}, nil
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
//	@Router			/news/wa [post]
func (h *newsHandler) SendNewsAPIWhatsapp(c *fiber.Ctx) error {

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

	message_whatsapp += "Powered by NewsAPI | Kelana Chandra Helyandika | kelanach.xyz"

	// for security, initiate whatsapp just every need to send messages
	waClient, err := whatsapp.NewWhatsApp()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed Initiate Whatsapp: "+err.Error())
	}
	defer func() {
		if err := waClient.Disconnect(); err != nil {
			log.Println("Error disconnecting WhatsApp client:", err)
		}
		log.Println("WhatsApp client disconnected")
	}()

	// send the messages to all the number
	for _, number := range req_body.WhatsappNumbers {
		if err := waClient.SendMessage(number, message_whatsapp, false); err != nil {
			log.Println("Error sending message on number " + number + " error: " + err.Error())
		} else {
			log.Println("Message sent successfully")
		}
	}

	return utils.ResponseMessage(c, fiber.StatusOK, "News sent to WhatsApp successfully")
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
func (h *newsHandler) WhatsAppLogout(c *fiber.Ctx) error {

	waClient, err := whatsapp.NewWhatsApp()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed Initiate Whatsapp: "+err.Error())
	}

	if err := waClient.Logout(); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed Logout Whatsapp: "+err.Error())
	}

	return utils.ResponseMessage(c, fiber.StatusOK, "Logout success")
}

func (h *newsHandler) WAStatus(c *fiber.Ctx) error {
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
