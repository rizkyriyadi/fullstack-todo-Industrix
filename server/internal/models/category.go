package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a todo category in the system
// Categories help organize todos into different groups like Work, Personal, etc.
type Category struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null;size:100" binding:"required,min=1,max=100"`
	Color     string         `json:"color" gorm:"not null;size:7" binding:"required,hexcolor"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship: One category can have many todos
	Todos []Todo `json:"todos,omitempty" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName returns the table name for Category model
func (Category) TableName() string {
	return "categories"
}

// BeforeCreate hook runs before creating a category
// Sets default color if not provided
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.Color == "" {
		c.Color = "#3B82F6" // Default blue color
	}
	return nil
}