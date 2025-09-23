package repository

import (
	"todo-backend/internal/models"
)

// TodoRepository defines the interface for todo data operations
type TodoRepository interface {
	// Create creates a new todo
	Create(todo *models.Todo) error
	
	// GetByID retrieves a todo by its ID
	GetByID(id uint) (*models.Todo, error)
	
	// Update updates an existing todo
	Update(todo *models.Todo) error
	
	// Delete soft deletes a todo by ID
	Delete(id uint) error
	
	// List retrieves todos with pagination and filtering
	List(filters TodoFilters, pagination PaginationParams) ([]models.Todo, PaginationResult, error)
	
	// ToggleComplete toggles the completion status of a todo
	ToggleComplete(id uint) error
}

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	// Create creates a new category
	Create(category *models.Category) error
	
	// GetByID retrieves a category by its ID
	GetByID(id uint) (*models.Category, error)
	
	// Update updates an existing category
	Update(category *models.Category) error
	
	// Delete soft deletes a category by ID
	Delete(id uint) error
	
	// List retrieves categories with pagination and filtering
	List(filters CategoryFilters, pagination PaginationParams) ([]models.Category, PaginationResult, error)
	
	// GetAll retrieves all categories without pagination (for dropdowns)
	GetAll() ([]models.Category, error)
}