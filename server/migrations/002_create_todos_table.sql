-- Migration: Create todos table
-- This migration creates the todos table with relationships to categories
-- Includes proper indexes for filtering, searching, and sorting

-- +migrate Up
CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    priority VARCHAR(10) NOT NULL DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
    due_date TIMESTAMP WITH TIME ZONE,
    category_id INTEGER REFERENCES categories(id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better query performance
-- Index for soft deletes (GORM requirement)
CREATE INDEX IF NOT EXISTS idx_todos_deleted_at ON todos(deleted_at);

-- Index for category foreign key
CREATE INDEX IF NOT EXISTS idx_todos_category_id ON todos(category_id) WHERE deleted_at IS NULL;

-- Index for filtering by completion status
CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos(completed) WHERE deleted_at IS NULL;

-- Index for filtering by priority
CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos(priority) WHERE deleted_at IS NULL;

-- Index for sorting by due date
CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date) WHERE deleted_at IS NULL;

-- Index for sorting by creation date (most common sort)
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at) WHERE deleted_at IS NULL;

-- Composite index for filtering and sorting (common query pattern)
CREATE INDEX IF NOT EXISTS idx_todos_completed_created_at ON todos(completed, created_at DESC) WHERE deleted_at IS NULL;

-- Full-text search index for title (used for search functionality)
CREATE INDEX IF NOT EXISTS idx_todos_title_search ON todos USING gin(to_tsvector('english', title)) WHERE deleted_at IS NULL;

-- +migrate Down
DROP TABLE IF EXISTS todos;