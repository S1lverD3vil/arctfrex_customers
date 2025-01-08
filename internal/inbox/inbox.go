package inbox

import "arctfrex-customers/internal/base"

type Inbox struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   string `json:"userid"`
	Receiver string `gorm:"index;not null" json:"receiver"` // Email or User ID
	Sender   string `gorm:"not null" json:"sender"`         // Email or User ID
	Subject  string `gorm:"size:255" json:"subject"`
	Message  string `gorm:"type:text" json:"message"`
	IsRead   bool   `gorm:"default:false" json:"is_read"`
	// CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	// UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	base.BaseModel
}

type InboxRepository interface {
	Create(inbox *Inbox) error
	GetAll(receiver string) ([]Inbox, error)
	GetAllInboxesByUserId(userId string) ([]Inbox, error)
	MarkAsRead(userId string, id uint) error
	Delete(userId string, id uint) error
}
