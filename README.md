# Full-Stack Todo Application

A modern, responsive todo list application built for the Industrix coding challenge. This project demonstrates full-stack development skills using React with TypeScript for the frontend and Go with PostgreSQL for the backend.

## Project Overview

This todo application provides a comprehensive task management solution with advanced features and clean architecture:

### Core Features
- **Complete Todo Management**: Full CRUD operations for todo items
- **Category Organization**: Custom categories with color coding
- **Priority System**: Three-level priority system with visual indicators  
- **Advanced Filtering**: Filter by completion status, category, and priority
- **Real-time Search**: Search through todo titles with instant results
- **Pagination**: Efficient pagination for large todo lists
- **Due Date Management**: Set and track due dates with overdue indicators
- **Responsive Design**: Mobile-first design that works on all screen sizes

### Technical Achievements
- **Modern React 19**: With TypeScript and functional components
- **React Context API**: Proper state management with useReducer pattern
- **Professional UI**: Ant Design component library
- **RESTful API**: Well-structured Go backend with clean architecture
- **Database Design**: PostgreSQL with proper indexing and relationships
- **Type Safety**: Full TypeScript implementation

## Project Structure

```
fullstack-todo-Industrix/
├── server/          # Go backend - See server/README.md
├── web/            # React frontend - See web/README.md
└── README.md       # This overview file
```

## Quick Start

### Prerequisites
- Node.js (v16+)
- Go (v1.19+) 
- PostgreSQL (v12+)

### Setup Instructions

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd fullstack-todo-Industrix
   ```

2. **Setup Backend:**
   ```bash
   cd server
   # Follow detailed instructions in server/README.md
   ```

3. **Setup Frontend:**
   ```bash
   cd web  
   # Follow detailed instructions in web/README.md
   ```

For detailed setup instructions, see the individual README files in each directory.

## Documentation

- **[Backend API Documentation](./server/README.md)** - Complete API reference, setup, and architecture
- **[Frontend Documentation](./web/README.md)** - Component structure, state management, and features

## Technology Stack

### Frontend
- React 19 + TypeScript
- Ant Design UI Components
- React Context API + useReducer
- Axios for API calls
- Responsive design

### Backend  
- Go 1.21 + Gin framework
- GORM ORM + PostgreSQL
- RESTful API design
- Database migrations
- Structured logging

## Features Completed

### Core Requirements
- Todo CRUD operations
- Category management  
- Pagination
- Search functionality
- Responsive design

### Additional Features
- React Context API implementation
- Advanced filtering capabilities
- TypeScript integration
- Professional UI/UX design

---

## Technical Questions & Answers

### Database Design Questions

#### 1. What database tables did you create and why?

I designed a normalized two-table structure:

**Categories Table:**
- **Purpose**: Stores reusable todo categories for organization
- **Key Fields**: `id`, `name` (unique), `color`, timestamps
- **Why**: Prevents data duplication and allows flexible category management

**Todos Table:** 
- **Purpose**: Main entity storing all todo items with their attributes
- **Key Fields**: `id`, `title`, `description`, `completed`, `priority`, `due_date`, `category_id` (FK)
- **Why**: Contains all todo-specific data with proper foreign key relationships

**Relationship Design:**
- One-to-Many: Categories → Todos
- Foreign Key: `todos.category_id` → `categories.id`
- Referential Integrity: CASCADE updates, SET NULL on category deletion (preserves todos)

This structure prevents data duplication, maintains integrity, and supports efficient querying.

#### 2. How did you handle pagination and filtering in the database?

**Pagination Strategy:**
```go
// Efficient OFFSET/LIMIT with total count
query := db.Model(&Todo{}).Where(conditions)
query.Count(&total) // For pagination metadata
query.Offset((page - 1) * limit).Limit(limit).Find(&todos)
```

**Dynamic Filtering:**
```go
if search != "" {
    db = db.Where("title ILIKE ?", "%"+search+"%")
}
if completed != nil {
    db = db.Where("completed = ?", *completed)
}
// Additional filters for category, priority
```

**Performance Indexes:**
- `idx_todos_completed` - Status filtering
- `idx_todos_category_id` - Category filtering  
- `idx_todos_created_at` - Default sorting
- `idx_todos_title_search` - Full-text search (GIN index)
- Composite indexes for common filter+sort combinations

This approach ensures fast queries even with large datasets.

### Technical Decision Questions

#### 1. How did you implement responsive design?

**Breakpoint Strategy:**
I used Ant Design's responsive breakpoints:
- Mobile: 0-767px (Drawer navigation)
- Tablet: 768-991px (Collapsible sidebar)  
- Desktop: 992px+ (Fixed sidebar)

**Implementation Approach:**
```typescript
const useResponsive = () => {
  const [screenSize, setScreenSize] = useState('lg');
  // Window resize listener logic
  return { isMobile: screenSize === 'sm', isTablet: screenSize === 'md' };
};
```

**Key Components:**
- **Layout**: Ant Design's `Layout.Sider` with responsive collapse
- **Navigation**: Fixed sidebar → collapsible → mobile drawer
- **Touch Optimization**: Larger buttons on mobile, gesture-friendly interactions

The design adapts seamlessly across all device sizes with appropriate touch targets and navigation patterns.

#### 2. How did you structure your React components?

**Component Architecture:**
```
App
├── Providers (TodoProvider, CategoryProvider)  
└── AppLayout
    ├── Header (Title + mobile menu)
    ├── Sidebar/Drawer (Filters + actions)
    ├── Content (TodoList + pagination)
    └── Modals (TodoForm, CategoryForm)
```

**State Management:**
I implemented centralized state with React Context API + useReducer:
```typescript
const TodoProvider = ({ children }) => {
  const [state, dispatch] = useReducer(todoReducer, initialState);
  
  const createTodo = useCallback(async (data) => {
    dispatch({ type: 'SET_LOADING', payload: true });
    // API call + state update
  }, []);
  
  return <TodoContext.Provider value={{state, createTodo}}>;
};
```

**Benefits:**
- Clear separation of concerns
- Predictable state updates
- Easy testing and maintenance
- TypeScript safety throughout

#### 3. What backend architecture did you choose and why?

**Layered Architecture:**
```
Handlers (HTTP) → Services (Business Logic) → Repository (Data Access) → Models
```

**Structure Benefits:**
- **Handlers**: Thin controllers handling HTTP concerns only
- **Services**: All business logic and validation rules  
- **Repository**: Database operations abstracted behind interfaces
- **Models**: GORM models with validation tags

**Example Flow:**
```go
func (h *TodoHandler) CreateTodo(c *gin.Context) {
    var req models.Todo
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    todo, err := h.todoService.Create(&req) // Delegate to service
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return  
    }
    
    c.JSON(201, gin.H{"data": todo})
}
```

This architecture provides excellent testability, maintainability, and follows Go best practices.

#### 4. How did you handle data validation?

**Multi-Layer Validation:**

**Frontend (UX):**
```typescript
const rules = {
  title: [
    { required: true, message: 'Title required' },
    { max: 255, message: 'Max 255 characters' }
  ],
  due_date: [{
    validator: (_, value) => {
      if (value && dayjs(value).isBefore(dayjs())) {
        return Promise.reject('Cannot be in the past');
      }
    }
  }]
};
```

**Backend (Security):**
```go
type Todo struct {
    Title    string   `binding:"required,min=1,max=255"`
    Priority Priority `binding:"omitempty,oneof=low medium high"`
    // Additional validation in service layer
}

func (s *TodoService) validateTodo(todo *Todo) error {
    if todo.DueDate != nil && todo.DueDate.Before(time.Now()) {
        return errors.New("due date cannot be in the past")
    }
    return nil
}
```

**Database (Integrity):**
```sql
CHECK (priority IN ('low', 'medium', 'high'))
```

This three-layer approach ensures data integrity while providing excellent user experience.

### Future Improvements

#### If you had more time, what would you improve or add?

**Performance Improvements:**
- Virtual scrolling for large todo lists
- React.memo optimization for expensive renders  
- Database query optimization and caching

**Advanced Features:**
- Real-time collaboration with WebSockets
- Offline capability with service workers
- Advanced search with Elasticsearch
- User authentication and authorization
- Push notifications for due dates

**Infrastructure:**
- Docker containerization
- CI/CD pipeline setup
- Unit test coverage for critical business logic
- API documentation with Swagger
- Performance monitoring and logging

**Code Quality:**
- Error boundaries in React
- Better loading states and skeleton screens
- Code splitting for performance
- API versioning strategy

These improvements would transform the application into a production-ready, scalable solution suitable for real-world deployment.

---

## Summary

This todo application demonstrates proficiency in modern full-stack development with clean architecture, proper state management, and production-ready code quality. The technical decisions prioritize maintainability, performance, and user experience while following industry best practices.