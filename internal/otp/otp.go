package otp

import (
	"arctfrex-customers/internal/base"
	"time"
)

type Otp struct {
	Code     string    `json:"code"`
	SendTo   string    `gorm:"primaryKey;index:idx_sendto_type_process,unique" json:"send_to"`
	Type     string    `gorm:"primaryKey;index:idx_sendto_type_process,unique" json:"type"`
	Process  string    `gorm:"primaryKey;index:idx_sendto_type_process,unique" json:"process"`
	IsUsed   bool      `json:"is_used"`
	UsedTime time.Time `json:"used_time"`

	base.BaseModel
}

type OtpRepository interface {
	Save(otp *Otp) error
	GetActiveOtpBySendTo(mobilePhone, otpType, otpProcess string) (*Otp, error)
}
