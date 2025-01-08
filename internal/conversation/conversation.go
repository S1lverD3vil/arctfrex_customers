package conversation

import (
	"arctfrex-customers/internal/base"
	userBackoffice "arctfrex-customers/internal/user/backoffice"
	userMobile "arctfrex-customers/internal/user/mobile"
)

const (
	USER     = "USER"
	OPERATOR = "OPERATOR"
	SYSTEM   = "SYSTEM"
)

// ConversationSession represents a user/operator session
type ConversationSession struct {
	ID               string                `gorm:"primaryKey" json:"session_id"`
	UserID           string                `json:"user_id"`                              // Reference to the user
	BackofficeUserID string                `json:"backoffice_user_id"`                   // Reference to the backoffice user
	Messages         []ConversationMessage `gorm:"foreignKey:SessionID" json:"messages"` // Messages in the session
	CurrentStepID    string                `json:"current_step_id"`                      // Tracks the current step of the conversation

	User           *userMobile.Users               `gorm:"-" json:"user"`
	BackofficeUser *userBackoffice.BackofficeUsers `gorm:"-" json:"backoffice_user"`

	base.BaseModel
}

// ConversationMessage represents a single message in a conversation
type ConversationMessage struct {
	ID            string `gorm:"primary_key" json:"message_id"`
	Content       string `json:"content"`         // The message content
	SessionID     string `json:"session_id"`      // Reference to the conversation session
	CreatedByID   string `json:"created_by_id"`   // Reference to user/operator ID
	CreatedByType string `json:"created_by_type"` // Indicates "Users" or "BackofficeUsers"
	FromUser      string `json:"from_user"`       // Reference to human/system

	CreatedBy interface{} `gorm:"-" json:"created_by"` // Dynamically loaded field for either Users or BackofficeUsers

	base.BaseModel
}

// ConversationStep represents a step in the conversation tree
type ConversationStep struct {
	ID      string `gorm:"primary_key" json:"step_id"`
	Content string `json:"content"` // Text to display to the user

	Options []ConversationOption `gorm:"foreignKey:StepID;references:ID" json:"options"` // Possible options for the step
}

// ConversationOption represents a possible response at a conversation step
type ConversationOption struct {
	ID         string  `gorm:"primary_key"`
	Number     string  `json:"number"`
	Content    string  `json:"content"`      // Option displayed to the user
	NextStepID *string `json:"next_step_id"` // Points to the next step
	StepID     string  `json:"step_id"`      // Step this option belongs to

	ConversationStep ConversationStep `gorm:"foreignKey:StepID;references:ID"`
}

type ConversationRepository interface {
	GetConversationSessions() ([]ConversationSession, error)
	GetConversationSessionBySessionID(conversationSessionId string) (*ConversationSession, error)
	TakeConversationSessionBySessionID(conversationSessionId string, backOfficeUserID string) error
	GetConversationSessionsByBackofficeUserID(backofficeUserID string) ([]ConversationSession, error)
	GetUserByID(userId string) (*userMobile.Users, error)
	GetBackofficeUserByID(userId string) (*userBackoffice.BackofficeUsers, error)
	CreateSession(conversation *ConversationSession) error
	GetActiveSessionByUserID(userID string) (*ConversationSession, error)
	GetSessionByID(sessionID string, session *ConversationSession) error
	UpdateSession(session *ConversationSession) error // Add this
	GetOptionByNumber(number string, option *ConversationOption) error
	GetStepByID(stepID string, step *ConversationStep) error
	GetOptionByNumberAndStep(number string, stepID string, option *ConversationOption) error // Scoped lookup
	SaveMessage(message *ConversationMessage) error
	EndConversationSessionBySessionID(sessionID string, endedBy string) error
	SeedConversationSteps()
}
