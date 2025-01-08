package infrastructure

import (
	"arctfrex-customers/internal/common"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestResponseLogger logs request and response details
func RequestResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		separator := strings.Repeat("-", 100)
		// Log the separator
		fmt.Printf("%s\n", separator)

		// Log request headers
		log.Println("Request Headers:")
		for key, values := range c.Request.Header {
			log.Printf("%s: %s\n", key, strings.Join(values, ", "))
		}

		var requestBody string
		if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
			// Handle multipart file uploads
			err := c.Request.ParseMultipartForm(32 << 20) // 32 MB is the default max memory
			if err != nil {
				log.Printf("Error parsing multipart form: %v", err)
			} else {
				// Log details for each uploaded file
				for _, fileHeaders := range c.Request.MultipartForm.File {
					for _, fileHeader := range fileHeaders {
						fileType := fileHeader.Header.Get("Content-Type") // Get file type
						requestBody += fmt.Sprintf("Uploaded File: %s, Type: %s\n", fileHeader.Filename, fileType)
					}
				}

				// Log other form fields
				for key, values := range c.Request.MultipartForm.Value {
					for _, value := range values {
						requestBody += fmt.Sprintf("Form Field: %s, Value: %s\n", key, value)
					}
				}
			}
		} else {
			// Log the request body for non-multipart forms
			rawBody, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody)) // Restore the request body
			requestBody = string(rawBody)
		}

		log.Printf("Request: %s %s\nBody: %s\n", c.Request.Method, c.Request.URL.Path, requestBody)

		// Create a buffer to capture the response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		// Attempt to format the response body if it's JSON
		var formattedBody string
		var responseBody map[string]interface{}

		if jsonErr := json.Unmarshal(blw.body.Bytes(), &responseBody); jsonErr == nil {
			// If the response body is JSON, pretty-print it
			formattedBodyBytes, _ := json.MarshalIndent(responseBody, common.STRING_EMPTY, common.STRING_DOUBLE_SPACE)
			formattedBody = string(formattedBodyBytes)
		} else {
			// If it's not JSON, use the raw body
			formattedBody = blw.body.String()
		}

		// Log response headers
		log.Println("Response Headers:")
		for key, values := range c.Writer.Header() {
			log.Printf("%s: %s\n", key, strings.Join(values, ", "))
		}

		// Log response details
		log.Printf("Response: Status=%d, Duration=%s\nBody: %s\n", c.Writer.Status(), duration, formattedBody)
		// Log the separator
		fmt.Printf("%s\n", separator)
	}
}
