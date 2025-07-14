package otp

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/pquerna/otp/totp"

	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/email"
	"arctfrex-customers/internal/repository"
	"arctfrex-customers/internal/whatsapp"
)

type OtpUsecase interface {
	SendOtp(otp *Otp) error
	ValidateOtp(otp *Otp) error
}

type otpUsecase struct {
	otpRepository        OtpRepository
	userRepository       repository.UserRepository
	twilioWhatsappSender whatsapp.TwilioWhatsappSender
	gomailSender         email.GomailSender
}

func NewOtpUseCase(
	or OtpRepository,
	us repository.UserRepository,
	tws whatsapp.TwilioWhatsappSender,
	gs email.GomailSender,
) *otpUsecase {
	return &otpUsecase{
		otpRepository:        or,
		userRepository:       us,
		twilioWhatsappSender: tws,
		gomailSender:         gs,
	}
}

func (ou *otpUsecase) SendOtp(otp *Otp) error {
	secret := os.Getenv(common.OTP_GENERATOR_SECRET) // This should be kept secure and private
	emailFrom := os.Getenv(common.EMAIL_FROM)
	emailSubject := os.Getenv(common.OTP_EMAIL_SUBJECT)
	isSendEmail, err := strconv.ParseBool(os.Getenv(common.OTP_SEND_WITH_EMAIL))
	if err != nil {
		return err
	}

	// user, err := ou.userRepository.GetActiveUserByMobilePhone(otp.SendTo)
	user, err := ou.userRepository.GetUserByMobilePhone(otp.SendTo)
	if user == nil || err != nil {
		return errors.New("user not found")
	}

	otpDb, _ := ou.otpRepository.GetActiveOtpBySendTo(user.Email, otp.Type, otp.Process)

	if otpDb != nil && !otpDb.IsUsed && time.Since(otpDb.ModifiedAt).Minutes() <= 2 {
		return errors.New("otp exists please wait 2 minutes to request new")
	}

	// Generate a TOTP using the secret key, 6 digits, and a 30-second interval
	otpCode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period: 30,
		Digits: 4,
	})
	if err != nil {
		return err
	}

	otp.Code = otpCode
	otp.IsActive = true
	otp.CreatedBy = otp.SendTo
	if otpDb != nil {
		otp.CreatedAt = otpDb.CreatedAt
		otp.ModifiedBy = otp.SendTo
	}

	if isSendEmail {
		otp.SendTo = user.Email
		go ou.gomailSender.SendEmail(email.Email{
			From:    emailFrom,
			To:      otp.SendTo,
			Subject: emailSubject,
			Body:    "Gunakan kode OTP " + otpCode + " untuk verifikasi akun PaNen anda. Jangan bagikan kode OTP Anda ke orang lain. ",
		})

		return ou.otpRepository.Save(otp)
	}

	ou.twilioWhatsappSender.SendWhatsapp(whatsapp.TwilioWhatsapp{
		To:   "whatsapp:+62" + otp.SendTo,
		From: "whatsapp:+14155238886",
		Body: "Gunakan kode OTP " + otpCode + " untuk verifikasi akun PaNen anda. Jangan bagikan kode OTP Anda ke orang lain. ",
	})

	return ou.otpRepository.Save(otp)
}

func (ou *otpUsecase) ValidateOtp(otp *Otp) error {
	// user, err := ou.userRepository.GetActiveUserByMobilePhone(otp.SendTo)
	user, err := ou.userRepository.GetUserByMobilePhone(otp.SendTo)
	if user == nil || err != nil {
		return errors.New("user not found")
	}

	otpDb, err := ou.otpRepository.GetActiveOtpBySendTo(user.Email, otp.Type, otp.Process)
	if otpDb == nil || err != nil || otpDb.IsUsed || (!otpDb.IsUsed && time.Since(otpDb.ModifiedAt).Minutes() > 2) {
		return errors.New("record not found")
	}

	if otp.Code != otpDb.Code {
		return errors.New("otp not valid")
	}

	otpDb.IsUsed = true
	otpDb.UsedTime = time.Now()

	return ou.otpRepository.Save(otpDb)
}
