package otp

import (
	"arctfrex-customers/internal/base"

	"gorm.io/gorm"
)

type otpRepository struct {
	db *gorm.DB
}

func NewOtpRepository(db *gorm.DB) OtpRepository {
	return &otpRepository{db: db}
}

func (or *otpRepository) Save(otp *Otp) error {
	return or.db.Save(&otp).Error
}

func (or *otpRepository) GetActiveOtpBySendTo(mobilePhone, otpType, otpProcess string) (*Otp, error) {
	var otp Otp

	queryParams := Otp{
		Type:    otpType,
		SendTo:  mobilePhone,
		Process: otpProcess,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := or.db.Where(&queryParams).First(&otp).Error; err != nil {
		return nil, err
	}
	return &otp, nil
}
