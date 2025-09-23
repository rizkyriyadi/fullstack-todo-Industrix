package models

import (
	"time"

	"gorm.io/gorm"
)

// Priority represents the priority levels for todos
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// IsValid checks if the priority value is valid
func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}

// Todo represents a todo item in the system
// Each todo belongs to a category and has various attributes for organization
type Todo struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Title       string         `json:"title" gorm:"not null;size:255" binding:"required,min=1,max=255"`
	Description string         `json:"description" gorm:"type:text"`
	Completed   bool           `json:"completed" gorm:"default:false"`
	Priority    Priority       `json:"priority" gorm:"type:varchar(10);default:'medium'" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time     `json:"due_date,omitempty" gorm:"index"`
	CategoryID  *uint          `json:"category_id,omitempty" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship: Todo belongs to a category
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID"`
}

// TableName returns the table name for Todo model
func (Todo) TableName() string {
	return "todos"
}

// BeforeCreate hook runs before creating a todo
// Sets default priority if not provided
func (t *Todo) BeforeCreate(tx *gorm.DB) error {
	if t.Priority == "" {
		t.Priority = PriorityMedium
	}
	return nil
}

// BeforeUpdate hook runs before updating a todo
// Validates priority if it's being updated
func (t *Todo) BeforeUpdate(tx *gorm.DB) error {
	if t.Priority != "" && !t.Priority.IsValid() {
		return gorm.ErrInvalidValue
	}
	return nil
}

// ToggleComplete toggles the completed status of the todo
func (t *Todo) ToggleComplete() {
	t.Completed = !t.Completed
}