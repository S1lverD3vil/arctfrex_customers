package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/dto"
	"arctfrex-customers/internal/middleware"
	"arctfrex-customers/internal/usecase"
)

type workflowApproverHandler struct {
	jwtMiddleware           *middleware.JWTMiddleware
	workflowapproverUsecase usecase.WorkflowApproverUsecase
}

func NewWorkflowApproverHandler(
	engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	workflowapproverUsecase usecase.WorkflowApproverUsecase,
) *workflowApproverHandler {
	handler := &workflowApproverHandler{
		jwtMiddleware:           jmw,
		workflowapproverUsecase: workflowapproverUsecase,
	}

	protectedRoutes := engine.Group("backoffice/workflow-approver")
	{
		protectedRoutes.Use(jmw.ValidateToken())
		protectedRoutes.POST("/approve-reject", handler.ApproveReject)
	}

	return handler
}

func (wa *workflowApproverHandler) ApproveReject(c *gin.Context) {
	var approverReject dto.ApproveRejectRequest

	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
	}
	approverReject.UserID = userId

	if err := c.ShouldBindJSON(&approverReject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := approverReject.ValidationRequest(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := wa.workflowapproverUsecase.ApproveRejectWorkflow(approverReject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ApiResponse{ApiResponse: base.ApiResponse{Message: "success", Data: response}})
}
