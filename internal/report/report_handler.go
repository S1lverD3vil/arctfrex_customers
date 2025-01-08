package report

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type reportHandler struct {
	jwtMiddleware *middleware.JWTMiddleware
	reportUsecase ReportUsecase
}

func NewReportHandler(
	engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	ru ReportUsecase,
) *reportHandler {
	handler := &reportHandler{
		jwtMiddleware: jmw,
		reportUsecase: ru,
	}
	unprotectedRoutesBackOffice := engine.Group("backoffice/reports")
	protectedRoutes := engine.Group("/reports")

	unprotectedRoutesBackOffice.GET("/all", handler.BackOfficeGetActiveReports)
	unprotectedRoutesBackOffice.GET("/:reportCode", handler.BackOfficeGetActiveReportsByCode)

	protectedRoutes.Use(jmw.ValidateToken())
	{
	}

	return handler
}

func (rh *reportHandler) BackOfficeGetActiveReports(c *gin.Context) {
	reports, err := rh.reportUsecase.GetActiveReports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		ReportApiResponse{base.ApiResponse{Message: "success", Data: reports}},
	)
}

func (rh *reportHandler) BackOfficeGetActiveReportsByCode(c *gin.Context) {
	reports, err := rh.reportUsecase.GetActiveReportsByCode(c.Param("reportCode"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		ReportApiResponse{base.ApiResponse{Message: "success", Data: reports}},
	)
}
