package conversation

import (
	userBackoffice "arctfrex-customers/internal/user/backoffice"
	userMobile "arctfrex-customers/internal/user/mobile"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type conversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{db: db}
}

func (cr *conversationRepository) GetConversationSessions() ([]ConversationSession, error) {
	var sessions []ConversationSession

	// Fetch all conversation sessions (without the user and backoffice user)
	if err := cr.db.Where("is_active = ?", true).Find(&sessions).Error; err != nil {
		return nil, err
	}

	// Now, manually load the related user and backoffice user data
	for i := range sessions {
		// Load the user data by MobilePhone or other reference
		var user userMobile.Users
		if sessions[i].UserID != "" {
			cr.db.Where("id = ?", sessions[i].UserID).First(&user)
			sessions[i].User = &user
		}

		// Load the backoffice user data
		var backofficeUser userBackoffice.BackofficeUsers
		if sessions[i].BackofficeUserID != "" {
			cr.db.Where("id = ?", sessions[i].BackofficeUserID).First(&backofficeUser)
			sessions[i].BackofficeUser = &backofficeUser
		}
	}

	return sessions, nil
}

func (cr *conversationRepository) GetConversationSessionsByBackofficeUserID(backofficeUserID string) ([]ConversationSession, error) {
	var sessions []ConversationSession

	// Fetch all conversation sessions (without the user and backoffice user)
	if err := cr.db.Where("is_active = ? AND backoffice_user_id = ?", true, backofficeUserID).Find(&sessions).Error; err != nil {
		return nil, err
	}

	// Now, manually load the related user and backoffice user data
	for i := range sessions {
		// Load the user data by MobilePhone or other reference
		var user userMobile.Users
		if sessions[i].UserID != "" {
			cr.db.Where("id = ?", sessions[i].UserID).First(&user)
			sessions[i].User = &user
		}

		// Load the backoffice user data
		var backofficeUser userBackoffice.BackofficeUsers
		if sessions[i].BackofficeUserID != "" {
			cr.db.Where("id = ?", sessions[i].BackofficeUserID).First(&backofficeUser)
			sessions[i].BackofficeUser = &backofficeUser
		}
	}

	return sessions, nil
}

func (cr *conversationRepository) GetConversationSessionBySessionID(conversationSessionId string) (*ConversationSession, error) {
	var session ConversationSession

	if err := cr.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC") // Sort messages by created_at in ascending order
	}).First(&session, "id = ? AND is_active = ?", conversationSessionId, true).Error; err != nil {
		return nil, err
	}

	// Load the user data by MobilePhone or other reference
	var user userMobile.Users
	cr.db.Where("id = ?", session.UserID).First(&user)
	session.User = &user

	// Load the backoffice user data
	var backofficeUser userBackoffice.BackofficeUsers
	cr.db.Where("id = ?", session.BackofficeUserID).First(&backofficeUser)
	session.BackofficeUser = &backofficeUser

	return &session, nil
}

func (cr *conversationRepository) TakeConversationSessionBySessionID(conversationSessionId string, backOfficeUserID string) error {
	conversationSession, err := cr.GetConversationSessionBySessionID(conversationSessionId)
	if err != nil {
		return err
	}

	return cr.db.Model(&conversationSession).Where("id = ?", conversationSessionId).Update("backoffice_user_id", backOfficeUserID).Error
}

func (cr *conversationRepository) GetUserByID(id string) (*userMobile.Users, error) {
	var user userMobile.Users
	if err := cr.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (cr *conversationRepository) GetBackofficeUserByID(id string) (*userBackoffice.BackofficeUsers, error) {
	var userBackoffice userBackoffice.BackofficeUsers
	if err := cr.db.First(&userBackoffice, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &userBackoffice, nil
}

func (cr *conversationRepository) CreateSession(conversationSession *ConversationSession) error {
	// Set default values base model
	conversationSession.IsActive = true // Set is_active to true on creation
	conversationSession.CreatedAt = time.Now()
	conversationSession.ModifiedAt = time.Now()

	// Start a new DB transaction
	tx := cr.db.Begin()

	// Set the initilal step ID to 1
	conversationSession.CurrentStepID = "1"

	// Save the conversation session
	if err := tx.Save(conversationSession).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Fetch the first conversation step (assuming step with ID 1 is the first step)
	var firstStep ConversationStep
	if err := tx.First(&firstStep, "1").Error; err != nil {
		tx.Rollback()
		return err
	}

	// Initialize the message content with the first step's content
	messageContent := firstStep.Content + "\n\n"

	// Fetch all options related to the first step
	var options []ConversationOption
	if err := tx.Where("step_id = ?", firstStep.ID).Find(&options).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Append options as a numbered list to the message content
	for i, option := range options {
		messageContent += fmt.Sprintf("%d. %s\n", i+1, option.Content)
	}

	// Generate a unique ID for the message
	messageID, err := uuid.NewUUID()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create and save the message with the combined content
	firstMessage := ConversationMessage{
		ID:        messageID.String(),
		Content:   messageContent,
		SessionID: conversationSession.ID,
		CreatedBy: SYSTEM, // System indicates this is a system-generated message
		FromUser:  SYSTEM,
	}

	if err := tx.Save(&firstMessage).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (cr *conversationRepository) GetActiveSessionByUserID(userID string) (*ConversationSession, error) {
	var session ConversationSession
	if err := cr.db.Where("user_id = ? AND is_active = ?", userID, true).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active session found
		}
		return nil, err
	}

	return &session, nil
}

func (cr *conversationRepository) GetSessionByID(sessionID string, session *ConversationSession) error {
	return cr.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC") // Sort messages by created_at in ascending order
	}).First(&session, "id = ?", sessionID).Error
}

func (cr *conversationRepository) UpdateSession(session *ConversationSession) error {
	return cr.db.Save(session).Error
}

func (cr *conversationRepository) GetOptionByNumber(number string, option *ConversationOption) error {
	return cr.db.Where("number = ?", number).First(option).Error
}

func (cr *conversationRepository) GetStepByID(stepID string, step *ConversationStep) error {
	return cr.db.Preload("Options").First(step, stepID).Error
}

func (cr *conversationRepository) GetOptionByNumberAndStep(number string, stepID string, option *ConversationOption) error {
	return cr.db.Where("number = ? AND step_id = ?", number, stepID).First(option).Error
}

func (cr *conversationRepository) SaveMessage(message *ConversationMessage) error {
	return cr.db.Save(message).Error
}

func (cr *conversationRepository) EndConversationSessionBySessionID(sessionID string, endedBy string) error {
	conversationSession, err := cr.GetConversationSessionBySessionID(sessionID)
	if err != nil {
		return err
	}

	return cr.db.Model(&conversationSession).Where("id = ?", sessionID).Update("is_active", false).Error
}

func (cr *conversationRepository) SeedConversationSteps() {
	// Check if the steps have already been seeded
	var count int64
	cr.db.Model(&ConversationStep{}).Count(&count)
	if count > 0 {
		log.Println("Conversation steps already seeded")
		return
	}

	// Define the conversation steps
	steps := []ConversationStep{
		{
			ID:      "1",
			Content: "Welcome! What do you need help with today?",
		},
		{
			ID:      "2",
			Content: "Here are some common questions you might have:",
		},
		{
			ID:      "3",
			Content: "For questions about your order, check your order status under 'My Account' > 'Orders'. If you need further assistance, contact support.",
		},
		{
			ID:      "4",
			Content: "For questions about withdrawals, find withdrawal details in 'My Account' > 'Withdrawals'. Contact support if issues persist.",
		},
	}

	// Define the conversation options
	options := []ConversationOption{
		{ID: "1_1", StepID: "1", Number: "1", Content: "I have a question", NextStepID: stringPtr("2")},
		{ID: "1_2", StepID: "1", Number: "2", Content: "I need assistance from an operator/admin", NextStepID: nil},
		{ID: "2_1", StepID: "2", Number: "1", Content: "Questions about my order", NextStepID: stringPtr("3")},
		{ID: "2_2", StepID: "2", Number: "2", Content: "Questions about withdrawal", NextStepID: stringPtr("4")},
		{ID: "3_1", StepID: "3", Number: "1", Content: "I need assistance from an operator/admin", NextStepID: nil},
		{ID: "4_1", StepID: "4", Number: "1", Content: "I need assistance from an operator/admin", NextStepID: nil},
	}

	// Seed conversation steps
	for _, step := range steps {
		if err := cr.db.Create(&step).Error; err != nil {
			// log.Fatalf("Failed to seed step: %v", err)
			log.Println("Failed to seed step", err)

			// return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	// Seed conversation options
	for _, option := range options {
		if err := cr.db.Create(&option).Error; err != nil {
			// log.Fatalf("Failed to seed option: %v", err)
			log.Println("Failed to seed step", err)
		}
	}

	log.Println("Conversation steps and options seeded successfully!")
}

// Helper function to get a pointer to a string
func stringPtr(s string) *string {
	return &s
}
