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

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	categoryService services.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory handles POST /api/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category

	// Bind JSON to category struct with validation
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Create the category using service
	if err := h.categoryService.CreateCategory(&category); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			utils.ConflictErrorResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Category created successfully", category)
}

// GetCategory handles GET /api/categories/:id
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	// Extract ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Get category using service
	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.NotFoundErrorResponse(c, "Category")
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

// UpdateCategory handles PUT /api/categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	// Extract ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	var category models.Category

	// Bind JSON to category struct with validation
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Set the ID from URL parameter
	category.ID = uint(id)

	// Update the category using service
	if err := h.categoryService.UpdateCategory(&category); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.NotFoundErrorResponse(c, "Category")
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			utils.ConflictErrorResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category updated successfully", category)
}

// DeleteCategory handles DELETE /api/categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	// Extract ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Delete the category using service
	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.NotFoundErrorResponse(c, "Category")
			return
		}
		if strings.Contains(err.Error(), "associated todos") {
			utils.ConflictErrorResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}

// ListCategories handles GET /api/categories
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	// Parse query parameters for filtering and pagination
	var filters repository.CategoryFilters
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

	// Get categories using service
	categories, paginationResult, err := h.categoryService.ListCategories(filters, pagination)
	if err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Categories retrieved successfully", categories, paginationResult)
}

// GetAllCategories handles GET /api/categories/all
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	// Get all categories using service (for dropdowns)
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "All categories retrieved successfully", categories)
}