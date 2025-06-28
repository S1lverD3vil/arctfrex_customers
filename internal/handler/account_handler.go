package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/middleware"
	"arctfrex-customers/internal/model"
	"arctfrex-customers/internal/usecase"
)

type accountHandler struct {
	jwtMiddleware  *middleware.JWTMiddleware
	accountUsecase usecase.AccountUsecase
}

func NewAccountHandler(
	engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	au usecase.AccountUsecase,
) *accountHandler {
	handler := &accountHandler{
		jwtMiddleware:  jmw,
		accountUsecase: au,
	}

	unprotectedRoutesBackOffice := engine.Group("/backoffice/account")
	protectedRoutes := engine.Group("/account")

	unprotectedRoutesBackOffice.GET("/all", handler.BackOfficeAll)
	unprotectedRoutesBackOffice.GET("/pending", handler.BackOfficePending)
	// unprotectedRoutesBackOffice.POST("/pending", handler.BackOfficePending)
	unprotectedRoutesBackOffice.POST("/pending/approval", handler.BackOfficePendingApproval)
	protectedRoutes.Use(jmw.ValidateToken())
	{
		protectedRoutes.POST("/submit", handler.Submit)
		protectedRoutes.GET("/pending", handler.Pending)
		protectedRoutes.GET("/pending/check", handler.PendingCheck)
		protectedRoutes.GET("/", handler.GetAccounts)
		protectedRoutes.PATCH("/topup", handler.TopUpAccount)
	}

	return handler
}

func (ah *accountHandler) Submit(c *gin.Context) {
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

	if err := ah.accountUsecase.Submit(&model.Account{UserID: userId}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (ah *accountHandler) Pending(c *gin.Context) {
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
	accounts, err := ah.accountUsecase.Pending(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		model.AccountApiResponse{Message: "success", Data: accounts},
	)
}

func (ah *accountHandler) PendingCheck(c *gin.Context) {
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
	pendingAccounts, err := ah.accountUsecase.PendingCheck(userId)
	if pendingAccounts != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	fmt.Printf("Pending Account: %+v\n", pendingAccounts)

	c.Status(http.StatusOK)
}

func (ah *accountHandler) GetAccounts(c *gin.Context) {
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
	accounts, err := ah.accountUsecase.GetAccounts(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		model.AccountApiResponse{Message: "success", Data: accounts},
	)
}

func (ah *accountHandler) TopUpAccount(c *gin.Context) {
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

	var topUpAccount *model.TopUpAccount
	if err := c.ShouldBindJSON(&topUpAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topUpAccount.UserID = userId

	if err := ah.accountUsecase.TopUpAccount(*topUpAccount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (ah *accountHandler) BackOfficeAll(c *gin.Context) {
	request := model.BackOfficeAllAccountRequest{}
	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.Pagination.Norm()

	response, err := ah.accountUsecase.BackOfficeAll(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		base.ApiPaginatedResponse{Message: "success", Data: response.Data, Pagination: *response.Pagination},
	)
}

func (ah *accountHandler) BackOfficePending(c *gin.Context) {
	request := model.BackOfficePendingAccountRequest{}
	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.Pagination.Norm()

	response, err := ah.accountUsecase.BackOfficePending(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		base.ApiPaginatedResponse{Message: "success", Data: response.Data, Pagination: *response.Pagination},
	)
}

func (ah *accountHandler) BackOfficePendingApproval(c *gin.Context) {
	var backOfficeApproval *model.BackOfficePendingAccountApprovalRequest
	if err := c.ShouldBindJSON(&backOfficeApproval); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ah.accountUsecase.BackOfficePendingApproval(*backOfficeApproval); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		model.AccountApiResponse{Message: "success"},
	)
}
