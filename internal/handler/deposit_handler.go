package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/middleware"
	"arctfrex-customers/internal/model"
	"arctfrex-customers/internal/usecase"
)

type depositHandler struct {
	jwtMiddleware  *middleware.JWTMiddleware
	depositUsecase usecase.DepositUsecase
}

func NewDepositHandler(
	engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	du usecase.DepositUsecase,
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

	// Get pending deposit SPA
	unprotectedRoutesBackOffice.GET("/pending/spa/:menutype", handler.BackOfficePendingSPA)
	unprotectedRoutesBackOffice.GET("/credit/spa/:menutype", handler.BackOfficeCreditSPA)

	// Get pending deposit multi
	unprotectedRoutesBackOffice.GET("/pending/multi/:menutype", handler.BackOfficePendingMulti)
	unprotectedRoutesBackOffice.GET("/credit/multi/:menutype", handler.BackOfficeCreaditMulti)

	unprotectedRoutesBackOffice.PATCH(":depositid/credit-type", handler.BackOfficeUpdateCreditType)

	// Get credit deposit SPA
	unprotectedRoutesBackOffice.GET("/credit/spa/:menutype/:depositid", handler.BackOfficeCreditSPADetail)

	// Get credit deposit SPA
	unprotectedRoutesBackOffice.GET("/credit/multi/:menutype/:depositid", handler.BackOfficeCreditMultiDetail)

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
	var deposit *model.Deposit
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
		model.ApiResponse{ApiResponse: base.ApiResponse{Message: "success", Data: deposits}},
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
		model.ApiResponse{ApiResponse: base.ApiResponse{Message: "success", Data: deposits}},
	)
}

func (dh *depositHandler) BackOfficePendingDetail(c *gin.Context) {
	pendingDetail, err := dh.depositUsecase.BackOfficePendingDetail(c.Param("depositid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ApiResponse{ApiResponse: base.ApiResponse{
		Message: "success",
		Data:    pendingDetail,
	}})
}

func (dh *depositHandler) BackOfficePendingApproval(c *gin.Context) {
	var backOfficeApproval *model.BackOfficePendingApprovalRequest
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
		model.ApiResponse{ApiResponse: base.ApiResponse{Message: "success"}},
	)
}

func (dh *depositHandler) BackOfficePendingSPA(c *gin.Context) {
	menutype := c.Param("menutype")
	request := model.DepositBackOfficeParam{
		Menutype: menutype,
	}

	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.Pagination.Norm()

	deposits, err := dh.depositUsecase.BackOfficePendingSPA(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		model.ApiPaginatedResponse{ApiPaginatedResponse: base.ApiPaginatedResponse{Message: "success", Data: deposits.Data, Pagination: *deposits.Pagination}},
	)
}

func (dh *depositHandler) BackOfficePendingMulti(c *gin.Context) {
	menutype := c.Param("menutype")
	request := model.DepositBackOfficeParam{
		Menutype: menutype,
	}

	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.Pagination.Norm()

	deposits, err := dh.depositUsecase.BackOfficePendingMulti(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		model.ApiPaginatedResponse{ApiPaginatedResponse: base.ApiPaginatedResponse{Message: "success", Data: deposits.Data, Pagination: *deposits.Pagination}},
	)
}

func (dh *depositHandler) BackOfficeCreditSPA(c *gin.Context) {
	menutype := c.Param("menutype")
	request := model.CreditBackOfficeParam{
		Menutype: menutype,
	}

	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.Pagination.Norm()

	deposits, err := dh.depositUsecase.BackOfficeCreditSPA(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		model.ApiPaginatedResponse{ApiPaginatedResponse: base.ApiPaginatedResponse{Message: "success", Data: deposits.Data, Pagination: *deposits.Pagination}},
	)
}

func (dh *depositHandler) BackOfficeCreaditMulti(c *gin.Context) {
	menutype := c.Param("menutype")
	request := model.CreditBackOfficeParam{
		Menutype: menutype,
	}

	err := c.ShouldBindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.Pagination.Norm()

	deposits, err := dh.depositUsecase.BackOfficeCreditMulti(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		model.ApiPaginatedResponse{ApiPaginatedResponse: base.ApiPaginatedResponse{Message: "success", Data: deposits.Data, Pagination: *deposits.Pagination}},
	)
}

func (dh *depositHandler) BackOfficeUpdateCreditType(c *gin.Context) {
	var backOfficeUpdateCreditTypeRequest *model.BackOfficeUpdateCreditTypeRequest

	if err := c.ShouldBindJSON(&backOfficeUpdateCreditTypeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	backOfficeUpdateCreditTypeRequest.Depositid = c.Param("depositid")

	err := backOfficeUpdateCreditTypeRequest.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dh.depositUsecase.BackOfficeUpdateCreditType(*backOfficeUpdateCreditTypeRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		model.ApiResponse{ApiResponse: base.ApiResponse{Message: "success", Data: backOfficeUpdateCreditTypeRequest}},
	)
}

func (dh *depositHandler) BackOfficeCreditSPADetail(c *gin.Context) {
	param := model.CreditBackOfficeDetailParam{
		Menutype:  c.Param("menutype"),
		DepositID: c.Param("depositid"),
	}

	if err := param.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creditDetail, err := dh.depositUsecase.BackOfficeCreditSPADetail(c.Request.Context(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ApiResponse{ApiResponse: base.ApiResponse{
		Message: "success",
		Data:    creditDetail,
	}})
}

func (dh *depositHandler) BackOfficeCreditMultiDetail(c *gin.Context) {
	param := model.CreditBackOfficeDetailParam{
		Menutype:  c.Param("menutype"),
		DepositID: c.Param("depositid"),
	}

	if err := param.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creditDetail, err := dh.depositUsecase.BackOfficeCreditMultiDetail(c.Request.Context(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.ApiResponse{ApiResponse: base.ApiResponse{
		Message: "success",
		Data:    creditDetail,
	}})
}
