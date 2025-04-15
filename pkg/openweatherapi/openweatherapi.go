package openweatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// checker function for base query value that every request must have
func checkBaseQuery(query_req OpenWeatherAPIV3OneCallBaseReq) error {
	if query_req.Lat < -90 || query_req.Lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}

	if query_req.Lon < -180 || query_req.Lon > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}

	if query_req.AppID == "" {
		return fmt.Errorf("appid or API Key is required")
	}

	if query_req.Units != "" && query_req.Units != "metric" && query_req.Units != "imperial" && query_req.Units != "standard" {
		return fmt.Errorf("units must be either metric or imperial or standard")
	}

	return nil
}

// OpenWeatherV3OneCallAPI fetches weather data from the OpenWeather API v3 One Call endpoint.
// This function makes a request to the OpenWeather API using the provided query parameters
// and returns the weather information as a structured response.
//
// Parameters:
//   - query_req: OpenWeatherAPIV3OneCallReq containing request parameters such as:
//   - Lat: Latitude (required)
//   - Lon: Longitude (required)
//   - AppID: API key (required)
//   - Units: Units format (optional, defaults to "standard")
//   - Exclude: Array of sections to exclude from the response (optional)
//     Valid values: "current", "minutely", "hourly", "daily", "alerts"
//   - Lang: Language for output (optional)
//
// Returns:
//   - OpenWeatherAPIV3OneCallResp: A structured response containing weather data
//   - error: Error if the request fails or validation fails
//
// Example Response Structure:
//
//	{
//	  "lat": 33.44,                    // Geographical coordinates (latitude)
//	  "lon": -94.04,                   // Geographical coordinates (longitude)
//	  "timezone": "America/Chicago",   // Local timezone
//	  "timezone_offset": -18000,       // Timezone offset in seconds from UTC
//	  "current": {                     // Current weather data
//	    "dt": 1684929490,              // Current time (Unix timestamp)
//	    "sunrise": 1684926645,         // Sunrise time (Unix timestamp)
//	    "temp": 292.55,                // Temperature
//	    "weather": [{                  // Weather condition details
//	      "id": 803,                   // Weather condition ID
//	      "main": "Clouds",            // Group of weather parameters
//	      "description": "broken clouds", // Weather description
//	      "icon": "04d"                // Weather icon ID
//	    }]
//	    // Additional fields: feels_like, pressure, humidity, etc.
//	  },
//	  "minutely": [                    // Minute forecast (if not excluded)
//	    {
//	      "dt": 1684929540,            // Time of forecasted data (Unix timestamp)
//	      "precipitation": 0           // Precipitation volume (mm)
//	    },
//	    // ... minute by minute forecast
//	  ],
//	  "hourly": [                      // Hourly forecast (if not excluded)
//	    {
//	      "dt": 1684926000,            // Time of forecasted data (Unix timestamp)
//	      "temp": 292.01,              // Temperature
//	      "weather": [{ /* weather details */ }],
//	      // Additional fields: feels_like, pressure, humidity, etc.
//	    },
//	    // ... hour by hour forecast
//	  ],
//	  "daily": [                       // Daily forecast (if not excluded)
//	    {
//	      "dt": 1684951200,            // Time of forecasted data (Unix timestamp)
//	      "summary": "Expect a day of partly cloudy with rain", // Summary
//	      "temp": {                    // Temperature breakdown
//	        "day": 299.03,             // Day temperature
//	        "min": 290.69,             // Min temperature
//	        "max": 300.35              // Max temperature
//	        // Additional fields: night, eve, morn
//	      },
//	      "weather": [{ /* weather details */ }],
//	      // Additional fields: precipitation, pressure, humidity, etc.
//	    },
//	    // ... day by day forecast
//	  ],
//	  "alerts": [                      // Weather alerts (if not excluded)
//	    {
//	      "sender_name": "NWS Philadelphia - Mount Holly",
//	      "event": "Small Craft Advisory",
//	      "start": 1684952747,         // Start time (Unix timestamp)
//	      "end": 1684988747,           // End time (Unix timestamp)
//	      "description": "...SMALL CRAFT ADVISORY REMAINS IN EFFECT..."
//	    },
//	    // ... additional alerts
//	  ]
//	}
//
// Errors:
//   - Returns error if required parameters (Lat, Lon, AppID) are missing
//   - Returns error if any exclude value is invalid
//   - Returns error if API request fails or returns non-OK status
//
// source: https://openweathermap.org/api/one-call-3
func OpenWeatherV3OneCallAPI(query_req OpenWeatherAPIV3OneCallReq) (OpenWeatherAPIV3OneCallResp, error) {

	var resp OpenWeatherAPIV3OneCallResp

	// base value checker
	base_req := OpenWeatherAPIV3OneCallBaseReq{
		Lat:   query_req.Lat,
		Lon:   query_req.Lon,
		AppID: query_req.AppID,
		Units: query_req.Units,
	}

	if err := checkBaseQuery(base_req); err != nil {
		return resp, err
	}

	// check additional param value

	// make sure also that exclude data is all unique
	excludeUnique := make(map[string]bool)
	var excludeData []string
	exclude_query_value := "" // to store the exclude query value
	if len(query_req.Exclude) > 0 {

		if len(query_req.Exclude) > 5 {
			return resp, fmt.Errorf("exclude value must be max 5, because there are 5 data available (current, minutely, hourly, daily and alerts)")
		}

		for _, exclude := range query_req.Exclude {
			if exclude != "current" && exclude != "minutely" && exclude != "hourly" && exclude != "daily" && exclude != "alerts" {
				return resp, fmt.Errorf("exclude value must be either current, minutely, hourly, daily or alerts")
			}

			// check if the exclude data is unique, with flow is to check if the exclude data is already in the map, if not then add to the map (make data  "marked") and append to the excludeData slice
			if _, ok := excludeUnique[exclude]; !ok {
				excludeUnique[exclude] = true
				excludeData = append(excludeData, exclude)

				exclude_query_value += exclude + "," // add the exclude data to the query value
			}
		}

		// if there any exclude data, then remove the last comma from the exclude_query_value
		if exclude_query_value != "" {
			exclude_query_value = exclude_query_value[:len(exclude_query_value)-1] // remove the last comma
		}
	}

	// http client set
	httpClient := http.Client{}
	req, err := http.NewRequest("GET", OpenWeatherAPIV3OneCall, nil)
	if err != nil {
		return resp, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// set base query param
	q := req.URL.Query()

	q.Add("lat", fmt.Sprintf("%f", query_req.Lat))
	q.Add("lon", fmt.Sprintf("%f", query_req.Lon))
	q.Add("appid", query_req.AppID)

	if query_req.Units != "" {
		q.Add("units", query_req.Units)
	}

	// set additional param
	if len(query_req.Exclude) > 0 {
		q.Add("exclude", exclude_query_value)
	}

	if query_req.Lang != "" {
		q.Add("lang", query_req.Lang)
	}

	// set the query param
	req.URL.RawQuery = q.Encode()

	// request to api
	response, err := httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	// make sure to read the response body
	defer func() {
		if response.StatusCode != http.StatusOK {
			io.ReadAll(response.Body)
		}

		response.Body.Close()
	}()

	// check response and decode response
	if response.StatusCode != http.StatusOK {
		var error OpenWeatherAPIError
		if err := json.NewDecoder(response.Body).Decode(&error); err != nil {
			return resp, err
		}

		return resp, fmt.Errorf("error message: %s, error code: %s, error parameters: %s", error.Message, error.Code, error.Parameters)

	} else {
		if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

// OpenWeatherV3OneCallTimestampAPI fetches historical weather data for a specific timestamp from the OpenWeather API v3 One Call endpoint.
// This function makes a request to the OpenWeatherMap API using the provided query parameters
// and returns the weather information for a specific historical timestamp.
//
// Parameters:
//   - query_req: OpenWeatherAPIV3OneCallTimestampReq containing request parameters such as:
//   - Lat: Latitude (required)
//   - Lon: Longitude (required)
//   - AppID: API key (required)
//   - Dt: Timestamp for historical weather data in Unix UTC format (required)
//   - Units: Units format (optional, defaults to "standard")
//   - Lang: Language for output (optional)
//
// Returns:
//   - OpenWeatherAPIV3OneCallTimestampResp: A structured response containing weather data for the specified timestamp
//   - error: Error if the request fails or validation fails
//
// Example Response Structure:
//
//	{
//	  "lat": 52.2297,                 // Geographical coordinates (latitude)
//	  "lon": 21.0122,                 // Geographical coordinates (longitude)
//	  "timezone": "Europe/Warsaw",    // Local timezone
//	  "timezone_offset": 3600,        // Timezone offset in seconds from UTC
//	  "data": [
//	    {
//	      "dt": 1645888976,           // Data point time (Unix timestamp)
//	      "sunrise": 1645853361,      // Sunrise time (Unix timestamp)
//	      "sunset": 1645891727,       // Sunset time (Unix timestamp)
//	      "temp": 279.13,             // Temperature
//	      "feels_like": 276.44,       // Human perception of temperature
//	      "pressure": 1029,           // Atmospheric pressure (hPa)
//	      "humidity": 64,             // Humidity (%)
//	      "dew_point": 272.88,        // Dew point temperature
//	      "uvi": 0.06,                // UV index
//	      "clouds": 0,                // Cloudiness (%)
//	      "visibility": 10000,        // Average visibility (meters)
//	      "wind_speed": 3.6,          // Wind speed
//	      "wind_deg": 340,            // Wind direction (degrees)
//	      "weather": [
//	        {
//	          "id": 800,              // Weather condition ID
//	          "main": "Clear",        // Group of weather parameters
//	          "description": "clear sky", // Weather description
//	          "icon": "01d"           // Weather icon ID
//	        }
//	      ]
//	    }
//	  ]
//	}
//
// Errors:
//   - Returns error if required parameters (Lat, Lon, AppID, Dt) are missing
//   - Returns error if API request fails or returns non-OK status
//
// source: https://openweathermap.org/api/one-call-3
func OpenWeatherV3OneCallTimestampAPI(query_req OpenWeatherAPIV3OneCallTimestampReq) (OpenWeatherAPIV3OneCallTimestampResp, error) {
	var resp OpenWeatherAPIV3OneCallTimestampResp

	// base value checker
	base_req := OpenWeatherAPIV3OneCallBaseReq{
		Lat:   query_req.Lat,
		Lon:   query_req.Lon,
		AppID: query_req.AppID,
		Units: query_req.Units,
	}

	if err := checkBaseQuery(base_req); err != nil {
		return resp, err
	}

	// additional value checker
	if query_req.Dt == 0 {
		return resp, fmt.Errorf("dt is required")
	}

	// http client setup
	httpClient := http.Client{}
	req, err := http.NewRequest("GET", OpenWeatherAPIV3OneCallTimestamp, nil)
	if err != nil {
		return resp, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// set base query param
	q := req.URL.Query()

	q.Add("lat", fmt.Sprintf("%f", query_req.Lat))
	q.Add("lon", fmt.Sprintf("%f", query_req.Lon))
	q.Add("appid", query_req.AppID)

	if query_req.Units != "" {
		q.Add("units", query_req.Units)
	}

	// set additional query param
	q.Add("dt", fmt.Sprintf("%d", query_req.Dt))

	if query_req.Lang != "" {
		q.Add("lang", query_req.Lang)
	}

	// set the query param
	req.URL.RawQuery = q.Encode()

	// request to api
	response, err := httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	// make sure to read the response body
	defer func() {
		if response.StatusCode != http.StatusOK {
			io.ReadAll(response.Body)
		}

		response.Body.Close()
	}()

	// check response and decode response
	if response.StatusCode != http.StatusOK {
		var error OpenWeatherAPIError
		if err := json.NewDecoder(response.Body).Decode(&error); err != nil {
			return resp, err
		}

		return resp, fmt.Errorf("error message: %s, error code: %s, error parameters: %s", error.Message, error.Code, error.Parameters)

	} else {
		if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

// OpenWeatherV3OneCallDailySummaryAPI fetches summarized daily weather information for a specific location and date from the OpenWeatherMap API.
//
// This function makes an HTTP GET request to the OpenWeatherMap API Daily Summary endpoint and returns weather data
// for the specified coordinates and date. The response includes details such as temperature, cloud cover,
// humidity, precipitation, pressure, and wind information for the requested day.
//
// Parameters:
//   - query_req: OpenWeatherAPIV3OneCallDailySummaryReq containing required parameters:
//   - Lat: Location latitude (required)
//   - Lon: Location longitude (required)
//   - AppID: Your OpenWeatherMap API key (required)
//   - Date: Date for which to get the weather summary in format YYYY-MM-DD (required)
//   - Optional parameters:
//   - Units: Units format (standard, metric, imperial)
//   - Lang: Language for text information
//
// Returns:
//   - OpenWeatherAPIV3OneCallDailySummaryResp: Contains the weather data with fields like:
//   - lat/lon: Geographical coordinates
//   - tz: Timezone
//   - date: Date of weather data
//   - units: Units format
//   - cloud_cover: Cloud coverage information
//   - humidity: Humidity values
//   - precipitation: Precipitation information
//   - temperature: Temperature values (min, max, afternoon, night, evening, morning)
//   - pressure: Atmospheric pressure
//   - wind: Wind information including maximum speed and direction
//   - error: An error if the request fails or if required parameters are missing
//
// Example response:
//
//	{
//	  "lat":33,
//	  "lon":35,
//	  "tz":"+02:00",
//	  "date":"2020-03-04",
//	  "units":"standard",
//	  "cloud_cover":{
//	     "afternoon":0
//	  },
//	  "humidity":{
//	     "afternoon":33
//	  },
//	  "precipitation":{
//	     "total":0
//	  },
//	  "temperature":{
//	     "min":286.48,
//	     "max":299.24,
//	     "afternoon":296.15,
//	     "night":289.56,
//	     "evening":295.93,
//	     "morning":287.59
//	  },
//	  "pressure":{
//	     "afternoon":1015
//	  },
//	  "wind":{
//	     "max":{
//	        "speed":8.7,
//	        "direction":120
//	     }
//	  }
//	}
//
// source: https://openweathermap.org/api/one-call-3
func OpenWeatherV3OneCallDailySummaryAPI(query_req OpenWeatherAPIV3OneCallDailySummaryReq) (OpenWeatherAPIV3OneCallDailySummaryResp, error) {

	var resp OpenWeatherAPIV3OneCallDailySummaryResp

	// base value checker
	base_req := OpenWeatherAPIV3OneCallBaseReq{
		Lat:   query_req.Lat,
		Lon:   query_req.Lon,
		AppID: query_req.AppID,
		Units: query_req.Units,
	}

	if err := checkBaseQuery(base_req); err != nil {
		return resp, err
	}

	// additional value checker
	if query_req.Date == "" {
		return resp, fmt.Errorf("date is required")
	}

	// * make a request to the NewsAPI with the given parameters
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", OpenWeatherAPIV3OneCallDailySummary, nil)
	if err != nil {
		return resp, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// add base query values to the request
	q := req.URL.Query()

	// base query values
	q.Add("lat", fmt.Sprintf("%f", query_req.Lat))
	q.Add("lon", fmt.Sprintf("%f", query_req.Lon))
	q.Add("appid", query_req.AppID)

	if query_req.Units != "" {
		q.Add("units", query_req.Units)
	}

	// add additional query base on type of request
	q.Add("date", query_req.Date)

	if query_req.Lang != "" {
		q.Add("lang", query_req.Lang)
	}

	// assign the query values to the request
	req.URL.RawQuery = q.Encode()

	// send request to the server
	response, err := httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	// make sure read the response body
	defer func() {
		if response.StatusCode != http.StatusOK {
			io.ReadAll(response.Body)
		}

		response.Body.Close()
	}()

	// Decode response from the read body and check for errors
	if response.StatusCode != http.StatusOK {
		var error OpenWeatherAPIError
		if err := json.NewDecoder(response.Body).Decode(&error); err != nil {
			return resp, err
		}

		return resp, fmt.Errorf("error message: %s, error code: %s, error parameters: %s", error.Message, error.Code, error.Parameters)

	} else {
		if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

// OpenWeatherV3OneCallOverviewAPI makes a request to the OpenWeather API V3 One Call Overview
// endpoint to retrieve a comprehensive overview of weather conditions for a specific location.
//
// This function takes an OpenWeatherAPIV3OneCallOverviewReq struct as input, which must
// contain the latitude, longitude, and API key. Optional parameters include units of
// measurement and a specific date for historical data.
//
// The function returns an OpenWeatherAPIV3OneCallOverviewResp struct containing detailed
// weather information including current conditions presented as a natural language summary.
// The response includes geographical coordinates, timezone, date, units, and a comprehensive
// weather overview text that describes current conditions, temperature, feels-like temperature,
// wind data, atmospheric conditions, visibility, UV index, cloud coverage, precipitation
// outlook, and an overall weather assessment.
//
// Example response:
//
//	{
//	  "lat": 51.509865,
//	  "lon": -0.118092,
//	  "tz": "+01:00",
//	  "date": "2024-05-13",
//	  "units": "metric",
//	  "weather_overview": "The current weather is overcast with a temperature of 16°C and a feels-like temperature of 16°C..."
//	}
//
// If an error occurs during the API request or response processing, the function returns
// an empty response struct and the corresponding error.
//
// source: https://openweathermap.org/api/one-call-3
func OpenWeatherV3OneCallOverviewAPI(query_req OpenWeatherAPIV3OneCallOverviewReq) (OpenWeatherAPIV3OneCallOverviewResp, error) {

	var resp OpenWeatherAPIV3OneCallOverviewResp

	// base value checker
	base_req := OpenWeatherAPIV3OneCallBaseReq{
		Lat:   query_req.Lat,
		Lon:   query_req.Lon,
		AppID: query_req.AppID,
		Units: query_req.Units,
	}

	if err := checkBaseQuery(base_req); err != nil {
		return resp, err
	}

	// * make a request to the NewsAPI with the given parameters
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", OpenWeatherAPIV3OneCallOverview, nil)
	if err != nil {
		return resp, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// get current query values and add the base query values to the request
	q := req.URL.Query()

	// base query values
	q.Add("lat", fmt.Sprintf("%f", query_req.Lat))
	q.Add("lon", fmt.Sprintf("%f", query_req.Lon))
	q.Add("appid", query_req.AppID)

	if query_req.Units != "" {
		q.Add("units", query_req.Units)
	}

	// add additional query base on type of request
	if query_req.Date != "" {
		q.Add("date", query_req.Date)
	}

	// assign the query values to the request
	req.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	// read the response body
	defer func() {
		if response.StatusCode != http.StatusOK {
			io.ReadAll(response.Body)
		}

		response.Body.Close()
	}()

	// Decode response from the read body and check for errors
	if response.StatusCode != http.StatusOK {
		var error OpenWeatherAPIError
		if err := json.NewDecoder(response.Body).Decode(&error); err != nil {
			return resp, err
		}

		return resp, fmt.Errorf("error message: %s, error code: %s, error parameters: %s", error.Message, error.Code, error.Parameters)

	} else {
		if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
			return resp, err
		}
	}

	return resp, nil
}
