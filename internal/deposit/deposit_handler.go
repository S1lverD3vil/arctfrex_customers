package deposit

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type depositHandler struct {
	jwtMiddleware  *middleware.JWTMiddleware
	depositUsecase DepositUsecase
}

func NewDepositHandler(
	engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	du DepositUsecase,
) *depositHandler {
	handler := &depositHandler{
		jwtMiddleware:  jmw,
		depositUsecase: du,
	}
	unprotectedRoutesBackOffice := engine.Group("backoffice/deposit")
	protectedRoutes := engine.Group("/deposit")

	unprotectedRoutesBackOffice.GET("/pending", handler.BackOfficePending)
	unprotectedRoutesBackOffice.POST("/pending", handler.BackOfficePending)
	unprotectedRoutesBackOffice.GET("/pending/:depositid", handler.BackOfficePendingDetail)
	unprotectedRoutesBackOffice.POST("/pending/:depositid", handler.BackOfficePendingDetail)
	unprotectedRoutesBackOffice.POST("/pending/approval", handler.BackOfficePendingApproval)

	protectedRoutes.Use(jmw.ValidateToken())
	{
		protectedRoutes.POST("/submit", handler.Submit)
		protectedRoutes.GET("/pending/:accountId", handler.Pending)
		protectedRoutes.GET("/:accountId", handler.DepositByAccountId)
		protectedRoutes.GET("/detail/:depositId", handler.Detail)
	}

	return handler
}

func (dh *depositHandler) Submit(c *gin.Context) {
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
	var deposit *Deposit
	if err := c.ShouldBindJSON(&deposit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deposit.UserID = userId

	depositId, err := dh.depositUsecase.Submit(deposit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"depositid": depositId})
}

func (dh *depositHandler) Pending(c *gin.Context) {
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

	if err := dh.depositUsecase.Pending(userId, c.Param("accountId")); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

func (dh *depositHandler) DepositByAccountId(c *gin.Context) {
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

	deposits, err := dh.depositUsecase.DepositByAccountId(userId, c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		DepositApiResponse{base.ApiResponse{Message: "success", Data: deposits}},
	)
}

func (dh *depositHandler) Detail(c *gin.Context) {
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
	pendingDetail, err := dh.depositUsecase.Detail(userId, c.Param("depositId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pendingDetail)
}

func (dh *depositHandler) BackOfficePending(c *gin.Context) {

	deposits, err := dh.depositUsecase.BackOfficePending()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		DepositApiResponse{base.ApiResponse{Message: "success", Data: deposits}},
	)
}

func (dh *depositHandler) BackOfficePendingDetail(c *gin.Context) {
	pendingDetail, err := dh.depositUsecase.BackOfficePendingDetail(c.Param("depositid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, DepositApiResponse{base.ApiResponse{
		Message: "success",
		Data:    pendingDetail,
	}})
}

func (dh *depositHandler) BackOfficePendingApproval(c *gin.Context) {
	var backOfficeApproval *BackOfficePendingApprovalRequest
	if err := c.ShouldBindJSON(&backOfficeApproval); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dh.depositUsecase.BackOfficePendingApproval(*backOfficeApproval); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		DepositApiResponse{base.ApiResponse{Message: "success"}},
	)
}
