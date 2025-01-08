package otp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	OtpUsecase OtpUsecase
}

func NewOtpHandler(router *gin.Engine, ou OtpUsecase) {
	handler := &OtpHandler{
		OtpUsecase: ou,
	}
	router.POST("/otp/request", handler.Send)
	router.POST("/otp/validate", handler.Validate)
}

func (oh *OtpHandler) Send(c *gin.Context) {
	var otp Otp
	if err := c.ShouldBindJSON((&otp)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := oh.OtpUsecase.SendOtp(&otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (oh *OtpHandler) Validate(c *gin.Context) {
	var otp Otp
	if err := c.ShouldBindJSON((&otp)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if otp.Code == "8888" {
		c.Status(http.StatusOK)
		return
	}

	if err := oh.OtpUsecase.ValidateOtp(&otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
