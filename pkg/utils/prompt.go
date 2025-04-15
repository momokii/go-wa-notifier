package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/momokii/go-wa-notifier/pkg/openweatherapi"
)

type NewsType string

const (
	NewsTypeBusiness   NewsType = "business"
	NewsTypeTechnology NewsType = "technology"
	NewsTypeScience    NewsType = "science"
	NewsTypeGeneral    NewsType = "general"
)

func GetNewsType(news_type string) (NewsType, error) {
	switch strings.ToLower(news_type) {
	case "business":
		return NewsTypeBusiness, nil
	case "technology":
		return NewsTypeTechnology, nil
	case "science":
		return NewsTypeScience, nil
	case "general":
		return NewsTypeGeneral, nil
	default:
		return "", fmt.Errorf("invalid news type: %s", news_type)
	}
}

func GenerateNewsSummariesPrompt(news_data string, news_type NewsType) (string, error) {

	var promptSummaries, upperCaseNewsType, domain_specific_analytical_task, domain_specific_section string

	// Check if the news_type is valid
	switch news_type {
	case NewsTypeBusiness, NewsTypeTechnology, NewsTypeScience, NewsTypeGeneral:

		if NewsTypeBusiness == news_type {
			domain_specific_analytical_task = `
			4. Identify economic indicators or market signals
			5. Note corporate developments or policy changes affecting markets
			6. Analyze sector-specific performance or challenges
			`
			domain_specific_section = `
			* Market Implications
			* Sectors to Watch
			* Economic Indicators
			`

		} else if NewsTypeTechnology == news_type {
			domain_specific_analytical_task = `
			4. Identify emerging technologies or innovation trends
			5. Analyze competitive dynamics between tech companies or platforms
			6. Examine regulatory developments affecting technology
			`
			domain_specific_section = `
			- Innovation Highlights
			- Tech Industry Dynamics
			- Digital Transformation Impact
			`

		} else if NewsTypeScience == news_type {
			domain_specific_analytical_task = `
			4. Evaluate the significance of research breakthroughs
			5. Analyze potential applications of scientific developments
			6. Identify interdisciplinary implications
			`
			domain_specific_section = `
			- Research Breakthroughs
			- Practical Applications
			- Scientific Community Developments
			`

		} else if NewsTypeGeneral == news_type {
			domain_specific_analytical_task = `
			4. Identify cross-domain patterns or interconnections
			5. Highlight societal impacts across different sectors
			6. Note emerging broad trends affecting multiple areas
			`
			domain_specific_section = `
			- Market Implications
			- Sectors to Watch
			- Economic Indicators
			`

		}

	default:
		return "", fmt.Errorf("invalid news type: %s", news_type)
	}

	upperCaseNewsType = strings.ToUpper(string(news_type))

	promptSummaries = fmt.Sprintf(`You are an expert analyst specializing in %s.

	I'll provide you with a set of recent %s news headlines and summaries. Your task is to:
	1. Analyze these news items and identify key patterns or trends
	2. Extract actionable insights relevant to %s
	3. Highlight potential impacts for stakeholders in this field
	%s 

	Present your analysis in a clear format under the heading "DAILY %s INSIGHTS" with the following sections:
	* Key Trends Identified
	%s
	* Strategic Considerations

	Here are the news items to analyze:
	%s

	End your analysis with 2-3 key takeaways that summarize the most important insights from today's %s news.
	
	IMPORTANT FORMATTING INSTRUCTIONS:
	- Use WhatsApp formatting standards throughout your response
	- For headers and section titles, use *asterisks for bold text*
	- For emphasis within paragraphs, use _underscores for italic text_
	- For lists, use proper bullet points (•) or numbers followed by periods
	- For critical insights or statistics, use both *bold* and _italic_ formatting where appropriate
	- When referencing specific news items, use the same formatting style as seen in the provided examples
	- Make sure all key points and takeaways are formatted in *bold* for easy visibility
	- Format the final takeaways section as "*Key Takeaways:*" followed by numbered points
	`,
		upperCaseNewsType, upperCaseNewsType, upperCaseNewsType, domain_specific_analytical_task,
		upperCaseNewsType, domain_specific_section, news_data, upperCaseNewsType)

	return promptSummaries, nil
}

// generateWeatherPrompt creates a comprehensive prompt for OpenAI to generate weather reports
func GenerateWeatherPrompt(data map[string]interface{}) string {
	reportType := data["reportType"].(string)
	timeContext := "today"
	if reportType == "tomorrow" {
		timeContext = "tomorrow"
	}

	prompt := fmt.Sprintf(`
You are a professional weather forecaster providing accurate and useful weather reports for WhatsApp users.

## DATA CONTEXT
I will provide you with three types of weather data for coordinates [%.4f, %.4f]:
1. Overview summary
2. Daily aggregate statistics
3. Hour-by-hour forecast for the next 24 hours

Your task is to analyze this data and create a concise, informative, and visually engaging WhatsApp message for %s's weather (%s).

## LOCATION CONTEXT
First, determine the location name based on these coordinates: Latitude %.4f, Longitude %.4f
For example: "Jakarta, Indonesia" or "South Jakarta, Indonesia" - be as specific as possible.

## WEATHER DATA
1. Weather Overview: %s
2. Daily Aggregate:
	- Temperature: Min %.1f°C, Max %.1f°C
	- Morning: %.1f°C, Afternoon: %.1f°C, Evening: %.1f°C, Night: %.1f°C
	- Humidity (afternoon): %.0f%%
	- Cloud Cover (afternoon): %.0f%%
	- Precipitation Total: %.1fmm
	- Wind Speed (max): %.1f m/s, Direction: %.0f°
	- Pressure (afternoon): %.0f hPa

## HOUR-BY-HOUR DATA
%s

## OUTPUT FORMAT
Create a WhatsApp-ready message using emojis and formatting with the following sections (must follow and have these sections):
1. HEADER: Create an eye-catching title with location and date
2. OVERVIEW: A 2-3 sentence summary of the day's weather
3. KEY METRICS: Important temperature, precipitation, and wind data
4. HOURLY HIGHLIGHTS: Key weather changes throughout the day (morning, afternoon, evening, night)
5. RECOMMENDATIONS: 3-5 practical suggestions based on the forecast (what to wear, activities to consider/avoid, precautions)
6. Key Takeaways: 2-3 concise points summarizing the most important insights from the weather report
7. Quote: some inspirational quote related to weather or nature that matches the forecast

Use appropriate weather emojis (☀️🌤️⛅🌥️☁️🌧️⛈️❄️) to make the message visually engaging.
Keep your response concise (under 1000 characters) and optimized for mobile viewing.
Format temperatures in Celsius with the degree symbol (°C)

IMPORTANT FORMATTING INSTRUCTIONS:
- Use WhatsApp formatting standards throughout your response
- For headers and section titles, use *asterisks for bold text*
- For emphasis within paragraphs, use _underscores for italic text_
- For lists, use proper bullet points (•) or numbers followed by periods
- For critical insights or statistics, use both *bold* and _italic_ formatting where appropriate
- When referencing specific news items, use the same formatting style as seen in the provided examples
- Make sure all key points and takeaways are formatted in *bold* for easy visibility
- Format the final takeaways section as "*Key Takeaways:*" followed by numbered points
`,
		data["latitude"].(float64),
		data["longitude"].(float64),
		timeContext,
		data["date"].(string),
		data["latitude"].(float64),
		data["longitude"].(float64),
		data["weatherOverview"].(string),
		data["dailyAggregate"].(openweatherapi.OpenWeatherAPIV3OneCallDailySummaryResp).Temperature.Min,
		data["dailyAggregate"].(openweatherapi.OpenWeatherAPIV3OneCallDailySummaryResp).Temperature.Max,
		data["dailyAggregate"].(openweatherapi.OpenWeatherAPIV3OneCallDailySummaryResp).Temperature.Morning,
		data["dailyAggregate"].(openweatherapi.OpenWeatherAPIV3OneCallDailySummaryResp).Temperature.Afternoon,
		data["dailyAggregate"].(openweatherapi.OpenWeatherAPIV3OneCallDailySummaryResp).Temperature.Evening,
		data["dailyAggregate"].(openweatherapi.OpenWeatherAPIV3OneCallDailySummaryResp).Temperature.Night,
		data["dailyAggregate"].(openweatherapi.OpenWeatherAPIV3OneCallDailySummaryResp).Humidity.Afternoon,
		data["dailyAggregate"].(openweatherapi.OpenWeatherAPIV3OneCallDailySummaryResp).CloudCover.Afternoon,
		data["dailyAggregate"].(openweatherapi.OpenWeatherAPIV3OneCallDailySummaryResp).Precipitation.Total,
		data["dailyAggregate"].(openweatherapi.OpenWeatherAPIV3OneCallDailySummaryResp).Wind.Max.Speed,
		data["dailyAggregate"].(openweatherapi.OpenWeatherAPIV3OneCallDailySummaryResp).Wind.Max.Direction,
		data["dailyAggregate"].(openweatherapi.OpenWeatherAPIV3OneCallDailySummaryResp).Pressure.Afternoon,
		formatHourlyDataForPrompt(data["hourlyForecast"].([]openweatherapi.HourlyData)),
	)

	return prompt
}

// formatHourlyDataForPrompt converts hourly data into a readable format for the AI prompt
func formatHourlyDataForPrompt(hourlyData []openweatherapi.HourlyData) string {
	var result strings.Builder

	// Only include selected hours for brevity (every 3 hours)
	selectedHours := []int{0, 3, 6, 9, 12, 15, 18, 21}

	for _, hour := range selectedHours {
		if hour < len(hourlyData) {
			data := hourlyData[hour]
			timeStr := time.Unix(data.Dt, 0).Format("15:04")
			weather := "No weather data"
			if len(data.Weather) > 0 {
				weather = data.Weather[0].Main + " (" + data.Weather[0].Description + ")"
			}

			result.WriteString(fmt.Sprintf("- %s: %.1f°C, %s, Humidity: %d%%, Wind: %.1f m/s, Precipitation Chance: %.0f%%\n",
				timeStr,
				data.Temp,
				weather,
				data.Humidity,
				data.WindSpeed,
				data.Pop*100))
		}
	}

	return result.String()
}

// formatWeatherMessage adds appropriate headers and footers to the weather report
func FormatWeatherMessage(content string, data map[string]interface{}) string {
	reportType := "TODAY'S"
	if data["reportType"].(string) == "tomorrow" {
		reportType = "TOMORROW'S"
	}

	header := fmt.Sprintf("🌤️ *%s WEATHER FORECAST* 🌤️\n\n", reportType)
	footer := "\n\nPowered by OpenWeather | Kelana Chandra Helyandika | kelanach.xyz"

	return header + content + footer
}
