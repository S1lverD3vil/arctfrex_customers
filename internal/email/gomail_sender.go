package email

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type gomailSender struct {
	Dialer *gomail.Dialer
}

func NewGomailSender() *gomailSender {
	// Retrieve the value of an environment variable
	googleSmtpServer := os.Getenv("GOOGLE_SMTP_SERVER")
	googleSmtpUsername := os.Getenv("GOOGLE_SMTP_USERNAME")
	googleSmtpPassword := os.Getenv("GOOGLE_SMTP_PASSWORD")
	googleSmtpPort, err := strconv.ParseInt(os.Getenv("GOOGLE_SMTP_PORT"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}

	return &gomailSender{
		Dialer: gomail.NewDialer(googleSmtpServer, int(googleSmtpPort), googleSmtpUsername, googleSmtpPassword),
	}
}

func (g *gomailSender) SendEmail(email Email) error {
	m := gomail.NewMessage()
	m.SetHeader("From", email.From)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/plain", email.Body)

	return g.Dialer.DialAndSend(m)
}
