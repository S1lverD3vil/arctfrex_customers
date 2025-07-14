package conversation

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"arctfrex-customers/internal/model"
)

type ConversationUsecase interface {
	GetConversationSessions() ([]ConversationSession, error)
	GetConversationSessionBySessionID(conversationSessionId string) (*ConversationSession, error)
	TakeConversationSessionBySessionID(conversationSessionId string, backOfficeUserID string) error
	GetConversationSessionsByBackofficeUserID(backofficeUserID string) ([]ConversationSession, error)
	GetUserByID(userId string) (*model.Users, error)
	GetBackofficeUserByID(userId string) (*model.BackofficeUsers, error)
	CreateSession(conversationSession *ConversationSession) error
	GetActiveSessionByUserID(userID string) (*ConversationSession, error)
	GetSessionByID(sessionID string, session *ConversationSession) error
	UpdateSession(session *ConversationSession) error // Add this
	GetOptionByNumber(number string, option *ConversationOption) error
	GetStepByID(stepID string, step *ConversationStep) error
	GetOptionByNumberAndStep(number string, stepID string, option *ConversationOption) error // Scoped lookup
	SaveMessage(message *ConversationMessage) error
	ProcessMessage(sessionID string, userID string, inputMessage string) (ConversationMessage, string, error)
	EndConversationSessionBySessionID(sessionID string, endedBy string) error
}

type conversationUsecase struct {
	conversationRepository ConversationRepository
}

func NewConversationUsecase(cr ConversationRepository) ConversationUsecase {
	return &conversationUsecase{
		conversationRepository: cr,
	}
}

func (cu *conversationUsecase) CreateSession(session *ConversationSession) error {
	return cu.conversationRepository.CreateSession(session)
}

func (cu *conversationUsecase) GetActiveSessionByUserID(userID string) (*ConversationSession, error) {
	return cu.conversationRepository.GetActiveSessionByUserID(userID)
}

func (cu *conversationUsecase) GetConversationSessions() ([]ConversationSession, error) {
	return cu.conversationRepository.GetConversationSessions()
}

func (cu *conversationUsecase) GetConversationSessionBySessionID(conversationSessionId string) (*ConversationSession, error) {
	return cu.conversationRepository.GetConversationSessionBySessionID(conversationSessionId)
}

func (cu *conversationUsecase) TakeConversationSessionBySessionID(conversationSessionId string, backOfficeUserID string) error {
	return cu.conversationRepository.TakeConversationSessionBySessionID(conversationSessionId, backOfficeUserID)
}

func (cu *conversationUsecase) GetConversationSessionsByBackofficeUserID(backofficeUserID string) ([]ConversationSession, error) {
	return cu.conversationRepository.GetConversationSessionsByBackofficeUserID(backofficeUserID)
}

func (cu *conversationUsecase) GetUserByID(userId string) (*model.Users, error) {
	return cu.conversationRepository.GetUserByID(userId)
}

func (cu *conversationUsecase) GetBackofficeUserByID(userId string) (*model.BackofficeUsers, error) {
	return cu.conversationRepository.GetBackofficeUserByID(userId)
}

func (cu *conversationUsecase) UpdateSession(session *ConversationSession) error {
	return cu.conversationRepository.UpdateSession(session)
}

func (cu *conversationUsecase) GetSessionByID(sessionID string, session *ConversationSession) error {
	return cu.conversationRepository.GetSessionByID(sessionID, session)
}

func (cu *conversationUsecase) GetOptionByNumber(number string, option *ConversationOption) error {
	return cu.conversationRepository.GetOptionByNumber(number, option)
}

func (cu *conversationUsecase) GetStepByID(stepID string, step *ConversationStep) error {
	return cu.conversationRepository.GetStepByID(stepID, step)
}

func (cu *conversationUsecase) GetOptionByNumberAndStep(number string, stepID string, option *ConversationOption) error {
	return cu.conversationRepository.GetOptionByNumberAndStep(number, stepID, option)
}

func (cu *conversationUsecase) SaveMessage(message *ConversationMessage) error {
	return cu.conversationRepository.SaveMessage(message)
}

func (cu *conversationUsecase) EndConversationSessionBySessionID(sessionID string, endedBy string) error {
	return cu.conversationRepository.EndConversationSessionBySessionID(sessionID, endedBy)
}

func (cu *conversationUsecase) ProcessMessage(sessionID string, userID string, inputMessage string) (ConversationMessage, string, error) {
	var userMessage ConversationMessage
	var systemMessage ConversationMessage
	var systemMessageContent string

	// Fetch session
	session, err := cu.GetConversationSessionBySessionID(sessionID)
	if err != nil {
		return userMessage, systemMessageContent, err
	}

	// Fetch user
	mobileUser, _ := cu.GetUserByID(userID)
	backofficeUser, _ := cu.GetBackofficeUserByID(userID)

	// Check if the user exists
	if mobileUser == nil && backofficeUser == nil {
		return userMessage, systemMessageContent, errors.New("user does not exist")
	}

	// Prepare user data
	var userId string
	var fromUser string
	if mobileUser != nil {
		userId = mobileUser.ID
		fromUser = USER
	} else if backofficeUser != nil {
		userId = backofficeUser.ID
		fromUser = OPERATOR
	}

	// Save user message
	userMessage = ConversationMessage{
		ID:          uuid.New().String(),
		Content:     inputMessage,
		CreatedByID: userId,
		SessionID:   session.ID,
		CreatedBy:   userId,
		FromUser:    fromUser,
	}
	if err := cu.SaveMessage(&userMessage); err != nil {
		return userMessage, systemMessageContent, err
	}

	// Handle current step
	if session.CurrentStepID == "" {
		return userMessage, systemMessageContent, nil
	}

	// Fetch current step
	var currentStep ConversationStep
	if err := cu.GetStepByID(session.CurrentStepID, &currentStep); err != nil {
		return userMessage, systemMessageContent, err
	}

	// Validate option
	var selectedOption ConversationOption
	if err := cu.GetOptionByNumberAndStep(inputMessage, currentStep.ID, &selectedOption); err != nil {
		return userMessage, systemMessageContent, err
	}

	// Transition logic
	if selectedOption.NextStepID != nil {
		var nextStep ConversationStep
		if err := cu.GetStepByID(*selectedOption.NextStepID, &nextStep); err != nil {
			return userMessage, systemMessageContent, err
		}

		// Generate system message content
		systemMessageContent = nextStep.Content + "\n\n"
		for i, opt := range nextStep.Options {
			systemMessageContent += fmt.Sprintf("%d. %s\n", i+1, opt.Content)
		}

		// Save system message
		systemMessage = ConversationMessage{
			ID:        uuid.New().String(),
			Content:   systemMessageContent,
			SessionID: session.ID,
			CreatedBy: SYSTEM,
			FromUser:  SYSTEM,
		}
		if err := cu.SaveMessage(&systemMessage); err != nil {
			return userMessage, systemMessageContent, err
		}

		// Update session with the next step
		session.CurrentStepID = nextStep.ID
		if err := cu.UpdateSession(session); err != nil {
			return userMessage, systemMessageContent, err
		}

		return userMessage, systemMessage.Content, nil
	}

	// End of conversation
	session.CurrentStepID = ""
	if err := cu.UpdateSession(session); err != nil {
		return userMessage, systemMessageContent, err
	}

	// Save end-of-conversation system message
	systemMessageContent = "Operator/Admin will reach you, please wait for a while..."
	systemMessage = ConversationMessage{
		ID:        uuid.New().String(),
		Content:   systemMessageContent,
		SessionID: session.ID,
		CreatedBy: SYSTEM,
		FromUser:  SYSTEM,
	}
	if err := cu.SaveMessage(&systemMessage); err != nil {
		return userMessage, systemMessageContent, err
	}

	return userMessage, systemMessage.Content, nil
}
