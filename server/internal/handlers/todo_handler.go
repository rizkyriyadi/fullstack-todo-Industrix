package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"todo-backend/internal/models"
	"todo-backend/internal/repository"
	"todo-backend/internal/services"
	"todo-backend/pkg/utils"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	todoService services.TodoService
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(todoService services.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// CreateTodo handles POST /api/todos
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var todo models.Todo

	// Bind JSON to todo struct with validation
	if err := c.ShouldBindJSON(&todo); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Create the todo using service
	if err := h.todoService.CreateTodo(&todo); err != nil {
		if strings.Contains(err.Error(), "not found") || 
		   strings.Contains(err.Error(), "does not exist") {
			utils.ValidationErrorResponse(c, err)
			return
		}
		if strings.Contains(err.Error(), "past") ||
		   strings.Contains(err.Error(), "invalid") {
			utils.ValidationErrorResponse(c, err)
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Todo created successfully", todo)
}

// GetTodo handles GET /api/todos/:id
func (h *TodoHandler) GetTodo(c *gin.Context) {
	// Extract ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Get todo using service
	todo, err := h.todoService.GetTodoByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.NotFoundErrorResponse(c, "Todo")
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Todo retrieved successfully", todo)
}

// UpdateTodo handles PUT /api/todos/:id
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	// Extract ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	var todo models.Todo

	// Bind JSON to todo struct with validation
	if err := c.ShouldBindJSON(&todo); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Set the ID from URL parameter
	todo.ID = uint(id)

	// Update the todo using service
	if err := h.todoService.UpdateTodo(&todo); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.NotFoundErrorResponse(c, "Todo")
			return
		}
		if strings.Contains(err.Error(), "past") ||
		   strings.Contains(err.Error(), "invalid") ||
		   strings.Contains(err.Error(), "does not exist") {
			utils.ValidationErrorResponse(c, err)
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Todo updated successfully", todo)
}

// DeleteTodo handles DELETE /api/todos/:id
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	// Extract ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Delete the todo using service
	if err := h.todoService.DeleteTodo(uint(id)); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.NotFoundErrorResponse(c, "Todo")
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Todo deleted successfully", nil)
}

// ListTodos handles GET /api/todos
func (h *TodoHandler) ListTodos(c *gin.Context) {
	// Parse query parameters for filtering and pagination
	var filters repository.TodoFilters
	var pagination repository.PaginationParams

	// Bind query parameters
	if err := c.ShouldBindQuery(&filters); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := c.ShouldBindQuery(&pagination); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Get todos using service
	todos, paginationResult, err := h.todoService.ListTodos(filters, pagination)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			utils.ValidationErrorResponse(c, err)
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Todos retrieved successfully", todos, paginationResult)
}

// ToggleTodoComplete handles PATCH /api/todos/:id/complete
func (h *TodoHandler) ToggleTodoComplete(c *gin.Context) {
	// Extract ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Toggle todo completion using service
	if err := h.todoService.ToggleTodoComplete(uint(id)); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.NotFoundErrorResponse(c, "Todo")
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Todo completion status toggled successfully", nil)
}