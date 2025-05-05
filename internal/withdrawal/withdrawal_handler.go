package withdrawal

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/middleware"
)

type withdrawalHandler struct {
	jwtMiddleware     *middleware.JWTMiddleware
	withdrawalUsecase WithdrawalUsecase
}

func NewWithdrawalHandler(
	engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	du WithdrawalUsecase,
) *withdrawalHandler {
	handler := &withdrawalHandler{
		jwtMiddleware:     jmw,
		withdrawalUsecase: du,
	}
	unprotectedRoutesBackOffice := engine.Group("backoffice/withdrawal")
	protectedRoutes := engine.Group("/withdrawal")

	unprotectedRoutesBackOffice.GET("/pending", handler.BackOfficePending)
	unprotectedRoutesBackOffice.POST("/pending", handler.BackOfficePending)
	unprotectedRoutesBackOffice.GET("/pending/:withdrawalid", handler.BackOfficePendingDetail)
	unprotectedRoutesBackOffice.POST("/pending/:withdrawalid", handler.BackOfficePendingDetail)
	unprotectedRoutesBackOffice.POST("/pending/approval", handler.BackOfficePendingApproval)

	// Get pending deposit SPA
	unprotectedRoutesBackOffice.GET("/pending/spa/:menutype", handler.BackOfficePendingSPA)

	// Get pending deposit multi
	unprotectedRoutesBackOffice.GET("/pending/multi/:menutype", handler.BackOfficePendingMulti)

	protectedRoutes.Use(jmw.ValidateToken())
	{
		protectedRoutes.POST("/submit", handler.Submit)
		protectedRoutes.GET("/:accountId", handler.WithdrawalByAccountId)
		protectedRoutes.GET("/pending/:accountId", handler.Pending)
		protectedRoutes.GET("/detail/:withdrawalId", handler.Detail)
	}

	return handler
}

func (wh *withdrawalHandler) Submit(c *gin.Context) {
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
	var withdrawal *Withdrawal
	if err := c.ShouldBindJSON(&withdrawal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	withdrawal.UserID = userId
	withdrawalId, err := wh.withdrawalUsecase.Submit(withdrawal)

	if err != nil {
		if err.Error() == "insufficient balance" {
			c.JSON(412, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"withdrawalid": withdrawalId})
}

func (wh *withdrawalHandler) Pending(c *gin.Context) {
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

	if err := wh.withdrawalUsecase.Pending(userId, c.Param("accountId")); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

func (wh *withdrawalHandler) WithdrawalByAccountId(c *gin.Context) {
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

	withdrawals, err := wh.withdrawalUsecase.WithdrawalByAccountId(userId, c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		WithdrawalApiResponse{base.ApiResponse{Message: "success", Data: withdrawals}},
	)
}

func (wh *withdrawalHandler) Detail(c *gin.Context) {
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
	pendingDetail, err := wh.withdrawalUsecase.Detail(userId, c.Param("withdrawalId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pendingDetail)
}

func (wh *withdrawalHandler) BackOfficePending(c *gin.Context) {

	withdrawals, err := wh.withdrawalUsecase.BackOfficePending()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		WithdrawalApiResponse{base.ApiResponse{Message: "success", Data: withdrawals}},
	)
}

func (wh *withdrawalHandler) BackOfficePendingDetail(c *gin.Context) {
	pendingApprovalDetail, err := wh.withdrawalUsecase.BackOfficePendingDetail(c.Param("withdrawalid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, WithdrawalApiResponse{base.ApiResponse{
		Message: "success",
		Data:    pendingApprovalDetail,
	}})
}

func (wh *withdrawalHandler) BackOfficePendingApproval(c *gin.Context) {
	var backOfficeApproval *BackOfficePendingApprovalRequest
	if err := c.ShouldBindJSON(&backOfficeApproval); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := wh.withdrawalUsecase.BackOfficePendingApproval(*backOfficeApproval); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		WithdrawalApiResponse{base.ApiResponse{Message: "success"}},
	)
}

func (wh *withdrawalHandler) BackOfficePendingSPA(c *gin.Context) {
	menutype := c.Param("menutype")
	request := WithdrawalBackOfficeParam{
		Menutype: menutype,
	}

	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.Pagination.Norm()

	withdrawals, err := wh.withdrawalUsecase.BackOfficePendingSPA(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		base.ApiPaginatedResponse{Message: "success", Data: withdrawals.Data, Pagination: *withdrawals.Pagination},
	)
}

func (wh *withdrawalHandler) BackOfficePendingMulti(c *gin.Context) {
	menutype := c.Param("menutype")
	request := WithdrawalBackOfficeParam{
		Menutype: menutype,
	}

	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.Pagination.Norm()

	withdrawals, err := wh.withdrawalUsecase.BackOfficePendingMulti(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		base.ApiPaginatedResponse{Message: "success", Data: withdrawals.Data, Pagination: *withdrawals.Pagination},
	)
}
