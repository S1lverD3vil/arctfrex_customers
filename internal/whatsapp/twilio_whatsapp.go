package whatsapp

type TwilioWhatsapp struct {
	From string `json:"from"`
	To   string `json:"to"`
	Body string `json:"body"`
}

type TwilioWhatsappSender interface {
	SendWhatsapp(whatsapp TwilioWhatsapp) error
}
