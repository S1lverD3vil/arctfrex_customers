package conversation

import (
	"arctfrex-customers/internal/auth"
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/middleware"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var connectionManager = NewConnectionManager()

type conversationHandler struct {
	jwtMiddleware       *middleware.JWTMiddleware
	conversationUsecase ConversationUsecase
	tokenService        auth.TokenService
}

func NewConversationHandler(
	engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	cu ConversationUsecase,
	ts auth.TokenService,
) *conversationHandler {
	handler := &conversationHandler{
		jwtMiddleware:       jmw,
		conversationUsecase: cu,
		tokenService:        ts,
	}

	unprotectedRoutes, protectedRoutes, unprotectedRoutesBackOffice, protectedRoutesBackOffice := engine.Group("/conversations"), engine.Group("/conversations"), engine.Group("/backoffice/conversations"), engine.Group("/backoffice/conversations")

	unprotectedRoutesBackOffice.GET("/sessions", handler.BackOfficeConversationSessions)
	unprotectedRoutesBackOffice.GET("/sessions/:session_id", handler.BackOfficeConversationSessionByConversationSessionId)

	protectedRoutesBackOffice.Use(jmw.ValidateToken())
	{
		protectedRoutesBackOffice.POST("/sessions/:session_id/take", handler.BackOfficeConversationSessionOnTake)
		protectedRoutesBackOffice.POST("/sessions/:session_id/end", handler.EndSession)
	}

	unprotectedRoutes.GET("/sessions/:session_id/messages/ws", handler.WsSendMessage) // Websocket for mobile users and admins/operators
	protectedRoutes.Use(jmw.ValidateToken())
	{
		protectedRoutes.POST("/sessions/start", handler.Start) // Only for mobile users
		protectedRoutes.GET("/sessions/active", handler.GetActiveSession)
		protectedRoutes.POST("/sessions/:session_id/select-option", handler.SelectOption)
		protectedRoutes.POST("/sessions/:session_id/messages", handler.PostMessage)
		protectedRoutes.POST("/sessions/:session_id/end", handler.EndSession)
	}

	return handler
}

func (ch *conversationHandler) BackOfficeConversationSessions(c *gin.Context) {
	backofficeUserID, hasTypeInQuery := c.GetQuery("backoffice_user_id")

	// If the backoffice_user_id is in the query, fetch the sessions by backoffice user ID
	if hasTypeInQuery {
		sessions, err := ch.conversationUsecase.GetConversationSessionsByBackofficeUserID(backofficeUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
			return
		}

		c.JSON(
			http.StatusOK,
			base.ApiResponse{Message: "success", Data: sessions},
		)
		return
	}

	sessions, err := ch.conversationUsecase.GetConversationSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		base.ApiResponse{Message: "success", Data: sessions},
	)
}

func (ch *conversationHandler) BackOfficeConversationSessionByConversationSessionId(c *gin.Context) {
	sessionId := c.Param("session_id")

	session, err := ch.conversationUsecase.GetConversationSessionBySessionID(sessionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		base.ApiResponse{Message: "success", Data: session},
	)
}

func (ch *conversationHandler) BackOfficeConversationSessionOnTake(c *gin.Context) {
	sessionID := c.Param("session_id")

	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
	}

	if err := ch.conversationUsecase.TakeConversationSessionBySessionID(sessionID, userId); err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, base.ApiResponse{Message: "Conversation session taken successfully"})
}

func (ch *conversationHandler) Start(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
	}

	// Check if the user already has an active session
	activeSession, err := ch.conversationUsecase.GetActiveSessionByUserID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: "Error checking active session"})
		return
	}

	if activeSession != nil {
		// Return the existing active session
		c.JSON(http.StatusConflict, base.ApiResponse{
			Message: "Active session already exists",
			Data:    gin.H{"session_id": activeSession.ID},
		})
		return
	}

	// Generate a new message ID
	sessionID, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
	}

	user, err := ch.conversationUsecase.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
		return
	}

	// Create the conversation session
	conversationSession := ConversationSession{
		ID:     sessionID.String(),
		UserID: user.ID,
		BaseModel: base.BaseModel{
			IsActive:  true,
			CreatedBy: user.ID,
		},
	}

	if err := ch.conversationUsecase.CreateSession(&conversationSession); err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
		return
	}

	// Return the created session including the first step message
	c.JSON(http.StatusOK, base.ApiResponse{Message: "success", Data: conversationSession})
}

func (ch *conversationHandler) GetActiveSession(c *gin.Context) {
	// Retrieve userID from the token
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
		c.JSON(http.StatusUnauthorized, base.ApiResponse{Message: "Unauthorized"})
		return
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: "Invalid userID"})
		return
	}

	// Fetch active session for the user
	activeSession, err := ch.conversationUsecase.GetActiveSessionByUserID(userId)
	if err != nil {
		log.Println("Error fetching active session:", err)
		c.JSON(http.StatusInternalServerError, base.ApiResponse{
			Message: err.Error(),
		})
		return
	}

	// No active session found
	if activeSession == nil {
		c.JSON(http.StatusNotFound, base.ApiResponse{Message: "No active session found"})
		return
	}

	// Fetch the current session
	session, err := ch.conversationUsecase.GetConversationSessionBySessionID(activeSession.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
		return
	}

	// Return the active session
	c.JSON(http.StatusOK, base.ApiResponse{
		Message: "success",
		Data:    session,
	})
}

func (ch *conversationHandler) SelectOption(c *gin.Context) {
	// Get the session ID from the URL path
	sessionID := c.Param("session_id")

	// Bind the selected option from the request body
	var input struct {
		Number string `json:"number"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, base.ApiResponse{Message: err.Error()})
		return
	}

	// Fetch the current session
	session, err := ch.conversationUsecase.GetConversationSessionBySessionID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
		return
	}

	// Fetch the selected option
	var option ConversationOption
	if err := ch.conversationUsecase.GetOptionByNumber(input.Number, &option); err != nil {
		c.JSON(http.StatusNotFound, base.ApiResponse{Message: "Option not found"})
		return
	}

	// Check if there's a next step
	if option.NextStepID != nil {
		// Fetch the next conversation step
		var nextStep ConversationStep
		if err := ch.conversationUsecase.GetStepByID(*option.NextStepID, &nextStep); err != nil {
			c.JSON(http.StatusNotFound, base.ApiResponse{Message: "Next step not found"})
			return
		}

		userMessageId, err := uuid.NewUUID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
			return
		}

		// Create a new message in the session for the selected option
		userMessage := ConversationMessage{
			ID:        userMessageId.String(),
			Content:   option.Content, // The content of the selected option
			SessionID: session.ID,
			CreatedBy: session.UserID,
			FromUser:  USER, // Message from the user
		}
		if err := ch.conversationUsecase.SaveMessage(&userMessage); err != nil {
			c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
			return
		}

		nextMessageId, err := uuid.NewUUID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
			return
		}

		// Add a new message for the next step (system-generated message)
		nextMessage := ConversationMessage{
			ID:        nextMessageId.String(),
			Content:   nextStep.Content, // Content of the next step
			SessionID: session.ID,
			CreatedBy: "",
			FromUser:  SYSTEM, // Message from the system
		}
		if err := ch.conversationUsecase.SaveMessage(&nextMessage); err != nil {
			c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
			return
		}

		// Send the next step back to the user
		c.JSON(http.StatusOK, nextStep)
	} else {
		// No next step, end the conversation
		c.JSON(http.StatusOK, base.ApiResponse{Message: "Success", Data: gin.H{"message": option.Content, "next_step": false, "last_step": true}})
	}
}

func (ch *conversationHandler) PostMessage(c *gin.Context) {
	sessionID := c.Param("session_id")
	var input struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, base.ApiResponse{Message: err.Error()})
		return
	}

	// Fetch the current session
	session, err := ch.conversationUsecase.GetConversationSessionBySessionID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
		return
	}

	// Log user message
	userMessage := ConversationMessage{
		ID:        uuid.New().String(),
		Content:   input.Message,
		SessionID: session.ID,
		CreatedBy: session.UserID,
		FromUser:  USER,
	}
	if err := ch.conversationUsecase.SaveMessage(&userMessage); err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
		return
	}

	// Handle case where conversation has no more steps
	if session.CurrentStepID == "" {
		c.JSON(http.StatusOK, base.ApiResponse{Message: "Success", Data: gin.H{
			"message":    "",
			"next_step":  false,
			"last_step":  true,
			"user_input": input.Message,
		}})
		return
	}

	// Get the current step
	var currentStep ConversationStep
	if err := ch.conversationUsecase.GetStepByID(session.CurrentStepID, &currentStep); err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: "Current step not found"})
		return
	}

	// Validate user input as a valid option
	var selectedOption ConversationOption
	if err := ch.conversationUsecase.GetOptionByNumberAndStep(input.Message, currentStep.ID, &selectedOption); err != nil {
		c.JSON(http.StatusBadRequest, base.ApiResponse{Message: "Invalid option"})
		return
	}

	// If selected option has no next step, end the conversation
	if selectedOption.NextStepID == nil {
		// Log a final system message
		systemMessage := ConversationMessage{
			ID:        uuid.New().String(),
			Content:   selectedOption.Content,
			SessionID: session.ID,
			CreatedBy: SYSTEM,
			FromUser:  SYSTEM,
		}
		if err := ch.conversationUsecase.SaveMessage(&systemMessage); err != nil {
			c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
			return
		}

		// Mark session as ended by clearing CurrentStepID
		session.CurrentStepID = ""
		if err := ch.conversationUsecase.UpdateSession(session); err != nil {
			c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, base.ApiResponse{Message: "Success", Data: gin.H{
			"message":   systemMessage.Content,
			"next_step": false,
			"last_step": true,
		}})
		return
	}

	// Proceed to the next step if it exists
	var nextStep ConversationStep
	if err := ch.conversationUsecase.GetStepByID(*selectedOption.NextStepID, &nextStep); err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: "Next step not found"})
		return
	}

	// Log system message for the next step
	systemMessageContent := nextStep.Content + "\n\n"
	for i, opt := range nextStep.Options {
		systemMessageContent += fmt.Sprintf("%d. %s\n", i+1, opt.Content)
	}

	systemMessage := ConversationMessage{
		ID:        uuid.New().String(),
		Content:   systemMessageContent,
		SessionID: session.ID,
		CreatedBy: SYSTEM,
		FromUser:  SYSTEM,
	}
	if err := ch.conversationUsecase.SaveMessage(&systemMessage); err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
		return
	}

	// Update session with the next step
	session.CurrentStepID = nextStep.ID
	if err := ch.conversationUsecase.UpdateSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, base.ApiResponse{Message: "Success", Data: gin.H{
		"message":    systemMessageContent,
		"next_step":  true,
		"step_id":    nextStep.ID,
		"options":    nextStep.Options,
		"user_input": input.Message,
	}})
}

func (ch *conversationHandler) EndSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
	}

	if err := ch.conversationUsecase.EndConversationSessionBySessionID(sessionID, userId); err != nil {
		c.JSON(http.StatusInternalServerError, base.ApiResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, base.ApiResponse{Message: "Conversation session ended successfully"})
}

func (ch *conversationHandler) WsSendMessage(c *gin.Context) {
	sessionID := c.Param("session_id")

	accessToken, exists := c.GetQuery("access_token")
	if !exists {
		c.JSON(http.StatusBadRequest, base.ApiResponse{Message: "Access token not found in query"})
		return
	}

	userID, err := ch.tokenService.ValidateToken(accessToken)
	if err != nil {
		log.Println("Invalid user ID")
		c.JSON(http.StatusUnauthorized, base.ApiResponse{Message: "User is invalid"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer func() {
		connectionManager.RemoveConnection(sessionID, conn)
		conn.Close()
	}()

	connectionManager.AddConnection(sessionID, conn)

	for {
		// Read from WebSocket
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage error:", err)
			return
		}

		var input struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(message, &input); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}

		proceedMessage, systemResponse, _ := ch.conversationUsecase.ProcessMessage(sessionID, userID, input.Message)

		// Broadcast the user's message
		connectionManager.Broadcast(sessionID, proceedMessage)

		// Broadcast the system-generated response
		if systemResponse != "" {
			connectionManager.Broadcast(sessionID, ConversationMessage{
				ID:        uuid.New().String(),
				Content:   systemResponse,
				SessionID: sessionID,
				CreatedBy: SYSTEM,
				FromUser:  SYSTEM,
			})
		}
	}
}
