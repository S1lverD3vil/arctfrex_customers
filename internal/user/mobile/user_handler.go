package user

import (
	"log"
	"net/http"

	"arctfrex-customers/internal/middleware"

	"github.com/gin-gonic/gin"
)

// userHandler handles HTTP requests for user operations
type userHandler struct {
	jwtMiddleware *middleware.JWTMiddleware
	userUsecase   UserUsecase
}

// NewUserHandler sets up the HTTP handlers for user operations
func NewUserHandler(
	engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	uu *userUsecase,
) *userHandler {
	handler := &userHandler{
		jwtMiddleware: jmw,
		userUsecase:   uu,
	}

	unprotectedRoutesBackOffice := engine.Group("/backoffice/customers/users")
	unprotectedRoutes, protectedRoutes := engine.Group("/users"), engine.Group("/users")

	unprotectedRoutesBackOffice.GET("/profile/:userid", handler.BackOfficeCustomersGetProfile)
	unprotectedRoutesBackOffice.GET("/address/:userid", handler.BackOfficeCustomersGetAddress)
	unprotectedRoutesBackOffice.GET("/employment/:userid", handler.BackOfficeCustomersGetEmployment)
	unprotectedRoutesBackOffice.GET("/finance/:userid", handler.BackOfficeCustomersGetFinance)
	unprotectedRoutesBackOffice.GET("/emergencycontact/:userid", handler.BackOfficeCustomersGetEmergencyContact)
	unprotectedRoutesBackOffice.POST("/register", handler.Register)
	unprotectedRoutesBackOffice.GET("/check/:mobilePhone", handler.Check)
	unprotectedRoutesBackOffice.PATCH("/pin", handler.UpdatePin)
	unprotectedRoutesBackOffice.GET("/leads", handler.BackOfficeLeads)

	unprotectedRoutes.POST("/register", handler.Register)
	unprotectedRoutes.GET("/check/:mobilePhone", handler.Check)
	unprotectedRoutes.PATCH("/pin", handler.UpdatePin)
	unprotectedRoutes.POST("/login/session", handler.LoginSession)
	protectedRoutes.Use(jmw.ValidateToken())
	{
		protectedRoutes.POST("/session", handler.Session)
		protectedRoutes.POST("/logout/session", handler.LogoutSession)
		protectedRoutes.DELETE("/delete", handler.Delete)
		protectedRoutes.PATCH("/profile", handler.UpdateProfile)
		protectedRoutes.GET("/profile", handler.GetProfile)
		protectedRoutes.PATCH("/address", handler.UpdateAddress)
		protectedRoutes.GET("/address", handler.GetAddress)
		protectedRoutes.PATCH("/employment", handler.UpdateEmployment)
		protectedRoutes.GET("/employment", handler.GetEmployment)
		protectedRoutes.PATCH("/finance", handler.UpdateFinance)
		protectedRoutes.GET("/finance", handler.GetFinance)
		protectedRoutes.PATCH("/emergencycontact", handler.UpdateEmergencyContact)
		protectedRoutes.GET("/emergencycontact", handler.GetEmergencyContact)
	}

	return handler
}

// CreateUser handles the creation of a new user
func (uh *userHandler) Register(c *gin.Context) {
	var user Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uh.userUsecase.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

// 200 -> input pin
// 400 -> create pin
// 402 -> signup
func (uh *userHandler) Check(c *gin.Context) {

	user, err := uh.userUsecase.Check(c.Param("mobilePhone"))
	if err != nil {
		c.JSON(402, gin.H{"error": "not registered"})
		return
	}

	if user.Pin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need to set pin"})
		return
	}

	c.Status(http.StatusOK)
}

func (uh *userHandler) LoginSession(c *gin.Context) {
	var user *UserLoginSessionRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userLoginResposnse, err := uh.userUsecase.LoginSession(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userLoginResposnse)
}

func (uh *userHandler) Session(c *gin.Context) {
	var user *UserSessionRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := uh.userUsecase.Session(&Users{SessionId: user.SessionId, DeviceId: user.DeviceId})
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	//log.Println(userdb)
	c.Status(http.StatusOK)
}

func (uh *userHandler) LogoutSession(c *gin.Context) {
	var user *UserSessionRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uh.userUsecase.LogoutSession(&Users{SessionId: user.SessionId, DeviceId: user.DeviceId})
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	//log.Println(userdb)
	c.Status(http.StatusOK)
}

func (uh *userHandler) Delete(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
		return
	}

	err := uh.userUsecase.Delete(&Users{ID: userId})
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

func (uh *userHandler) UpdatePin(c *gin.Context) {
	var user *UpdatePinRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.userUsecase.UpdatePin(user.MobilePhone, user.Pin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (uh *userHandler) UpdateProfile(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
		return
	}

	var userProfile *UserProfile
	if err := c.ShouldBindJSON(&userProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.userUsecase.UpdateProfile(userId, userProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (uh *userHandler) GetProfile(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
		return
	}

	userProfile, err := uh.userUsecase.GetProfile(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userProfile)
}

func (uh *userHandler) BackOfficeCustomersGetProfile(c *gin.Context) {
	userProfile, err := uh.userUsecase.GetProfile(c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userProfile)
}

func (uh *userHandler) UpdateAddress(c *gin.Context) {
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

	var userAddress *UserAddress
	if err := c.ShouldBindJSON(&userAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.userUsecase.UpdateAddress(userId, userAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (uh *userHandler) GetAddress(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
		return
	}

	userAddress, err := uh.userUsecase.GetAddress(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userAddress)
}

func (uh *userHandler) BackOfficeCustomersGetAddress(c *gin.Context) {
	userAddress, err := uh.userUsecase.GetAddress(c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userAddress)
}

func (uh *userHandler) UpdateEmployment(c *gin.Context) {
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

	var userEmployment *UserEmployment
	if err := c.ShouldBindJSON(&userEmployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.userUsecase.UpdateEmployment(userId, userEmployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (uh *userHandler) GetEmployment(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
		return
	}

	userEmployment, err := uh.userUsecase.GetEmployment(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userEmployment)
}

func (uh *userHandler) BackOfficeCustomersGetEmployment(c *gin.Context) {
	userEmployment, err := uh.userUsecase.GetEmployment(c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userEmployment)
}

func (uh *userHandler) UpdateFinance(c *gin.Context) {
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

	var userFinance *UserFinance
	if err := c.ShouldBindJSON(&userFinance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.userUsecase.UpdateFinance(userId, userFinance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (uh *userHandler) GetFinance(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
		return
	}

	userFinance, err := uh.userUsecase.GetFinance(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userFinance)
}

func (uh *userHandler) BackOfficeCustomersGetFinance(c *gin.Context) {
	userFinance, err := uh.userUsecase.GetFinance(c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userFinance)
}

func (uh *userHandler) UpdateEmergencyContact(c *gin.Context) {
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

	var userEmergencyContact *UserEmergencyContact
	if err := c.ShouldBindJSON(&userEmergencyContact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.userUsecase.UpdateEmergencyContact(userId, userEmergencyContact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (uh *userHandler) GetEmergencyContact(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
		return
	}

	userEmergencyContact, err := uh.userUsecase.GetEmergencyContact(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userEmergencyContact)
}

func (uh *userHandler) BackOfficeCustomersGetEmergencyContact(c *gin.Context) {
	userEmergencyContact, err := uh.userUsecase.GetEmergencyContact(c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userEmergencyContact)
}

func (uh *userHandler) BackOfficeLeads(c *gin.Context) {
	backOfficeUserLeads, err := uh.userUsecase.BackOfficeGetLeads()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(
		http.StatusOK,
		UserApiResponse{Message: "success", Data: backOfficeUserLeads},
	)
}
