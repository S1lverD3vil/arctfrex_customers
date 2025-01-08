package whatsapp

import (
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type twilioWhatsappSender struct {
	client *twilio.RestClient
}

func NewTwilioWhatsappSender() *twilioWhatsappSender {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_WHATSAPP_USERNAME"),
		Password: os.Getenv("TWILIO_WHATSAPP_PASSWORD"),
	})

	return &twilioWhatsappSender{
		client: client,
	}
}

func (tws *twilioWhatsappSender) SendWhatsapp(whatsapp TwilioWhatsapp) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(whatsapp.To)
	params.SetFrom(whatsapp.From)
	params.SetBody(whatsapp.Body)

	_, err := tws.client.Api.CreateMessage(params)
	return err
}
