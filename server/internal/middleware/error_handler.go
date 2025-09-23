package middleware

import (
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	"todo-backend/pkg/utils"
)

// ErrorHandler handles panics and errors globally
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log the panic
		log.Printf("Panic recovered: %v", recovered)

		// Send error response
		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
	})
}

// NotFoundHandler handles 404 errors
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.ErrorResponse(c, http.StatusNotFound, "Endpoint not found")
	}
}

// MethodNotAllowedHandler handles 405 errors
func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.ErrorResponse(c, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// ErrorLogger logs errors that occur during request processing
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors
		if len(c.Errors) > 0 {
			// Log all errors
			for _, err := range c.Errors {
				log.Printf("Request error: %v", err.Error())
			}
		}
	}
}