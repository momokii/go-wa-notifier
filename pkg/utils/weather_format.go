package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/momokii/go-wa-notifier/pkg/openweatherapi"
)

// FormatWeatherMessage adds appropriate headers and footers to the weather report
// using it when using LLM
func FormatWeatherMessage(content string, data *openweatherapi.WeatherDataAggregate) string {
	reportType := "TODAY'S"
	if data.ReportType == "tomorrow" {
		reportType = "TOMORROW'S"
	}

	header := fmt.Sprintf("🌤️ *%s WEATHER FORECAST* 🌤️\n\n", reportType)
	footer := "\n\nPowered by OpenWeather | Kelana Chandra Helyandika | kelanach.xyz"

	return header + content + footer
}

// FormatWeatherMessageManual creates a formatted weather message without using LLM
func FormatWeatherMessageManual(weatherData *openweatherapi.WeatherDataAggregate) string {
	// Build the message
	var message strings.Builder

	// Format header based on whether it's today or tomorrow
	var reportTypeCaps string
	if weatherData.ReportType == "tomorrow" {
		reportTypeCaps = "TOMORROW'S"
	} else {
		reportTypeCaps = "TODAY'S" // Default to today if not specified
	}

	// Add header
	message.WriteString(fmt.Sprintf("🌤️ *%s WEATHER FORECAST* 🌤️\n", reportTypeCaps))
	message.WriteString(fmt.Sprintf("📍 Coordinates: [%.4f, %.4f]\n", weatherData.Latitude, weatherData.Longitude))
	message.WriteString(fmt.Sprintf("📅 Date: %s\n", weatherData.Date))
	message.WriteString(fmt.Sprintf("🌐 Timezone: %s\n\n", weatherData.Timezone))

	// Add overview
	message.WriteString("*📝 OVERVIEW*\n")
	message.WriteString(fmt.Sprintf("%s\n\n", weatherData.WeatherOverview))

	// Add temperature data
	message.WriteString("*🌡️ TEMPERATURE*\n")
	message.WriteString(fmt.Sprintf("• Min: %.1f°C | Max: %.1f°C\n",
		weatherData.DailyAggregate.Temperature.Min,
		weatherData.DailyAggregate.Temperature.Max))
	message.WriteString(fmt.Sprintf("• Morning: %.1f°C | Afternoon: %.1f°C\n",
		weatherData.DailyAggregate.Temperature.Morning,
		weatherData.DailyAggregate.Temperature.Afternoon))
	message.WriteString(fmt.Sprintf("• Evening: %.1f°C | Night: %.1f°C\n\n",
		weatherData.DailyAggregate.Temperature.Evening,
		weatherData.DailyAggregate.Temperature.Night))

	// Add other weather conditions
	message.WriteString("*☁️ CONDITIONS*\n")
	message.WriteString(fmt.Sprintf("• Humidity: %.0f%%\n", weatherData.DailyAggregate.Humidity.Afternoon))
	message.WriteString(fmt.Sprintf("• Cloud Cover: %.0f%%\n", weatherData.DailyAggregate.CloudCover.Afternoon))
	message.WriteString(fmt.Sprintf("• Precipitation: %.1fmm\n", weatherData.DailyAggregate.Precipitation.Total))
	message.WriteString(fmt.Sprintf("• Wind: %.1f m/s at %.0f°\n",
		weatherData.DailyAggregate.Wind.Max.Speed,
		weatherData.DailyAggregate.Wind.Max.Direction))
	message.WriteString(fmt.Sprintf("• Pressure: %.0f hPa\n\n", weatherData.DailyAggregate.Pressure.Afternoon))

	// Add hourly forecast (key times of day)
	message.WriteString("*⏰ KEY HOURS FORECAST*\n")

	// Select key hours based on time of day
	keyHours := []int{6, 12, 18}
	timeLabels := []string{"Morning", "Afternoon", "Evening"}

	// Only process if we have enough hourly data
	if len(weatherData.HourlyForecast) > 0 {
		for i, offsetHour := range keyHours {
			// Calculate the actual array index
			hourIndex := offsetHour
			if hourIndex < len(weatherData.HourlyForecast) {
				data := weatherData.HourlyForecast[hourIndex]

				// Format time in a user-friendly way
				timeObj := time.Unix(data.Dt, 0)
				timeStr := timeObj.Format("15:04")

				// Default values in case data is missing
				weatherDesc := "No data"
				weatherEmoji := "❓"

				// Extract weather info if available
				if len(data.Weather) > 0 {
					weatherDesc = data.Weather[0].Description

					// Select emoji based on weather condition
					switch data.Weather[0].Main {
					case "Clear":
						weatherEmoji = "☀️"
					case "Clouds":
						weatherEmoji = "☁️"
					case "Rain":
						weatherEmoji = "🌧️"
					case "Drizzle":
						weatherEmoji = "🌦️"
					case "Thunderstorm":
						weatherEmoji = "⛈️"
					case "Snow":
						weatherEmoji = "❄️"
					case "Mist", "Fog", "Haze":
						weatherEmoji = "🌫️"
					default:
						weatherEmoji = "🌤️"
					}
				}

				message.WriteString(fmt.Sprintf("• %s (%s): %s %.1f°C, %s, %d%% humidity, %.0f%% chance of rain\n",
					timeLabels[i],
					timeStr,
					weatherEmoji,
					data.Temp,
					weatherDesc,
					data.Humidity,
					data.Pop*100))
			}
		}
	} else {
		message.WriteString("• Hourly forecast data not available\n")
	}

	// Add recommendations based on weather
	message.WriteString("\n*💡 RECOMMENDATIONS*\n")

	// Check for rain
	if weatherData.DailyAggregate.Precipitation.Total > 0 {
		message.WriteString("• Carry an umbrella or raincoat ☔\n")
	}

	// Temperature recommendations
	if weatherData.DailyAggregate.Temperature.Max > 30 {
		message.WriteString("• Stay hydrated and wear light clothing 💧\n")
		message.WriteString("• Use sunscreen if going outdoors 🧴\n")
	} else if weatherData.DailyAggregate.Temperature.Min < 15 {
		message.WriteString("• Wear warm clothing, especially in the morning/evening 🧥\n")
	}

	// Wind recommendations
	if weatherData.DailyAggregate.Wind.Max.Speed > 10 {
		message.WriteString("• Expect strong winds - secure loose items outdoors 💨\n")
	}

	// Add key takeaways
	message.WriteString("\n*🔑 KEY TAKEAWAYS*\n")

	// Generate takeaways based on weather conditions
	// Rain forecast
	if weatherData.DailyAggregate.Precipitation.Total > 5 {
		message.WriteString("• Expect significant rainfall, plan indoor activities ☔\n")
	} else if weatherData.DailyAggregate.Precipitation.Total > 0 {
		message.WriteString("• Light rain possible, keep an umbrella handy 🌂\n")
	} else {
		message.WriteString("• Dry conditions expected, no rain gear needed 👍\n")
	}

	// Temperature takeaway
	tempDiff := weatherData.DailyAggregate.Temperature.Max - weatherData.DailyAggregate.Temperature.Min
	if tempDiff > 10 {
		message.WriteString("• Large temperature swings throughout the day, dress in layers 🧥➡️👕\n")
	}

	// Add a weather-appropriate quote
	message.WriteString("\n*💭 WEATHER WISDOM*\n")
	if weatherData.DailyAggregate.Precipitation.Total > 0 {
		message.WriteString("\"The best thing one can do when it's raining is to let it rain.\" - Henry W. Longfellow\n")
	} else if weatherData.DailyAggregate.CloudCover.Afternoon > 70 {
		message.WriteString("\"Clouds come floating into my life, no longer to carry rain or usher storm, but to add color to my sunset sky.\" - Rabindranath Tagore\n")
	} else {
		message.WriteString("\"Wherever you go, no matter what the weather, always bring your own sunshine.\" - Anthony J. D'Angelo\n")
	}

	// Add footer
	message.WriteString("\n*Weather data provided by OpenWeather*")
	message.WriteString("\nPowered by Kelana Chandra Helyandika | kelanach.xyz")

	return message.String()
}
