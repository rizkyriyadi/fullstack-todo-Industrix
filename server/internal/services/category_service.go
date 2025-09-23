package services

import (
	"errors"
	"strings"

	"todo-backend/internal/models"
	"todo-backend/internal/repository"
)

// categoryService implements CategoryService interface
type categoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory creates a new category with validation
func (s *categoryService) CreateCategory(category *models.Category) error {
	// Business logic validation
	if err := s.validateCategory(category); err != nil {
		return err
	}

	// Clean and format the data
	s.cleanCategoryData(category)

	return s.categoryRepo.Create(category)
}

// GetCategoryByID retrieves a category by its ID
func (s *categoryService) GetCategoryByID(id uint) (*models.Category, error) {
	if id == 0 {
		return nil, errors.New("invalid category ID")
	}
	return s.categoryRepo.GetByID(id)
}

// UpdateCategory updates an existing category with validation
func (s *categoryService) UpdateCategory(category *models.Category) error {
	if category.ID == 0 {
		return errors.New("invalid category ID")
	}

	// Business logic validation
	if err := s.validateCategory(category); err != nil {
		return err
	}

	// Clean and format the data
	s.cleanCategoryData(category)

	return s.categoryRepo.Update(category)
}

// DeleteCategory soft deletes a category by ID
func (s *categoryService) DeleteCategory(id uint) error {
	if id == 0 {
		return errors.New("invalid category ID")
	}
	return s.categoryRepo.Delete(id)
}

// ListCategories retrieves categories with pagination and filtering
func (s *categoryService) ListCategories(filters repository.CategoryFilters, pagination repository.PaginationParams) ([]models.Category, repository.PaginationResult, error) {
	// Set default pagination values
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}

	// Clean search filter
	if filters.Search != "" {
		filters.Search = strings.TrimSpace(filters.Search)
	}

	return s.categoryRepo.List(filters, pagination)
}

// GetAllCategories retrieves all categories without pagination (for dropdowns)
func (s *categoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

// validateCategory validates category data
func (s *categoryService) validateCategory(category *models.Category) error {
	if category == nil {
		return errors.New("category cannot be nil")
	}

	// Validate name
	if strings.TrimSpace(category.Name) == "" {
		return errors.New("category name is required")
	}

	if len(category.Name) > 100 {
		return errors.New("category name cannot exceed 100 characters")
	}

	// Validate color format (hex color)
	if category.Color != "" {
		if !isValidHexColor(category.Color) {
			return errors.New("invalid color format, must be a valid hex color (e.g., #FF0000)")
		}
	}

	return nil
}

// cleanCategoryData cleans and formats category data
func (s *categoryService) cleanCategoryData(category *models.Category) {
	// Trim whitespace from name
	category.Name = strings.TrimSpace(category.Name)

	// Ensure color has # prefix if not empty
	if category.Color != "" && !strings.HasPrefix(category.Color, "#") {
		category.Color = "#" + category.Color
	}
}

// isValidHexColor validates if the string is a valid hex color
func isValidHexColor(color string) bool {
	if len(color) != 7 {
		return false
	}
	if color[0] != '#' {
		return false
	}
	for i := 1; i < 7; i++ {
		c := color[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}