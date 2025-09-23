-- Migration: Create categories table
-- This migration creates the categories table with proper constraints and indexes
-- Categories are used to organize todos into different groups

-- +migrate Up
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(7) NOT NULL DEFAULT '#3B82F6',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at);
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name) WHERE deleted_at IS NULL;

-- Insert default categories
INSERT INTO categories (name, color) VALUES 
    ('Work', '#3B82F6'),
    ('Personal', '#10B981'),
    ('Shopping', '#F59E0B'),
    ('Health', '#EF4444')
ON CONFLICT (name) DO NOTHING;

-- +migrate Down
DROP TABLE IF EXISTS categories;