package models

type NewsSendWhatsappReq struct {
	WhatsappNumbers []string `json:"whatsapp_numbers" example:"6285727771234,6285667889887"` // list of numbers to send the news to and start with code number like 62 and not 0 like 08123456789
	Category        string   `json:"category" example:"business"`                            // options: business, entertainment, general, health, science, sports, technology
}
