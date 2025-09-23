package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"todo-backend/internal/models"
)

// categoryRepository implements CategoryRepository interface
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

// Create creates a new category
func (r *categoryRepository) Create(category *models.Category) error {
	if err := r.db.Create(category).Error; err != nil {
		// Handle unique constraint violation
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return errors.New("category name already exists")
		}
		return err
	}
	return nil
}

// GetByID retrieves a category by its ID
func (r *categoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// Update updates an existing category
func (r *categoryRepository) Update(category *models.Category) error {
	// Check if category exists
	var existingCategory models.Category
	if err := r.db.First(&existingCategory, category.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	// Update the category
	if err := r.db.Save(category).Error; err != nil {
		// Handle unique constraint violation
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return errors.New("category name already exists")
		}
		return err
	}
	return nil
}

// Delete soft deletes a category by ID
func (r *categoryRepository) Delete(id uint) error {
	// Check if category exists
	var category models.Category
	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	// Check if category has associated todos
	var todoCount int64
	r.db.Model(&models.Todo{}).Where("category_id = ?", id).Count(&todoCount)
	if todoCount > 0 {
		return errors.New("cannot delete category with associated todos")
	}

	// Soft delete the category
	return r.db.Delete(&category).Error
}

// List retrieves categories with pagination and filtering
func (r *categoryRepository) List(filters CategoryFilters, pagination PaginationParams) ([]models.Category, PaginationResult, error) {
	var categories []models.Category
	var total int64

	// Build the base query
	query := r.db.Model(&models.Category{})

	// Apply search filter
	if filters.Search != "" {
		searchTerm := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where("LOWER(name) LIKE ?", searchTerm)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, PaginationResult{}, err
	}

	// Apply sorting
	sortBy := "created_at"
	if pagination.SortBy != "" {
		// Validate sort field
		validSortFields := map[string]bool{
			"name":       true,
			"created_at": true,
			"updated_at": true,
		}
		if validSortFields[pagination.SortBy] {
			sortBy = pagination.SortBy
		}
	}
	query = query.Order(sortBy + " " + pagination.GetSortOrder())

	// Apply pagination
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	if err := query.Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		return nil, PaginationResult{}, err
	}

	// Calculate pagination result
	paginationResult := PaginationResult{
		CurrentPage: pagination.Page,
		PerPage:     limit,
		Total:       total,
	}
	paginationResult.CalculateTotalPages()

	return categories, paginationResult, nil
}

// GetAll retrieves all categories without pagination (for dropdowns)
func (r *categoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Order("name ASC").Find(&categories).Error
	return categories, err
}