package database

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	UpSQL   string
	DownSQL string
}

// MigrationRunner handles database migrations
type MigrationRunner struct {
	db *gorm.DB
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *gorm.DB) *MigrationRunner {
	return &MigrationRunner{db: db}
}

// CreateMigrationsTable creates the migrations tracking table
func (mr *MigrationRunner) CreateMigrationsTable() error {
	return mr.db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}

// GetAppliedMigrations returns list of applied migration versions
func (mr *MigrationRunner) GetAppliedMigrations() ([]int, error) {
	var versions []int
	err := mr.db.Raw("SELECT version FROM migrations ORDER BY version").Scan(&versions).Error
	return versions, err
}

// ApplyMigration applies a single migration
func (mr *MigrationRunner) ApplyMigration(migration Migration) error {
	// Start transaction
	tx := mr.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Execute the up migration
	if err := tx.Exec(migration.UpSQL).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
	}

	// Record the migration as applied
	if err := tx.Exec("INSERT INTO migrations (version, name) VALUES (?, ?)", 
		migration.Version, migration.Name).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
	}

	return tx.Commit().Error
}

// RollbackMigration rolls back a single migration
func (mr *MigrationRunner) RollbackMigration(migration Migration) error {
	// Start transaction
	tx := mr.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Execute the down migration
	if err := tx.Exec(migration.DownSQL).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to rollback migration %d: %w", migration.Version, err)
	}

	// Remove the migration record
	if err := tx.Exec("DELETE FROM migrations WHERE version = ?", migration.Version).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %d: %w", migration.Version, err)
	}

	return tx.Commit().Error
}

// LoadMigrationsFromFS loads migrations from an embedded filesystem
func LoadMigrationsFromFS(fsys fs.FS) ([]Migration, error) {
	var migrations []Migration

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}

		migration, err := parseMigration(path, string(content))
		if err != nil {
			return err
		}

		migrations = append(migrations, migration)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// parseMigration parses a migration file content
func parseMigration(filename, content string) (Migration, error) {
	// Extract version from filename (e.g., "001_create_table.sql" -> 1)
	basename := filepath.Base(filename)
	parts := strings.SplitN(basename, "_", 2)
	if len(parts) < 2 {
		return Migration{}, fmt.Errorf("invalid migration filename format: %s", filename)
	}

	var version int
	if _, err := fmt.Sscanf(parts[0], "%d", &version); err != nil {
		return Migration{}, fmt.Errorf("invalid version in filename %s: %w", filename, err)
	}

	name := strings.TrimSuffix(parts[1], ".sql")

	// Split content into up and down parts
	upMarker := "-- +migrate Up"
	downMarker := "-- +migrate Down"

	upIndex := strings.Index(content, upMarker)
	downIndex := strings.Index(content, downMarker)

	if upIndex == -1 {
		return Migration{}, fmt.Errorf("migration %s missing '-- +migrate Up' marker", filename)
	}

	var upSQL, downSQL string

	if downIndex == -1 {
		// No down migration
		upSQL = strings.TrimSpace(content[upIndex+len(upMarker):])
	} else {
		upSQL = strings.TrimSpace(content[upIndex+len(upMarker):downIndex])
		downSQL = strings.TrimSpace(content[downIndex+len(downMarker):])
	}

	return Migration{
		Version: version,
		Name:    name,
		UpSQL:   upSQL,
		DownSQL: downSQL,
	}, nil
}