package email

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	EmailUseCase EmailUseCase
}

func NewEmailHandler(router *gin.Engine, eu EmailUseCase) {
	handler := &EmailHandler{
		EmailUseCase: eu,
	}
	router.POST("/email/send", handler.SendEmail)
}

func (eh *EmailHandler) SendEmail(c *gin.Context) {
	var email Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := eh.EmailUseCase.SendEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
