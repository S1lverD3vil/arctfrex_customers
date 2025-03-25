package report

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/middleware"
	"arctfrex-customers/internal/utils"
	"encoding/csv"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
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
	unprotectedRoutesBackOffice.GET("/:report_code", handler.BackOfficeGetActiveReportsByCode)
	unprotectedRoutesBackOffice.GET("/:report_code/xlsx", handler.BackOfficeDownloadXlsxActiveReportsByCode)
	unprotectedRoutesBackOffice.GET("/:report_code/csv", handler.BackOfficeDownloadCsvActiveReportsByCode)
	unprotectedRoutesBackOffice.GET("/:report_code/pdf", handler.BackOfficeDownloadPdfActiveReportsByCode)

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
	reportCode := c.Param("report_code")
	startDate := c.Query("start_date") // Example format: "2024-03-01"
	endDate := c.Query("end_date")

	reports, err := rh.reportUsecase.GetActiveReportsByCode(reportCode, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		ReportApiResponse{base.ApiResponse{Message: "success", Data: reports}},
	)
}

func (rh *reportHandler) BackOfficeDownloadXlsxActiveReportsByCode(c *gin.Context) {
	reportCode := c.Param("report_code")
	startDate := c.Query("start_date") // Example format: "2024-03-01"
	endDate := c.Query("end_date")

	filename := utils.GenerateReportFilename(reportCode, startDate, endDate)

	// Set Response as xlsx
	c.Header("Content-type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))

	reports, err := rh.reportUsecase.GetActiveReportsByCode(reportCode, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dataSlice, ok := reports.Data.([]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid data format"})
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Add Column Headers
	for colIdx, colName := range reports.Column {
		cell := fmt.Sprintf("%s1", string(rune(65+colIdx))) // Convert 0 -> A, 1 -> B, etc.
		f.SetCellValue(sheetName, cell, colName)
	}

	// Add Data Rows
	for rowIdx, rowData := range dataSlice {
		v := reflect.ValueOf(rowData)

		// If rowData is a pointer, get the actual struct
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		// Ensure v is a struct before accessing fields
		if v.Kind() != reflect.Struct {
			continue
		}

		for colIdx, colName := range reports.Column {
			cell := fmt.Sprintf("%s%d", string(rune(65+colIdx)), rowIdx+2)

			fieldValue := v.FieldByName(colName)
			if fieldValue.IsValid() && fieldValue.CanInterface() {
				f.SetCellValue(sheetName, cell, fieldValue.Interface())
			} else {
				f.SetCellValue(sheetName, cell, "")
			}
		}
	}

	f.SetActiveSheet(index)

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusOK, err)
	}
}

func (rh *reportHandler) BackOfficeDownloadCsvActiveReportsByCode(c *gin.Context) {
	reportCode := c.Param("report_code")
	startDate := c.Query("start_date") // Example format: "2024-03-01"
	endDate := c.Query("end_date")

	filename := utils.GenerateReportFilename(reportCode, startDate, endDate)

	// Set Response as CSV
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))

	reports, err := rh.reportUsecase.GetActiveReportsByCode(reportCode, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dataSlice, ok := reports.Data.([]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid data format"})
		return
	}

	// Create CSV writer
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write Column Headers
	if err := writer.Write(reports.Column); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write CSV headers"})
		return
	}

	// Write Data Rows
	for _, rowData := range dataSlice {
		v := reflect.ValueOf(rowData)

		// If rowData is a pointer, dereference it
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		// Ensure it's a struct before accessing fields
		if v.Kind() != reflect.Struct {
			continue
		}

		// Extract values by column name
		var row []string
		for _, colName := range reports.Column {
			fieldValue := v.FieldByName(colName)
			if fieldValue.IsValid() && fieldValue.CanInterface() {
				row = append(row, fmt.Sprintf("%v", fieldValue.Interface()))
			} else {
				row = append(row, "")
			}
		}

		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write CSV data"})
			return
		}
	}
}

func (rh *reportHandler) BackOfficeDownloadPdfActiveReportsByCode(c *gin.Context) {
	reportCode := c.Param("report_code")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	filename := utils.GenerateReportFilename(reportCode, startDate, endDate)

	reports, err := rh.reportUsecase.GetActiveReportsByCode(reportCode, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dataSlice, ok := reports.Data.([]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid data format"})
		return
	}

	// Calculate paper size based on column width
	columnWidth := 40.0
	totalWidth := float64(len(reports.Column)) * columnWidth
	if totalWidth < 210.0 { // Minimum A4 size
		totalWidth = 210.0
	}

	// Initialize PDF with custom size using NewCustom
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "L",
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: totalWidth, Ht: 297.0},
	})
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)

	// Add Title
	pdf.CellFormat(totalWidth, 10, fmt.Sprintf("Report: %s", reportCode), "0", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Set Column Headers
	pdf.SetFont("Arial", "B", 10)
	for _, col := range reports.Column {
		pdf.CellFormat(columnWidth, 7, col, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Set Data Rows
	pdf.SetFont("Arial", "", 10)
	for _, rowData := range dataSlice {
		v := reflect.ValueOf(rowData)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			continue
		}

		for _, colName := range reports.Column {
			fieldValue := v.FieldByName(colName)
			var cellValue string
			if fieldValue.IsValid() && fieldValue.CanInterface() {
				cellValue = fmt.Sprintf("%v", fieldValue.Interface())
			}
			pdf.CellFormat(columnWidth, 7, cellValue, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Write PDF to Response
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", filename))
	if err := pdf.Output(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate PDF"})
	}
}
