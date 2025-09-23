package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"todo-backend/internal/models"
)

// todoRepository implements TodoRepository interface
type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a new todo repository
func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{
		db: db,
	}
}

// Create creates a new todo
func (r *todoRepository) Create(todo *models.Todo) error {
	// Validate category exists if provided
	if todo.CategoryID != nil {
		var category models.Category
		if err := r.db.First(&category, *todo.CategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("category not found")
			}
			return err
		}
	}

	return r.db.Create(todo).Error
}

// GetByID retrieves a todo by its ID
func (r *todoRepository) GetByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.Preload("Category").First(&todo, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}
	return &todo, nil
}

// Update updates an existing todo
func (r *todoRepository) Update(todo *models.Todo) error {
	// Check if todo exists
	var existingTodo models.Todo
	if err := r.db.First(&existingTodo, todo.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("todo not found")
		}
		return err
	}

	// Validate category exists if provided
	if todo.CategoryID != nil {
		var category models.Category
		if err := r.db.First(&category, *todo.CategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("category not found")
			}
			return err
		}
	}

	// Update the todo
	return r.db.Save(todo).Error
}

// Delete soft deletes a todo by ID
func (r *todoRepository) Delete(id uint) error {
	// Check if todo exists
	var todo models.Todo
	if err := r.db.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("todo not found")
		}
		return err
	}

	// Soft delete the todo
	return r.db.Delete(&todo).Error
}

// List retrieves todos with pagination and filtering
func (r *todoRepository) List(filters TodoFilters, pagination PaginationParams) ([]models.Todo, PaginationResult, error) {
	var todos []models.Todo
	var total int64

	// Build the base query with category preload
	query := r.db.Model(&models.Todo{}).Preload("Category")

	// Apply search filter (search in title using full-text search)
	if filters.Search != "" {
		searchTerm := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where("LOWER(title) LIKE ?", searchTerm)
	}

	// Apply completion filter
	if filters.Completed != nil {
		query = query.Where("completed = ?", *filters.Completed)
	}

	// Apply category filter
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}

	// Apply priority filter
	if filters.Priority != "" {
		query = query.Where("priority = ?", filters.Priority)
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
			"title":      true,
			"completed":  true,
			"priority":   true,
			"due_date":   true,
			"created_at": true,
			"updated_at": true,
		}
		if validSortFields[pagination.SortBy] {
			sortBy = pagination.SortBy
		}
	}

	// Handle sorting for priority (custom order: high, medium, low)
	if sortBy == "priority" {
		orderClause := "CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 END"
		if pagination.GetSortOrder() == "desc" {
			orderClause += " DESC"
		}
		query = query.Order(orderClause)
	} else {
		query = query.Order(sortBy + " " + pagination.GetSortOrder())
	}

	// Apply pagination
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	if err := query.Offset(offset).Limit(limit).Find(&todos).Error; err != nil {
		return nil, PaginationResult{}, err
	}

	// Calculate pagination result
	paginationResult := PaginationResult{
		CurrentPage: pagination.Page,
		PerPage:     limit,
		Total:       total,
	}
	paginationResult.CalculateTotalPages()

	return todos, paginationResult, nil
}

// ToggleComplete toggles the completion status of a todo
func (r *todoRepository) ToggleComplete(id uint) error {
	// Get the current todo
	var todo models.Todo
	if err := r.db.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("todo not found")
		}
		return err
	}

	// Toggle the completion status
	todo.ToggleComplete()

	// Save the updated todo
	return r.db.Save(&todo).Error
}