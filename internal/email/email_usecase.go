package email

type EmailUseCase interface {
	SendEmail(email Email) error
}

type emailUseCase struct {
	gomailSender GomailSender
}

func NewEmailUseCase(gs GomailSender) *emailUseCase {
	return &emailUseCase{gomailSender: gs}
}

func (eu *emailUseCase) SendEmail(email Email) error {
	return eu.gomailSender.SendEmail(email)
}
