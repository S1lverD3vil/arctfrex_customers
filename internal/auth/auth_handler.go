package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase AuthUsecase
}

func NewAuthHandler(engine *gin.Engine, au AuthUsecase) *AuthHandler {
	handler := &AuthHandler{authUsecase: au}
	engine.POST("/token", handler.Token)

	return handler
}

func (h *AuthHandler) Token(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	token, _, err := h.authUsecase.Token(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
