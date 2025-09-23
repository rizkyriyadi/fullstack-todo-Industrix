package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"todo-backend/internal/repository"
)

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool                         `json:"success"`
	Message    string                       `json:"message"`
	Data       interface{}                  `json:"data"`
	Pagination repository.PaginationResult `json:"pagination"`
	Error      string                       `json:"error,omitempty"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   message,
	})
}

// PaginatedSuccessResponse sends a successful paginated response
func PaginatedSuccessResponse(c *gin.Context, message string, data interface{}, pagination repository.PaginationResult) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusBadRequest, "Validation error: "+err.Error())
}

// NotFoundErrorResponse sends a not found error response
func NotFoundErrorResponse(c *gin.Context, resource string) {
	ErrorResponse(c, http.StatusNotFound, resource+" not found")
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusInternalServerError, "Internal server error: "+err.Error())
}

// ConflictErrorResponse sends a conflict error response
func ConflictErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, message)
}