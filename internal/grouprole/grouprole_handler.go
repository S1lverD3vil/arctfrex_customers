package grouprole

import (
	"arctfrex-customers/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type groupRoleHandler struct {
	jwtMiddleware    *middleware.JWTMiddleware
	groupRoleUsecase GroupRoleUseCase
}

func NewGroupRoleHandler(engine *gin.Engine, jmw *middleware.JWTMiddleware, gru *groupRoleUseCase) *groupRoleHandler {
	handler := &groupRoleHandler{
		jwtMiddleware:    jmw,
		groupRoleUsecase: gru,
	}

	unprotectedRoutes := engine.Group("/backoffice/group-roles")
	{
		unprotectedRoutes.GET("/all", handler.All)
		unprotectedRoutes.POST("", handler.Create)
		unprotectedRoutes.PUT("/:id", handler.Update)
		unprotectedRoutes.DELETE("/:id", handler.Delete)
		unprotectedRoutes.GET("/:id", handler.GetByID)
	}

	return handler
}

func (rh *groupRoleHandler) All(c *gin.Context) {
	roles, err := rh.groupRoleUsecase.All()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, GroupRoleApiResponse{
		Message: "succes", Data: roles,
	})
}

func (rh *groupRoleHandler) Create(c *gin.Context) {
	var groupRole CreateUserDTO
	if err := c.ShouldBindJSON(&groupRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newRole := GroupRole{
		ID:   groupRole.ID,
		Name: groupRole.Name,
	}

	if err := rh.groupRoleUsecase.Create(&newRole); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, GroupRoleApiResponse{
		Message: "succes",
		Data:    groupRole,
	})
}

func (rh *groupRoleHandler) Update(c *gin.Context) {
	roleID := c.Param("id")
	var groupRole GroupRole
	if err := c.ShouldBindJSON(&groupRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupRole.ID = roleID
	if err := rh.groupRoleUsecase.Update(&groupRole); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GroupRoleApiResponse{
		Message: "success",
		Data:    groupRole,
	})
}

func (rh *groupRoleHandler) Delete(c *gin.Context) {
	roleID := c.Param("id")
	if err := rh.groupRoleUsecase.Delete(roleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (rh *groupRoleHandler) GetByID(c *gin.Context) {
	roleID := c.Param("id")

	groupRole, err := rh.groupRoleUsecase.GetByID(roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GroupRoleApiResponse{
		Message: "success",
		Data:    groupRole,
	})
}
