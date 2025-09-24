# Todo Backend API

A RESTful API server built with Go, Gin framework, and PostgreSQL for the todo application. This backend provides comprehensive CRUD operations, advanced filtering, pagination, and clean architecture.

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 12 or higher  
- Make (optional, for using Makefile commands)

### Setup Instructions

1. **Clone and navigate to server directory:**
   ```bash
   cd server
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Setup environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Create PostgreSQL database:**
   ```bash
   createdb todoapp
   ```

5. **Run the server:**
   ```bash
   make run
   # OR
   go run cmd/api/main.go
   ```

The server will start on `http://localhost:8080` and automatically run database migrations.

## 🏗️ Architecture

### Project Structure
```
server/
├── cmd/api/              # Application entry point
├── internal/             # Private application code
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers (controllers)
│   ├── middleware/      # HTTP middleware functions
│   ├── models/          # Data models and structs
│   ├── repository/      # Data access layer
│   └── services/        # Business logic layer
├── migrations/          # Database migration files
├── pkg/                # Public packages
│   ├── database/       # Database connection and utilities
│   └── utils/          # Shared utility functions
├── Makefile            # Build and development commands
├── docker-compose.yml  # Docker configuration
└── .env.example        # Environment variables template
```

### Architecture Pattern
I implemented a **layered architecture** with clear separation of concerns:

```
HTTP Handlers ← Services ← Repository ← Database
      ↑             ↑           ↑
  (Controllers) (Business)  (Data Access)
```

**Benefits:**
- Easy to test each layer independently
- Clear separation of responsibilities
- Flexible and maintainable codebase
- Follows Go community best practices

## 🛠️ Available Commands

Using the included Makefile:

```bash
# Development
make run          # Start the server
make dev          # Start with hot reload (requires air)
make deps         # Download dependencies
make clean        # Clean build artifacts

# Testing
make test         # Run all tests
make test-coverage # Run tests with coverage report

# Code Quality
make fmt          # Format code
make lint         # Lint code (requires golangci-lint)

# Docker
make docker-build        # Build Docker image
make docker-compose-up   # Start with Docker Compose
make docker-compose-down # Stop Docker services
```

## 📊 Database Design

### Tables Structure

**Categories Table:**
```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(7) NOT NULL DEFAULT '#3B82F6',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft delete support
);
```

**Todos Table:**
```sql
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    priority VARCHAR(10) CHECK (priority IN ('low', 'medium', 'high')),
    due_date TIMESTAMP WITH TIME ZONE,
    category_id INTEGER REFERENCES categories(id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

### Performance Indexes
Strategic indexes for optimal query performance:
- `idx_todos_completed` - Filter by completion status
- `idx_todos_category_id` - Filter by category
- `idx_todos_priority` - Filter by priority
- `idx_todos_created_at` - Default sorting
- `idx_todos_title_search` - Full-text search using GIN index
- Composite indexes for common query patterns

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
Currently, no authentication is required. All endpoints are publicly accessible.

---

## 📝 Todos API

### GET /api/v1/todos
Retrieve paginated list of todos with optional filtering and sorting.

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number (starts from 1) |
| `limit` | integer | 10 | Items per page (max 50) |
| `search` | string | - | Search in todo titles |
| `completed` | boolean | - | Filter by completion status |
| `category_id` | integer | - | Filter by category ID |
| `priority` | string | - | Filter by priority (low, medium, high) |
| `sort_by` | string | created_at | Sort field (created_at, due_date, title) |
| `sort_order` | string | desc | Sort direction (asc, desc) |

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/todos?page=1&limit=5&completed=false&priority=high"
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Complete coding challenge",
      "description": "Build a full-stack todo application for Industrix",
      "completed": false,
      "priority": "high",
      "due_date": "2024-08-03T23:59:59Z",
      "category_id": 1,
      "category": {
        "id": 1,
        "name": "Work",
        "color": "#3B82F6",
        "created_at": "2024-07-31T10:00:00Z",
        "updated_at": "2024-07-31T10:00:00Z"
      },
      "created_at": "2024-07-31T10:00:00Z",
      "updated_at": "2024-07-31T10:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 5,
    "total": 25,
    "total_pages": 5
  }
}
```

### POST /api/v1/todos
Create a new todo item.

**Request Body:**
```json
{
  "title": "New todo item",
  "description": "Optional detailed description",
  "priority": "medium",
  "due_date": "2024-08-10T15:30:00Z",
  "category_id": 1
}
```

**Validation Rules:**
- `title`: Required, 1-255 characters
- `description`: Optional, max 1000 characters
- `priority`: Optional, must be "low", "medium", or "high"
- `due_date`: Optional, must be valid ISO 8601 timestamp
- `category_id`: Optional, must reference existing category

**Response (201 Created):**
```json
{
  "data": {
    "id": 15,
    "title": "New todo item",
    "description": "Optional detailed description",
    "completed": false,
    "priority": "medium",
    "due_date": "2024-08-10T15:30:00Z",
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Work",
      "color": "#3B82F6"
    },
    "created_at": "2024-08-01T14:20:30Z",
    "updated_at": "2024-08-01T14:20:30Z"
  }
}
```

### GET /api/v1/todos/:id
Get a specific todo by ID.

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/todos/1"
```

**Response (200 OK):**
```json
{
  "data": {
    "id": 1,
    "title": "Complete coding challenge",
    "description": "Build a full-stack todo application",
    "completed": false,
    "priority": "high",
    "due_date": "2024-08-03T23:59:59Z",
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Work", 
      "color": "#3B82F6"
    },
    "created_at": "2024-07-31T10:00:00Z",
    "updated_at": "2024-07-31T10:00:00Z"
  }
}
```

### PUT /api/v1/todos/:id
Update an existing todo item.

**Request Body:**
```json
{
  "title": "Updated todo title",
  "description": "Updated description", 
  "priority": "low",
  "due_date": "2024-08-05T12:00:00Z",
  "category_id": 2
}
```

**Response (200 OK):**
```json
{
  "data": {
    "id": 1,
    "title": "Updated todo title",
    "description": "Updated description",
    "completed": false,
    "priority": "low",
    "due_date": "2024-08-05T12:00:00Z",
    "category_id": 2,
    "category": {
      "id": 2,
      "name": "Personal",
      "color": "#10B981"
    },
    "created_at": "2024-07-31T10:00:00Z",
    "updated_at": "2024-08-01T14:25:15Z"
  }
}
```

### PATCH /api/v1/todos/:id/toggle
Toggle the completion status of a todo.

**Example Request:**
```bash
curl -X PATCH "http://localhost:8080/api/v1/todos/1/toggle"
```

**Response (200 OK):**
```json
{
  "message": "Todo completion status updated successfully",
  "data": {
    "completed": true
  }
}
```

### DELETE /api/v1/todos/:id
Delete a todo item (soft delete).

**Example Request:**
```bash
curl -X DELETE "http://localhost:8080/api/v1/todos/1"
```

**Response (200 OK):**
```json
{
  "message": "Todo deleted successfully"
}
```

---

## 🏷️ Categories API

### GET /api/v1/categories
Get all categories.

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/categories"
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Work",
      "color": "#3B82F6",
      "created_at": "2024-07-31T10:00:00Z",
      "updated_at": "2024-07-31T10:00:00Z"
    },
    {
      "id": 2, 
      "name": "Personal",
      "color": "#10B981",
      "created_at": "2024-07-31T10:00:00Z",
      "updated_at": "2024-07-31T10:00:00Z"
    }
  ]
}
```

### POST /api/v1/categories
Create a new category.

**Request Body:**
```json
{
  "name": "Shopping",
  "color": "#F59E0B"
}
```

**Validation Rules:**
- `name`: Required, 1-100 characters, must be unique
- `color`: Required, valid hex color code (e.g., "#FF0000")

**Response (201 Created):**
```json
{
  "data": {
    "id": 3,
    "name": "Shopping", 
    "color": "#F59E0B",
    "created_at": "2024-08-01T14:30:00Z",
    "updated_at": "2024-08-01T14:30:00Z"
  }
}
```

### PUT /api/v1/categories/:id
Update an existing category.

**Request Body:**
```json
{
  "name": "Work Projects",
  "color": "#6366F1"
}
```

**Response (200 OK):**
```json
{
  "data": {
    "id": 1,
    "name": "Work Projects",
    "color": "#6366F1", 
    "created_at": "2024-07-31T10:00:00Z",
    "updated_at": "2024-08-01T14:35:20Z"
  }
}
```

### DELETE /api/v1/categories/:id
Delete a category (soft delete). Associated todos will have their category_id set to null.

**Example Request:**
```bash
curl -X DELETE "http://localhost:8080/api/v1/categories/1"
```

**Response (200 OK):**
```json
{
  "message": "Category deleted successfully"
}
```

---

## ❌ Error Responses

All error responses follow a consistent format:

**Validation Error (400 Bad Request):**
```json
{
  "error": "Validation failed",
  "type": "validation_error",
  "details": [
    {
      "field": "title",
      "message": "Title is required"
    }
  ]
}
```

**Not Found (404):**
```json
{
  "error": "Todo not found",
  "type": "not_found"
}
```

**Server Error (500):**
```json
{
  "error": "Internal server error",
  "type": "server_error"
}
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file based on `.env.example`:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=todoapp
DB_SSLMODE=disable

# Server Configuration
PORT=8080
GIN_MODE=debug

# CORS Configuration (comma-separated origins)
ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
```

### Database Migration

Database migrations run automatically on server startup. Migration files are located in `migrations/`:

- `001_create_categories_table.sql` - Creates categories table with default data
- `002_create_todos_table.sql` - Creates todos table with indexes

## 🐳 Docker Support

### Using Docker Compose (Recommended)

Start the entire stack (database + API):
```bash
make docker-compose-up
# OR
docker-compose up -d
```

This starts:
- PostgreSQL database on port 5432
- API server on port 8080

### Using Docker Only

Build and run just the API:
```bash
make docker-build
make docker-run
# OR
docker build -t todo-backend .
docker run -p 8080:8080 --env-file .env todo-backend
```

## 🧪 Testing

### Running Tests
```bash
make test           # Run all tests
make test-coverage  # Run with coverage report
```

### Test Structure
```
internal/
├── handlers/
│   └── *_test.go    # HTTP handler tests
├── services/  
│   └── *_test.go    # Business logic tests
└── repository/
    └── *_test.go    # Data access tests
```

Tests cover:
- API endpoint behavior
- Business logic validation
- Error handling scenarios
- Database operations

## 🚀 Development

### Hot Reload
Install Air for hot reload during development:
```bash
go install github.com/cosmtrek/air@latest
make dev
```

### Code Quality
```bash
make fmt   # Format code
make lint  # Lint code (requires golangci-lint)
```

### Adding New Features

1. **Add Model**: Define in `internal/models/`
2. **Add Repository**: Implement in `internal/repository/`
3. **Add Service**: Business logic in `internal/services/`
4. **Add Handler**: HTTP logic in `internal/handlers/`
5. **Add Routes**: Register in `internal/handlers/routes.go`
6. **Add Migration**: Create SQL migration file

## 🏆 Key Features

- **Clean Architecture**: Layered design with clear separation of concerns
- **Type Safety**: Comprehensive struct validation with GORM tags
- **Performance**: Strategic database indexing and efficient pagination
- **Error Handling**: Consistent error responses with proper HTTP status codes
- **CORS Support**: Configurable cross-origin resource sharing
- **Soft Deletes**: Data preservation with GORM soft delete functionality
- **Database Migrations**: Automatic schema management
- **Docker Ready**: Complete containerization support
- **Hot Reload**: Development-friendly with Air integration

This backend provides a solid foundation for the todo application with room for easy expansion and modification.