package storage

import (
	"arctfrex-customers/internal/middleware"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type storageHandler struct {
	jwtMiddleware  *middleware.JWTMiddleware
	storageUsecase StorageUsecase
}

func NewStorageHandler(
	engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	su StorageUsecase,
) *storageHandler {
	handler := &storageHandler{
		jwtMiddleware:  jmw,
		storageUsecase: su,
	}
	unprotectedRoutesBackOffice := engine.Group("/backoffice/storages")
	unprotectedRoutes, protectedRoutes := engine.Group("/storages"), engine.Group("/storages")

	unprotectedRoutesBackOffice.POST("/upload", handler.BackOfficeUploadFile)

	unprotectedRoutes.GET("/download/:fileName", handler.DownloadFile)
	unprotectedRoutes.GET("/preview/:fileName", handler.PreviewFile)
	unprotectedRoutes.POST("/preview/:fileName", handler.PreviewFile)
	protectedRoutes.Use(jmw.ValidateToken())
	{
		protectedRoutes.POST("/upload", handler.UploadFile)
	}

	return handler
}

func (sh *storageHandler) UploadFile(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized"})
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
	}
	documentType := c.PostForm("documentType")
	accountId := c.PostForm("accountId")
	file, _ := c.FormFile("file")
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	defer f.Close()

	contentType := file.Header.Get("Content-Type")
	err = sh.storageUsecase.UploadFile(userId, accountId, documentType, contentType, f, file.Filename, file.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func (sh *storageHandler) BackOfficeUploadFile(c *gin.Context) {
	userId := c.PostForm("userId")
	documentType := c.PostForm("documentType")
	accountId := c.PostForm("accountId")
	file, _ := c.FormFile("file")
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	defer f.Close()

	contentType := file.Header.Get("Content-Type")
	err = sh.storageUsecase.UploadFile(userId, accountId, documentType, contentType, f, file.Filename, file.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func (sh *storageHandler) DownloadFile(c *gin.Context) {
	fileName := c.Param("fileName")

	data, err := sh.storageUsecase.DownloadFile(fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", data)
}

func (sh *storageHandler) PreviewFile(c *gin.Context) {
	fileName := c.Param("fileName")

	data, err := sh.storageUsecase.PreviewFile(fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to preview file"})
		return
	}

	// Set the appropriate content type based on the file extension
	contentType := "application/octet-stream" // Default content type
	if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(fileName, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(fileName, ".gif") {
		contentType = "image/gif"
	}

	// Serve the image data
	c.Data(http.StatusOK, contentType, data)
}
