# Todo Backend API

A full-stack todo list backend API built with Go, Gin, GORM, and PostgreSQL. This project implements a clean architecture with proper separation of concerns, comprehensive error handling, and robust data validation.

## Features

### Core Features
- ✅ **Todo Management**: Create, read, update, delete todos
- ✅ **Category Management**: Organize todos into categories
- ✅ **Search & Filtering**: Search todos by title, filter by completion status, category, and priority
- ✅ **Pagination**: Efficient pagination with configurable page sizes
- ✅ **Priority Levels**: High, medium, low priority support
- ✅ **Due Dates**: Optional due dates for todos

### Technical Features
- ✅ **Clean Architecture**: Repository pattern, service layer, dependency injection
- ✅ **Database Migrations**: SQL migrations with up/down support
- ✅ **Input Validation**: Comprehensive validation using Gin's binding tags
- ✅ **Error Handling**: Structured error responses with proper HTTP status codes
- ✅ **CORS Support**: Configurable CORS for frontend integration
- ✅ **Logging**: Structured request/response logging
- ✅ **Security Headers**: Basic security headers implementation

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **ORM**: GORM v2
- **Database**: PostgreSQL
- **Configuration**: Environment variables with godotenv
- **Architecture**: Clean Architecture (Repository + Service pattern)

## Project Structure

```
server/
├── cmd/
│   └── api/                    # Application entry point
│       └── main.go
├── internal/
│   ├── config/                 # Configuration management
│   ├── handlers/               # HTTP handlers (controllers)
│   ├── middleware/             # Custom middleware
│   ├── models/                 # Database models
│   ├── repository/             # Data access layer
│   └── services/               # Business logic layer
├── pkg/
│   ├── database/               # Database utilities
│   └── utils/                  # Common utilities
├── migrations/                 # SQL migration files
├── .env.example               # Environment variables template
├── Makefile                   # Development commands
└── README.md                  # This file
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

## Quick Start

### 1. Clone and Setup

```bash
# Clone the repository (if using git)
cd server

# Install dependencies
go mod download
```

### 2. Database Setup

Create a PostgreSQL database:

```sql
CREATE DATABASE todo_db;
```

### 3. Environment Configuration

Copy the example environment file and configure your database:

```bash
cp .env.example .env
```

Edit `.env` with your database credentials:

```env
# Server Configuration
PORT=8080
APP_ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=todo_db
DB_SSL_MODE=disable
DB_TIMEZONE=UTC
```

### 4. Run the Application

```bash
# Using Go directly
go run cmd/api/main.go

# OR using Makefile
make run

# OR for development with hot reload (requires air)
make dev
```

The API will be available at `http://localhost:8080`

### 5. Verify Installation

Test the health endpoint:

```bash
curl http://localhost:8080/api/health
```

Expected response:
```json
{
  "status": "OK",
  "message": "Todo API is running"
}
```

## API Documentation

### Base URL
```
http://localhost:8080/api
```

### Authentication
Currently, no authentication is required (this would be a future enhancement).

### Response Format

All API responses follow this structure:

**Success Response:**
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { /* response data */ }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

**Paginated Response:**
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [ /* array of items */ ],
  "pagination": {
    "current_page": 1,
    "per_page": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

### Todo Endpoints

#### Create Todo
```http
POST /api/todos
Content-Type: application/json

{
  "title": "Complete coding challenge",
  "description": "Build a full-stack todo application",
  "priority": "high",
  "due_date": "2024-12-31T23:59:59Z",
  "category_id": 1
}
```

#### List Todos
```http
GET /api/todos?page=1&limit=10&search=coding&completed=false&category_id=1&priority=high&sort_by=created_at&sort_order=desc
```

**Query Parameters:**
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 10, max: 100)
- `search` (string): Search in todo titles
- `completed` (bool): Filter by completion status
- `category_id` (int): Filter by category
- `priority` (string): Filter by priority (low, medium, high)
- `sort_by` (string): Sort field (title, completed, priority, due_date, created_at, updated_at)
- `sort_order` (string): Sort order (asc, desc)

#### Get Todo
```http
GET /api/todos/{id}
```

#### Update Todo
```http
PUT /api/todos/{id}
Content-Type: application/json

{
  "title": "Updated title",
  "description": "Updated description",
  "completed": true,
  "priority": "medium",
  "due_date": "2024-12-31T23:59:59Z",
  "category_id": 2
}
```

#### Delete Todo
```http
DELETE /api/todos/{id}
```

#### Toggle Todo Completion
```http
PATCH /api/todos/{id}/complete
```

### Category Endpoints

#### Create Category
```http
POST /api/categories
Content-Type: application/json

{
  "name": "Work",
  "color": "#3B82F6"
}
```

#### List Categories
```http
GET /api/categories?page=1&limit=10&search=work&sort_by=name&sort_order=asc
```

#### Get All Categories (for dropdowns)
```http
GET /api/categories/all
```

#### Get Category
```http
GET /api/categories/{id}
```

#### Update Category
```http
PUT /api/categories/{id}
Content-Type: application/json

{
  "name": "Updated Work",
  "color": "#FF0000"
}
```

#### Delete Category
```http
DELETE /api/categories/{id}
```

### HTTP Status Codes

- `200 OK`: Success
- `201 Created`: Resource created successfully
- `400 Bad Request`: Validation error or malformed request
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., duplicate name)
- `500 Internal Server Error`: Server error

## Data Models

### Todo Model
```json
{
  "id": 1,
  "title": "Complete coding challenge",
  "description": "Build a full-stack todo application",
  "completed": false,
  "priority": "high",
  "due_date": "2024-12-31T23:59:59Z",
  "category_id": 1,
  "category": {
    "id": 1,
    "name": "Work",
    "color": "#3B82F6"
  },
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

### Category Model
```json
{
  "id": 1,
  "name": "Work",
  "color": "#3B82F6",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

## Development

### Available Commands

```bash
# Run the application
make run

# Run with hot reload (requires air)
make dev

# Build the application
make build

# Download dependencies
make deps

# Format code
make fmt

# Clean build artifacts
make clean

# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Installing Air for Hot Reload

```bash
go install github.com/cosmtrek/air@latest
```

### Database Migrations

The application automatically runs migrations in development mode. For production, you should run migrations manually using the provided SQL files in the `migrations/` directory.

## Testing

Run the test suite:

```bash
make test
```

Run tests with coverage:

```bash
make test-coverage
```

This will generate a `coverage.html` file you can open in your browser.

## Deployment

### Environment Variables for Production

For production deployment, set these environment variables:

```env
APP_ENV=production
PORT=8080
DB_HOST=your_prod_db_host
DB_USER=your_prod_db_user
DB_PASSWORD=your_prod_db_password
DB_NAME=your_prod_db_name
DB_SSL_MODE=require
```

### Building for Production

```bash
make build
```

This creates a binary in `bin/todo-api` that you can deploy.

## Troubleshooting

### Common Issues

1. **Database Connection Error**
   - Verify PostgreSQL is running
   - Check database credentials in `.env`
   - Ensure database exists

2. **Port Already in Use**
   - Change the `PORT` in `.env`
   - Kill any process using port 8080: `lsof -ti:8080 | xargs kill`

3. **Migration Errors**
   - Check database user permissions
   - Verify database exists
   - Check migration SQL syntax

### Logs

The application logs all requests and errors. Check the console output for debugging information.

## Contributing

1. Follow Go conventions and best practices
2. Write tests for new features
3. Update documentation for API changes
4. Use meaningful commit messages

## License

This project is created for the Industrix coding challenge.