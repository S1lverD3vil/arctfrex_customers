package role

import (
	"arctfrex-customers/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type roleHandler struct {
	jwtMiddleware *middleware.JWTMiddleware
	roleUsecase   RoleUseCase
}

func NewRoleHandler(engine *gin.Engine, jmw *middleware.JWTMiddleware, ru *roleUseCase) *roleHandler {
	handler := &roleHandler{
		jwtMiddleware: jmw,
		roleUsecase:   ru,
	}

	unprotectedRoutes := engine.Group("/backoffice/roles")

	unprotectedRoutes.GET("/all", handler.All)
	unprotectedRoutes.POST("", handler.Create)
	unprotectedRoutes.PUT("/:id", handler.Update)
	unprotectedRoutes.DELETE("/:id", handler.Delete)
	unprotectedRoutes.GET("/:id", handler.GetByID)

	return handler
}

func (rh *roleHandler) All(c *gin.Context) {
	roles, err := rh.roleUsecase.All()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, RoleApiResponse{
		Message: "succes", Data: roles,
	})
}

func (rh *roleHandler) Create(c *gin.Context) {
	var role CreateUserDTO
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newRole := Role{
		ID:             role.ID,
		Name:           role.Name,
		CommissionRate: role.CommissionRate,
		ParentRoleID:   role.ParentRoleID,
	}

	if err := rh.roleUsecase.Create(&newRole); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, RoleApiResponse{
		Message: "succes",
		Data:    role,
	})
}

func (rh *roleHandler) Update(c *gin.Context) {
	roleID := c.Param("id")
	var role Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role.ID = roleID
	if err := rh.roleUsecase.Update(&role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, RoleApiResponse{
		Message: "success",
		Data:    role,
	})
}

func (rh *roleHandler) Delete(c *gin.Context) {
	roleID := c.Param("id")
	if err := rh.roleUsecase.Delete(roleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (rh *roleHandler) GetByID(c *gin.Context) {
	roleID := c.Param("id")

	role, err := rh.roleUsecase.GetByID(roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, RoleApiResponse{
		Message: "success",
		Data:    role,
	})
}
