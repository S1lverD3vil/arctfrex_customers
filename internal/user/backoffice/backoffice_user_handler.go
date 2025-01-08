package user

import (
	"arctfrex-customers/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type backofficeUserHandler struct {
	jwtMiddleware         *middleware.JWTMiddleware
	backofficeUserUsecase BackofficeUserUsecase
}

func NewBackofficeHandler(engine *gin.Engine, jmw *middleware.JWTMiddleware, buu *backofficeUserUsecase) *backofficeUserHandler {
	handler := &backofficeUserHandler{
		jwtMiddleware:         jmw,
		backofficeUserUsecase: buu,
	}

	unprotectedRoutes := engine.Group("/backoffice/users")
	unprotectedRoutes.POST("/register", handler.Register)
	unprotectedRoutes.POST("login/session", handler.LoginSession)
	unprotectedRoutes.GET("/all", handler.All)

	return handler
}

func (buh *backofficeUserHandler) Register(c *gin.Context) {
	var backofficeUser BackofficeUsers
	if err := c.ShouldBindJSON(&backofficeUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := buh.backofficeUserUsecase.Register(&backofficeUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (buh *backofficeUserHandler) LoginSession(c *gin.Context) {
	var backofficeUser BackofficeUsers
	if err := c.ShouldBindJSON(&backofficeUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	backofficeUserLoginResponse, err := buh.backofficeUserUsecase.LoginSession(&backofficeUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, backofficeUserLoginResponse)
}

func (buh *backofficeUserHandler) All(c *gin.Context) {
	users, err := buh.backofficeUserUsecase.All()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		BackofficeUserApiResponse{Message: "success", Data: users},
	)
}
