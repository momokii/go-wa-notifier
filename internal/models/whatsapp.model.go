package models

type NewsSendWhatsappReq struct {
	WhatsappNumbers []string `json:"whatsapp_numbers" example:"6285727771234,6285667889887"` // list of numbers to send the news to and start with code number like 62 and not 0 like 08123456789
	Category        string   `json:"category" example:"business"`                            // options: business, entertainment, general, health, science, sports, technology
	UsingLLM        bool     `json:"using_llm" example:"true"`                               // options: true, false, if set to true, the message news will be add with llm and if false, the message news will be add with the default message
}

type WhatsappMessagesReq struct {
	Messages        string   `json:"messages" example:"Hello, this is a test message"`       // message to be sent to the whatsapp numbers
	WhatsappNumbers []string `json:"whatsapp_numbers" example:"6285727771234,6285667889887"` // list of numbers to send the news to and start with code number like 62 and not 0 like 08123456789
}

type WeatherSendWhatsappReq struct {
	Type            string   `json:"type" example:"today"`                                   // required, options: today, tomorrow
	Lat             float64  `json:"lat" example:"-6.2617"`                                  // required, latitude of the location
	Lon             float64  `json:"lon" example:"106.8103"`                                 // required, longitude of the location
	WhatsappNumbers []string `json:"whatsapp_numbers" example:"6285727771234,6285667889887"` // list of numbers to send the news to and start with code number like 62 and not 0 like 08123456789
	UsingLLM        bool     `json:"using_llm" example:"true"`                               // options: true, false, if set to true, the message news will be add with llm and if false, the message news will be add with the default message
}
