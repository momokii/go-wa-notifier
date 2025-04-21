package openweatherapi

const (
	OpenWeatherAPIBaseURL               = "https://api.openweathermap.org"
	OpenWeatherAPIV3                    = OpenWeatherAPIBaseURL + "/data/3.0"
	OpenWeatherAPIV3OneCall             = OpenWeatherAPIV3 + "/onecall"
	OpenWeatherAPIV3OneCallOverview     = OpenWeatherAPIV3OneCall + "/overview"
	OpenWeatherAPIV3OneCallDailySummary = OpenWeatherAPIV3OneCall + "/day_summary"
	OpenWeatherAPIV3OneCallTimestamp    = OpenWeatherAPIV3OneCall + "/timemachine"
)

// source: https://openweathermap.org/api/one-call-3

// ====== REQUEST STRUCTURES ====== //
// OpenWeatherAPIV3OneCallReq is the request structure for the OpenWeather API v3 One Call endpoint.

type OpenWeatherAPIV3OneCallBaseReq struct {
	Lat   float64 `json:"lat"`   // required,Latitude, decimal (-90; 90). If you need the geocoder to automatic convert city names and zip-codes to geo coordinates and the other way around, use the geocoding API.
	Lon   float64 `json:"lon"`   // required, Longitude, decimal (-180; 180). If you need the geocoder to automatic convert city names and zip-codes to geo coordinates and the other way around, use the geocoding API.
	AppID string  `json:"appid"` // required, Your unique API key (you can always find it on your account page under the
	Units string  `json:"units"` // optional, units of measurement. Possible values: standard, metric, imperial. Default is standard. (Fahrenheit=imperial, Celsius=metric, Kelvin=standard)
}

type OpenWeatherAPIV3OneCallReq struct {
	OpenWeatherAPIV3OneCallBaseReq
	Exclude []string `json:"exclude"` // optional, By using this parameter you can exclude some parts of the weather data from the API response. It should be a comma-delimited list (without spaces). Possible values: current, minutely, hourly, daily, alerts. Default is empty.
	Lang    string   `json:"lang"`    // optional, language of the response. Default is en.
}

type OpenWeatherAPIV3OneCallTimestampReq struct {
	OpenWeatherAPIV3OneCallBaseReq
	Dt   int64  `json:"dt"`   // required, Timestamp (Unix time, UTC time zone), e.g. dt=1586468027. Data is available from January 1st, 1979 till 4 days ahead of the current date.
	Lang string `json:"lang"` // optional, language of the response. Default is en.
}

type OpenWeatherAPIV3OneCallDailySummaryReq struct {
	OpenWeatherAPIV3OneCallBaseReq
	Date string `json:"date"` // required, Date in the `YYYY-MM-DD` format for which data is requested. Date is available for 46+ years archive (starting from 1979-01-02) up to the 1,5 years ahead forecast to the current date
	Lang string `json:"lang"` // optional, language of the response. Default is en.
}

type OpenWeatherAPIV3OneCallOverviewReq struct {
	OpenWeatherAPIV3OneCallBaseReq
	Date string `json:"date"` // optional, The date the user wants to get a weather summary in the YYYY-MM-DD format. Data is available for today and tomorrow. If not specified, the current date will be used by default. Please note that the date is determined by the timezone relevant to the coordinates specified in the API request
}

// ====== RESPONSE STRUCTURES ====== //
// OpenWeatherAPIV3OneCallRes is the response structure for the OpenWeather API v3 One Call endpoint.

// one call response
type OpenWeatherAPIV3OneCallResp struct {
	Lat            float64        `json:"lat"`
	Lon            float64        `json:"lon"`
	Timezone       string         `json:"timezone"`
	TimezoneOffset int            `json:"timezone_offset"`
	Current        CurrentData    `json:"current"`
	Minutely       []MinutelyData `json:"minutely"`
	Hourly         []HourlyData   `json:"hourly"`
	Daily          []DailyData    `json:"daily"`
	Alerts         []AlertData    `json:"alerts"`
}

type CurrentData struct {
	Dt         int64         `json:"dt"`
	Sunrise    int64         `json:"sunrise"`
	Sunset     int64         `json:"sunset"`
	Temp       float64       `json:"temp"`
	FeelsLike  float64       `json:"feels_like"`
	Pressure   int           `json:"pressure"`
	Humidity   int           `json:"humidity"`
	DewPoint   float64       `json:"dew_point"`
	Uvi        float64       `json:"uvi"`
	Clouds     int           `json:"clouds"`
	Visibility int           `json:"visibility"`
	WindSpeed  float64       `json:"wind_speed"`
	WindDeg    int           `json:"wind_deg"`
	WindGust   float64       `json:"wind_gust"`
	Weather    []WeatherData `json:"weather"`
}

type MinutelyData struct {
	Dt            int64   `json:"dt"`
	Precipitation float64 `json:"precipitation"`
}

type HourlyData struct {
	Dt         int64         `json:"dt"`
	Temp       float64       `json:"temp"`
	FeelsLike  float64       `json:"feels_like"`
	Pressure   int           `json:"pressure"`
	Humidity   int           `json:"humidity"`
	DewPoint   float64       `json:"dew_point"`
	Uvi        float64       `json:"uvi"`
	Clouds     int           `json:"clouds"`
	Visibility int           `json:"visibility"`
	WindSpeed  float64       `json:"wind_speed"`
	WindDeg    int           `json:"wind_deg"`
	WindGust   float64       `json:"wind_gust"`
	Weather    []WeatherData `json:"weather"`
	Pop        float64       `json:"pop"`
}

type DailyData struct {
	Dt        int64         `json:"dt"`
	Sunrise   int64         `json:"sunrise"`
	Sunset    int64         `json:"sunset"`
	Moonrise  int64         `json:"moonrise"`
	Moonset   int64         `json:"moonset"`
	MoonPhase float64       `json:"moon_phase"`
	Summary   string        `json:"summary"`
	Temp      TempData      `json:"temp"`
	FeelsLike FeelsLikeData `json:"feels_like"`
	Pressure  int           `json:"pressure"`
	Humidity  int           `json:"humidity"`
	DewPoint  float64       `json:"dew_point"`
	WindSpeed float64       `json:"wind_speed"`
	WindDeg   int           `json:"wind_deg"`
	WindGust  float64       `json:"wind_gust"`
	Weather   []WeatherData `json:"weather"`
	Clouds    int           `json:"clouds"`
	Pop       float64       `json:"pop"`
	Rain      float64       `json:"rain,omitempty"`
	Uvi       float64       `json:"uvi"`
}

type TempData struct {
	Day   float64 `json:"day"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type FeelsLikeData struct {
	Day   float64 `json:"day"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type WeatherData struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type AlertData struct {
	SenderName  string   `json:"sender_name"`
	Event       string   `json:"event"`
	Start       int64    `json:"start"`
	End         int64    `json:"end"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// OpenWeatherAPIV3OneCallResp is the response structure for the OpenWeather API v3 One Call endpoint.
type OpenWeatherAPIV3OneCallTimestampResp struct {
	Lat            float64                    `json:"lat"`
	Lon            float64                    `json:"lon"`
	Timezone       string                     `json:"timezone"`
	TimezoneOffset int                        `json:"timezone_offset"`
	Data           []WeatherDataTimestampResp `json:"data"`
}

type WeatherDataTimestampResp struct {
	Dt         int64           `json:"dt"`
	Sunrise    int64           `json:"sunrise"`
	Sunset     int64           `json:"sunset"`
	Temp       float64         `json:"temp"`
	FeelsLike  float64         `json:"feels_like"`
	Pressure   int             `json:"pressure"`
	Humidity   int             `json:"humidity"`
	DewPoint   float64         `json:"dew_point"`
	UVI        float64         `json:"uvi"`
	Clouds     int             `json:"clouds"`
	Visibility int             `json:"visibility"`
	WindSpeed  float64         `json:"wind_speed"`
	WindDeg    int             `json:"wind_deg"`
	Weather    []WeatherDetail `json:"weather"`
}

type WeatherDetail struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// daily summary response
type OpenWeatherAPIV3OneCallDailySummaryResp struct {
	Lat           float64           `json:"lat"`
	Lon           float64           `json:"lon"`
	TZ            string            `json:"tz"`
	Date          string            `json:"date"`
	Units         string            `json:"units"`
	CloudCover    CloudCoverData    `json:"cloud_cover"`
	Humidity      HumidityData      `json:"humidity"`
	Precipitation PrecipitationData `json:"precipitation"`
	Temperature   TemperatureData   `json:"temperature"`
	Pressure      PressureData      `json:"pressure"`
	Wind          WindData          `json:"wind"`
}

type CloudCoverData struct {
	Afternoon float64 `json:"afternoon"`
}

type HumidityData struct {
	Afternoon float64 `json:"afternoon"`
}

type PrecipitationData struct {
	Total float64 `json:"total"`
}

type TemperatureData struct {
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
	Afternoon float64 `json:"afternoon"`
	Night     float64 `json:"night"`
	Evening   float64 `json:"evening"`
	Morning   float64 `json:"morning"`
}

type PressureData struct {
	Afternoon float64 `json:"afternoon"`
}

type WindData struct {
	Max WindDetail `json:"max"`
}

type WindDetail struct {
	Speed     float64 `json:"speed"`
	Direction float64 `json:"direction"`
}

// overview response
type OpenWeatherAPIV3OneCallOverviewResp struct {
	Lat             float64 `json:"lat"`
	Lon             float64 `json:"lon"`
	TZ              string  `json:"tz"`
	Date            string  `json:"date"`
	Units           string  `json:"units"`
	WeatherOverview string  `json:"weather_overview"`
}

// ===== ERROR STRUCTURES ====== //
type OpenWeatherAPIError struct {
	Code       string   `json:"code"`
	Message    string   `json:"message"`
	Parameters []string `json:"parameters"`
}

// ===== AGGREGATED WEATHER DATA STRUCTURE ===== //
// WeatherDataAggregate combines all the weather information for easier handling
type WeatherDataAggregate struct {
	Date             string                                  `json:"date"`
	ReportType       string                                  `json:"report_type"` // "today" or "tomorrow"
	Latitude         float64                                 `json:"latitude"`
	Longitude        float64                                 `json:"longitude"`
	WeatherOverview  string                                  `json:"weather_overview"`
	Timezone         string                                  `json:"timezone"`
	DailyAggregate   OpenWeatherAPIV3OneCallDailySummaryResp `json:"daily_aggregate"`
	HourlyForecast   []HourlyData                            `json:"hourly_forecast"`
	CurrentTimeLocal string                                  `json:"current_time_local"`
}
