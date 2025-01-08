package inbox

import (
	"arctfrex-customers/internal/middleware"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type inboxHandler struct {
	jwtMiddleware *middleware.JWTMiddleware
	inboxUsecase  InboxUseCase
}

func NewInboxHandler(
	// router *gin.RouterGroup,
	engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	iu InboxUseCase,
) *inboxHandler {
	handler := &inboxHandler{
		jwtMiddleware: jmw,
		inboxUsecase:  iu,
	}

	unprotectedRoutes, protectedRoutes := engine.Group("/inbox"), engine.Group("/inbox")
	unprotectedRoutes.POST("/add", handler.CreateInbox)
	protectedRoutes.Use(jmw.ValidateToken())
	{
		// protectedRoutes.POST("/add", handler.CreateInbox)
		protectedRoutes.GET("/all", handler.GetUserInboxes)
		// protectedRoutes.GET("/inboxes/:receiver", handler.GetUserInboxes)
		protectedRoutes.PATCH("/read/:inboxid", handler.MarkAsRead)
		protectedRoutes.DELETE("/:inboxid", handler.DeleteInbox)
	}

	return handler

}

func (ih *inboxHandler) CreateInbox(c *gin.Context) {
	// // Retrieve the userID from context
	// userID, exists := c.Get("userID")
	// if !exists {
	// 	log.Println("userID not found in context")
	// }

	// // Convert userID to string
	// userId, ok := userID.(string)
	// if !ok {
	// 	log.Println("userID is not of type string")
	// 	return
	// }

	var inbox Inbox
	if err := c.ShouldBindJSON(&inbox); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// inbox.UserID = userId
	if err := ih.inboxUsecase.CreateInbox(&inbox); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inbox"})
		return
	}

	c.JSON(http.StatusCreated, inbox)
}

func (ih *inboxHandler) GetUserInboxes(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
		return
	}
	// receiver := c.Param("receiver")
	inboxes, err := ih.inboxUsecase.GetUserInboxes(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inboxes"})
		return
	}

	c.JSON(http.StatusOK, inboxes)
}

func (ih *inboxHandler) MarkAsRead(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
		return
	}

	id, _ := strconv.Atoi(c.Param("inboxid"))

	if err := ih.inboxUsecase.MarkInboxAsRead(userId, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark inbox as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inbox marked as read"})
}

func (ih *inboxHandler) DeleteInbox(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
		return
	}

	id, _ := strconv.Atoi(c.Param("inboxid"))

	if err := ih.inboxUsecase.DeleteInbox(userId, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete inbox"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inbox deleted"})
}
