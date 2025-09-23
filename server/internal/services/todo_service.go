package services

import (
	"errors"
	"strings"
	"time"

	"todo-backend/internal/models"
	"todo-backend/internal/repository"
)

// todoService implements TodoService interface
type todoService struct {
	todoRepo     repository.TodoRepository
	categoryRepo repository.CategoryRepository
}

// NewTodoService creates a new todo service
func NewTodoService(todoRepo repository.TodoRepository, categoryRepo repository.CategoryRepository) TodoService {
	return &todoService{
		todoRepo:     todoRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateTodo creates a new todo with validation
func (s *todoService) CreateTodo(todo *models.Todo) error {
	// Business logic validation
	if err := s.validateTodo(todo); err != nil {
		return err
	}

	// Clean and format the data
	s.cleanTodoData(todo)

	// Additional business logic validation
	if err := s.validateTodoBusinessRules(todo); err != nil {
		return err
	}

	return s.todoRepo.Create(todo)
}

// GetTodoByID retrieves a todo by its ID
func (s *todoService) GetTodoByID(id uint) (*models.Todo, error) {
	if id == 0 {
		return nil, errors.New("invalid todo ID")
	}
	return s.todoRepo.GetByID(id)
}

// UpdateTodo updates an existing todo with validation
func (s *todoService) UpdateTodo(todo *models.Todo) error {
	if todo.ID == 0 {
		return errors.New("invalid todo ID")
	}

	// Business logic validation
	if err := s.validateTodo(todo); err != nil {
		return err
	}

	// Clean and format the data
	s.cleanTodoData(todo)

	// Additional business logic validation
	if err := s.validateTodoBusinessRules(todo); err != nil {
		return err
	}

	return s.todoRepo.Update(todo)
}

// DeleteTodo soft deletes a todo by ID
func (s *todoService) DeleteTodo(id uint) error {
	if id == 0 {
		return errors.New("invalid todo ID")
	}
	return s.todoRepo.Delete(id)
}

// ListTodos retrieves todos with pagination and filtering
func (s *todoService) ListTodos(filters repository.TodoFilters, pagination repository.PaginationParams) ([]models.Todo, repository.PaginationResult, error) {
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

	// Validate priority filter
	if filters.Priority != "" {
		priority := models.Priority(filters.Priority)
		if !priority.IsValid() {
			return nil, repository.PaginationResult{}, errors.New("invalid priority filter")
		}
	}

	return s.todoRepo.List(filters, pagination)
}

// ToggleTodoComplete toggles the completion status of a todo
func (s *todoService) ToggleTodoComplete(id uint) error {
	if id == 0 {
		return errors.New("invalid todo ID")
	}
	return s.todoRepo.ToggleComplete(id)
}

// validateTodo validates basic todo data
func (s *todoService) validateTodo(todo *models.Todo) error {
	if todo == nil {
		return errors.New("todo cannot be nil")
	}

	// Validate title
	if strings.TrimSpace(todo.Title) == "" {
		return errors.New("todo title is required")
	}

	if len(todo.Title) > 255 {
		return errors.New("todo title cannot exceed 255 characters")
	}

	// Validate description length
	if len(todo.Description) > 5000 {
		return errors.New("todo description cannot exceed 5000 characters")
	}

	// Validate priority
	if todo.Priority != "" && !todo.Priority.IsValid() {
		return errors.New("invalid priority value")
	}

	return nil
}

// validateTodoBusinessRules validates business rules for todos
func (s *todoService) validateTodoBusinessRules(todo *models.Todo) error {
	// Validate due date is not in the past (only for new todos or when due date is being changed)
	if todo.DueDate != nil {
		now := time.Now().UTC()
		// Allow some tolerance for timezone differences (1 day)
		tolerance := 24 * time.Hour
		if todo.DueDate.Before(now.Add(-tolerance)) {
			return errors.New("due date cannot be in the past")
		}
	}

	// Validate category exists if provided
	if todo.CategoryID != nil {
		if _, err := s.categoryRepo.GetByID(*todo.CategoryID); err != nil {
			return errors.New("specified category does not exist")
		}
	}

	return nil
}

// cleanTodoData cleans and formats todo data
func (s *todoService) cleanTodoData(todo *models.Todo) {
	// Trim whitespace from title and description
	todo.Title = strings.TrimSpace(todo.Title)
	todo.Description = strings.TrimSpace(todo.Description)

	// Set default priority if empty
	if todo.Priority == "" {
		todo.Priority = models.PriorityMedium
	}

	// Normalize due date to UTC if provided
	if todo.DueDate != nil {
		utcTime := todo.DueDate.UTC()
		todo.DueDate = &utcTime
	}
}