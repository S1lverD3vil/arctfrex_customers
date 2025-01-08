package email

type Email struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type GomailSender interface {
	SendEmail(email Email) error
}
