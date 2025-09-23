package models

// This file serves as the main entry point for all database models
// It provides a convenient way to import all models from other packages

// AllModels returns a slice of all model structs for database migration
// This is used by GORM AutoMigrate to create/update database tables
func AllModels() []interface{} {
	return []interface{}{
		&Category{},
		&Todo{},
	}
}