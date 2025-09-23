# Technical Questions & Answers

This document answers the required technical questions from the challenge requirements.

## Database Design Questions

### 1. What database tables did you create and why?

#### Tables Created:

**categories table:**
- **Purpose**: Organizes todos into different groups (Work, Personal, Shopping, Health, etc.)
- **Structure**: 
  - `id` (Primary Key): Unique identifier for each category
  - `name` (VARCHAR(100), UNIQUE): Category name with uniqueness constraint
  - `color` (VARCHAR(7)): Hex color code for UI representation
  - `created_at`, `updated_at`: Timestamps for audit trail
  - `deleted_at`: Soft delete support
- **Why**: Categories provide organization and filtering capabilities, making it easier for users to manage large numbers of todos

**todos table:**
- **Purpose**: Stores the main todo items that users create and manage
- **Structure**:
  - `id` (Primary Key): Unique identifier for each todo
  - `title` (VARCHAR(255), NOT NULL): Brief description of the task
  - `description` (TEXT): Detailed information about the todo
  - `completed` (BOOLEAN): Completion status
  - `priority` (VARCHAR(10)): Priority level (low, medium, high) with CHECK constraint
  - `due_date` (TIMESTAMP): Optional deadline
  - `category_id` (Foreign Key): Reference to categories table with CASCADE behavior
  - `created_at`, `updated_at`: Timestamps for audit trail
  - `deleted_at`: Soft delete support

#### Relationships:
- **One-to-Many**: Category → Todos (one category can have many todos)
- **Foreign Key Constraint**: `todos.category_id` references `categories.id` with `ON DELETE SET NULL`
- **Why this structure**: Maintains referential integrity while allowing todos to exist without categories

### 2. How did you handle pagination and filtering in the database?

#### Pagination Implementation:
- **LIMIT/OFFSET approach**: Used PostgreSQL's LIMIT and OFFSET for efficient pagination
- **Parameters**: 
  - `page`: Current page number (1-based)
  - `limit`: Number of items per page (default: 10, max: 100)
- **Calculation**: `OFFSET = (page - 1) * limit`
- **Total Count**: Separate COUNT query to determine total pages

#### Filtering Implementation:
**Todos Filtering:**
- **Search**: Full-text search using `ILIKE` on title field
- **Completion Status**: Boolean filter on `completed` field
- **Category**: Filter by `category_id`
- **Priority**: Enum filter on priority field
- **Soft Delete**: All queries include `WHERE deleted_at IS NULL`

**Categories Filtering:**
- **Search**: Case-insensitive search on category name using `ILIKE`

#### Indexes Added:
```sql
-- Performance indexes
CREATE INDEX idx_todos_completed ON todos(completed) WHERE deleted_at IS NULL;
CREATE INDEX idx_todos_priority ON todos(priority) WHERE deleted_at IS NULL;
CREATE INDEX idx_todos_category_id ON todos(category_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_todos_created_at ON todos(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_todos_title_search ON todos USING gin(to_tsvector('english', title));

-- Composite index for common filtering
CREATE INDEX idx_todos_completed_category ON todos(completed, category_id) WHERE deleted_at IS NULL;
```

**Why these indexes**: 
- Single-column indexes for common filters
- Composite index for frequently combined filters
- GIN index for full-text search on titles
- Partial indexes to exclude soft-deleted records

## Technical Decision Questions

### 1. How did you implement responsive design?

**Note**: This backend API doesn't implement responsive design directly, as that's a frontend concern. However, the API is designed to support responsive frontends through:

- **Flexible Pagination**: Configurable page sizes allow frontends to request different amounts of data based on screen size
- **Efficient Data Transfer**: Only necessary fields are returned to minimize bandwidth usage
- **Search & Filtering**: Rich filtering options allow frontends to implement efficient mobile-friendly interfaces
- **RESTful Design**: Standard REST endpoints make integration straightforward across different devices

### 2. How did you structure your React components?

**Note**: This is the backend implementation. The frontend structure would be covered in the React application.

### 3. What backend architecture did you choose and why?

#### Architecture: Clean Architecture with Repository Pattern

**Layers:**
1. **Handlers Layer** (`internal/handlers/`): HTTP request/response handling
2. **Services Layer** (`internal/services/`): Business logic and validation
3. **Repository Layer** (`internal/repository/`): Data access abstraction
4. **Models Layer** (`internal/models/`): Database entities and business objects

#### Code Organization:
```
internal/
├── handlers/          # HTTP handlers (controllers)
│   ├── todo_handler.go
│   ├── category_handler.go
│   └── routes.go
├── services/          # Business logic
│   ├── todo_service.go
│   └── category_service.go
├── repository/        # Data access
│   ├── todo_repository.go
│   └── category_repository.go
├── models/           # Data models
│   ├── todo.go
│   └── category.go
├── middleware/       # HTTP middleware
└── config/          # Configuration
```

#### API Routes Structure:
- **RESTful Design**: Standard HTTP methods (GET, POST, PUT, DELETE, PATCH)
- **Resource-based URLs**: `/api/todos`, `/api/categories`
- **Nested Resources**: `/api/todos/:id/complete`
- **Versioning**: `/api` prefix for future versioning

#### Error Handling Approach:
- **Centralized Error Handling**: Custom middleware for panic recovery
- **Structured Error Responses**: Consistent JSON error format
- **HTTP Status Codes**: Proper status codes for different error types
- **Validation Errors**: Detailed validation error messages
- **Logging**: Comprehensive request/error logging

**Why this architecture**:
- **Separation of Concerns**: Each layer has a single responsibility
- **Testability**: Easy to unit test individual components
- **Maintainability**: Clear dependencies and interfaces
- **Scalability**: Easy to modify or replace individual layers

### 4. How did you handle data validation?

#### Validation Layers:

**1. Struct-level Validation (Model Layer):**
```go
type Todo struct {
    Title    string `binding:"required,min=1,max=255"`
    Priority Priority `binding:"omitempty,oneof=low medium high"`
    // ...
}
```

**2. Business Logic Validation (Service Layer):**
- Due date validation (not in the past)
- Category existence validation
- Data cleanup and normalization
- Priority enum validation

**3. Database Constraints:**
- NOT NULL constraints
- CHECK constraints for enums
- UNIQUE constraints
- Foreign key constraints

#### Where Validation Occurs:
- **Frontend**: Basic client-side validation (would be implemented in React)
- **Backend API**: Comprehensive server-side validation
- **Database**: Final constraints and data integrity

#### Validation Rules Implemented:
**Todos:**
- Title: Required, 1-255 characters
- Description: Optional, max 5000 characters
- Priority: Must be 'low', 'medium', or 'high'
- Due Date: Cannot be in the past
- Category: Must exist if specified

**Categories:**
- Name: Required, unique, max 100 characters
- Color: Valid hex color format

**Why this approach**:
- **Multiple layers** prevent invalid data at different stages
- **Client-side validation** provides immediate feedback
- **Server-side validation** ensures data integrity
- **Database constraints** provide final safety net

## Testing & Quality Questions

### 1. What did you choose to unit test and why?

#### Testing Strategy (Recommended Implementation):

**Service Layer Tests** (Highest Priority):
- Business logic validation
- Edge cases and error handling
- Data transformation logic
- Integration between services and repositories

**Repository Layer Tests**:
- Database operations
- Query correctness
- Transaction handling
- Error scenarios

**Handler Tests**:
- HTTP request/response handling
- Input validation
- Status code correctness
- JSON serialization

#### Test Structure:
```go
func TestTodoService_CreateTodo(t *testing.T) {
    tests := []struct {
        name    string
        todo    *models.Todo
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid todo",
            todo: &models.Todo{Title: "Test", Priority: "medium"},
            wantErr: false,
        },
        {
            name: "empty title",
            todo: &models.Todo{Title: "", Priority: "medium"},
            wantErr: true,
            errMsg: "title is required",
        },
        // More test cases...
    }
}
```

#### Edge Cases Considered:
- Empty/invalid input data
- Database connection failures
- Concurrent access scenarios
- Large dataset performance
- Boundary value testing

**Why focus on these areas**:
- **Service layer** contains most business logic
- **Repository layer** handles critical data operations
- **Edge cases** reveal potential production issues
- **Integration tests** ensure components work together

### 2. If you had more time, what would you improve or add?

#### Technical Debt to Address:
1. **Comprehensive Unit Tests**: Add full test coverage for all layers
2. **Integration Tests**: Test database interactions end-to-end
3. **API Documentation**: Add OpenAPI/Swagger documentation
4. **Rate Limiting**: Implement proper rate limiting middleware
5. **Authentication**: Add JWT-based authentication system
6. **Caching**: Implement Redis caching for frequently accessed data

#### Features to Add:
1. **Advanced Filtering**: Date range filters, multiple category selection
2. **Bulk Operations**: Bulk update/delete operations
3. **File Attachments**: Support for file uploads on todos
4. **Notifications**: Due date reminder system
5. **Audit Trail**: Track all changes to todos and categories
6. **Data Export**: Export todos to CSV/PDF
7. **Search Enhancement**: Full-text search across all fields
8. **Real-time Updates**: WebSocket support for live updates

#### Infrastructure Improvements:
1. **Docker Support**: Complete containerization with docker-compose
2. **CI/CD Pipeline**: Automated testing and deployment
3. **Monitoring**: Application metrics and health checks
4. **Database Migrations**: Proper migration management tool
5. **Environment Management**: Support for multiple environments
6. **Performance Optimization**: Query optimization and indexing
7. **Security Hardening**: Input sanitization, SQL injection prevention
8. **API Versioning**: Support for multiple API versions

#### Refactoring Opportunities:
1. **Error Handling**: More granular error types and handling
2. **Configuration**: More flexible configuration management
3. **Logging**: Structured logging with log levels
4. **Validation**: Custom validation rules and messages
5. **Code Organization**: Further modularization of large files

This comprehensive backend provides a solid foundation for the todo application with room for growth and enhancement.