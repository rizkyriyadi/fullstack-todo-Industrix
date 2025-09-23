package services

import (
	"todo-backend/internal/models"
	"todo-backend/internal/repository"
)

// TodoService defines the interface for todo business logic
type TodoService interface {
	// CreateTodo creates a new todo with validation
	CreateTodo(todo *models.Todo) error
	
	// GetTodoByID retrieves a todo by its ID
	GetTodoByID(id uint) (*models.Todo, error)
	
	// UpdateTodo updates an existing todo with validation
	UpdateTodo(todo *models.Todo) error
	
	// DeleteTodo soft deletes a todo by ID
	DeleteTodo(id uint) error
	
	// ListTodos retrieves todos with pagination and filtering
	ListTodos(filters repository.TodoFilters, pagination repository.PaginationParams) ([]models.Todo, repository.PaginationResult, error)
	
	// ToggleTodoComplete toggles the completion status of a todo
	ToggleTodoComplete(id uint) error
}

// CategoryService defines the interface for category business logic
type CategoryService interface {
	// CreateCategory creates a new category with validation
	CreateCategory(category *models.Category) error
	
	// GetCategoryByID retrieves a category by its ID
	GetCategoryByID(id uint) (*models.Category, error)
	
	// UpdateCategory updates an existing category with validation
	UpdateCategory(category *models.Category) error
	
	// DeleteCategory soft deletes a category by ID
	DeleteCategory(id uint) error
	
	// ListCategories retrieves categories with pagination and filtering
	ListCategories(filters repository.CategoryFilters, pagination repository.PaginationParams) ([]models.Category, repository.PaginationResult, error)
	
	// GetAllCategories retrieves all categories without pagination (for dropdowns)
	GetAllCategories() ([]models.Category, error)
}