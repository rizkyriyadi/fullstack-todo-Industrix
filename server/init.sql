-- Initialization script for Docker PostgreSQL
-- This script runs when the PostgreSQL container starts for the first time

-- Ensure the database exists
SELECT 'CREATE DATABASE todo_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'todo_db')\gexec

-- Connect to the todo_db database
\c todo_db;

-- Create any additional users or permissions if needed
-- GRANT ALL PRIVILEGES ON DATABASE todo_db TO postgres;