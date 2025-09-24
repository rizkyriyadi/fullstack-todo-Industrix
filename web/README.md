# Todo Frontend Application

A modern, responsive React application built with TypeScript and Ant Design. This frontend provides an intuitive interface for managing todos with advanced filtering, real-time search, and mobile-first design.

## Quick Start

### Prerequisites
- Node.js 16 or higher
- npm or yarn package manager
- Backend API running on http://localhost:8080

### Setup Instructions

1. **Navigate to web directory:**
   ```bash
   cd web
   ```

2. **Install dependencies:**
   ```bash
   npm install
   # OR
   yarn install
   ```

3. **Start development server:**
   ```bash
   npm start
   # OR
   yarn start
   ```

4. **Open in browser:**
   Navigate to `http://localhost:3000`

### Build for Production
```bash
npm run build     # Creates optimized production build
npm run serve     # Serve production build locally
```

## Project Structure

```
web/
├── public/                 # Static assets
├── src/
│   ├── components/        # React components
│   │   ├── layout/       # Layout components
│   │   ├── todo/         # Todo-related components
│   │   └── category/     # Category components
│   ├── contexts/         # React Context providers
│   │   ├── TodoContext.tsx
│   │   └── CategoryContext.tsx
│   ├── services/         # API service layer
│   │   └── api.ts
│   ├── types/            # TypeScript type definitions
│   │   └── index.ts
│   ├── utils/            # Utility functions
│   │   ├── constants.ts
│   │   └── helpers.ts
│   ├── App.tsx           # Main application component
│   └── index.tsx         # Application entry point
├── package.json
└── README.md
```

## Features Overview

### Core Features
- **Complete Todo Management**: Create, edit, delete, and toggle todos
- **Category Organization**: Color-coded categories for better organization
- **Priority System**: Visual priority indicators (high, medium, low)
- **Real-time Search**: Instant search with debouncing
- **Advanced Filtering**: Filter by status, category, and priority
- **Pagination**: Efficient navigation through large todo lists
- **Due Date Management**: Set and track due dates with overdue indicators
- **Responsive Design**: Mobile-first design that works on all screen sizes

### UI/UX Enhancements
- **Modern Interface**: Clean, professional design with Ant Design
- **Loading States**: Skeleton loading and progress indicators
- **Error Handling**: User-friendly error messages and recovery
- **Animations**: Smooth transitions and interactions
- **Accessibility**: Keyboard navigation and screen reader support

## Technology Stack

### Core Technologies
- **React 19** - Latest React with concurrent features and improved performance
- **TypeScript** - Full type safety throughout the application
- **Ant Design 5.x** - Professional UI component library
- **Axios** - HTTP client with interceptors and error handling

### Development Tools
- **Create React App** - Zero-configuration React setup with TypeScript
- **ESLint & Prettier** - Code quality and consistent formatting
- **Day.js** - Lightweight date manipulation library

### State Management
- **React Context API** - Centralized state management
- **useReducer** - Predictable state updates with Redux-like patterns

## Architecture

### Component Architecture

I structured the application using a **feature-based component architecture**:

```
App
├── Providers
│   ├── TodoProvider (Todo state management)
│   └── CategoryProvider (Category state management)
└── AppLayout (Main layout container)
    ├── Header (App title + mobile menu trigger)
    ├── Sidebar/Drawer (Navigation + filters)
    │   ├── ActionButtons (New todo/category buttons)
    │   ├── FilterMenu (Status filter options)
    │   └── CategoryList (Category filters with colors)
    ├── Content (Main content area)
    │   ├── TodoList (Todo display with cards)
    │   │   ├── TodoCard (Individual todo items)
    │   │   └── PaginationControls
    │   ├── TodoForm (Create/edit modal)
    │   └── CategoryForm (Category management modal)
    └── GlobalModals (Application-wide modals)
```

### State Management Strategy

I implemented **centralized state management** using React Context API with useReducer:

```typescript
// TodoContext.tsx - Main state management
const TodoContext = createContext<TodoContextType | undefined>(undefined);

export const TodoProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [state, dispatch] = useReducer(todoReducer, initialState);
  
  // Async actions with error handling
  const createTodo = useCallback(async (todoData: TodoFormData): Promise<boolean> => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const response = await todoApi.createTodo(todoData);
      dispatch({ type: 'ADD_TODO', payload: response.data });
      message.success('Todo created successfully!');
      return true;
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: error.message });
      message.error(error.message);
      return false;
    }
  }, []);

  return (
    <TodoContext.Provider value={{ state, createTodo, ... }}>
      {children}
    </TodoContext.Provider>
  );
};
```

**Benefits of this approach:**
- Predictable state updates
- Easy debugging with Redux DevTools
- Type-safe state management
- Centralized error handling
- Optimistic UI updates

### Responsive Design Implementation

I implemented a **mobile-first responsive design** using Ant Design's breakpoint system:

```typescript
// useResponsive hook
const useResponsive = () => {
  const [screenSize, setScreenSize] = useState<'sm' | 'md' | 'lg'>('lg');
  
  useEffect(() => {
    const updateScreenSize = () => {
      const width = window.innerWidth;
      if (width < 768) setScreenSize('sm');      // Mobile
      else if (width < 992) setScreenSize('md');  // Tablet  
      else setScreenSize('lg');                   // Desktop
    };
    
    updateScreenSize();
    window.addEventListener('resize', updateScreenSize);
    return () => window.removeEventListener('resize', updateScreenSize);
  }, []);
  
  return {
    isMobile: screenSize === 'sm',
    isTablet: screenSize === 'md',
    isDesktop: screenSize === 'lg'
  };
};
```

**Responsive Behavior:**
- **Desktop (≥992px)**: Fixed sidebar with full functionality
- **Tablet (768-991px)**: Collapsible sidebar with icon-only mode
- **Mobile (<768px)**: Hidden sidebar, slide-out drawer navigation

## Components Documentation

### Layout Components

#### AppLayout.tsx
Main layout container that handles responsive behavior and navigation.

**Key Features:**
- Responsive sidebar/drawer switching
- Mobile menu toggle
- Consistent spacing and layout
- Theme integration

#### Header.tsx
Application header with title and mobile menu trigger.

### Todo Components

#### TodoList.tsx
Main todo display component with pagination and empty states.

**Props:**
```typescript
interface TodoListProps {
  loading?: boolean;
  error?: string | null;
}
```

#### TodoCard.tsx
Individual todo item display with actions.

**Features:**
- Priority color indicators
- Due date with overdue highlighting
- Completion toggle with optimistic updates
- Edit/delete action buttons
- Category badges

**Props:**
```typescript
interface TodoCardProps {
  todo: Todo;
  onEdit: (todo: Todo) => void;
  onDelete: (id: number) => void;
  onToggleComplete: (id: number) => void;
}
```

#### TodoForm.tsx
Modal form for creating and editing todos.

**Features:**
- Form validation with Ant Design rules
- Date picker with past date validation
- Category and priority selection
- Loading states during submission

### Category Components

#### CategoryForm.tsx
Modal form for creating and editing categories.

**Features:**
- Color picker for category colors
- Name uniqueness validation
- Delete confirmation dialogs

## State Management

### Todo State Structure
```typescript
interface TodoState {
  todos: Todo[];
  loading: boolean;
  error: string | null;
  filters: FilterState;
  pagination: PaginationState;
}

interface FilterState {
  search: string;
  completed?: boolean;
  categoryId?: number;
  priority?: Priority;
  sortBy: string;
  sortOrder: 'asc' | 'desc';
}

interface PaginationState {
  current: number;
  pageSize: number;
  total: number;
}
```

### Action Types
```typescript
type TodoAction =
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_TODOS'; payload: Todo[] }
  | { type: 'ADD_TODO'; payload: Todo }
  | { type: 'UPDATE_TODO'; payload: Todo }
  | { type: 'DELETE_TODO'; payload: number }
  | { type: 'TOGGLE_TODO'; payload: { id: number; completed: boolean } }
  | { type: 'SET_FILTERS'; payload: Partial<FilterState> }
  | { type: 'SET_PAGINATION'; payload: Partial<PaginationState> }
  | { type: 'SET_ERROR'; payload: string | null };
```

### Custom Hooks

#### useTodos()
Main hook for accessing todo state and actions.

```typescript
const useTodos = () => {
  const context = useContext(TodoContext);
  if (!context) {
    throw new Error('useTodos must be used within TodoProvider');
  }
  return context;
};
```

#### useCategories()
Hook for accessing category state and actions.

#### useResponsive()
Hook for responsive design utilities.

## API Integration

### Service Layer
I created a centralized API service layer for clean separation of concerns:

```typescript
// services/api.ts
const API_BASE_URL = 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

// Request interceptor for consistent headers
api.interceptors.request.use((config) => {
  config.headers['Content-Type'] = 'application/json';
  return config;
});

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const message = error.response?.data?.error || 'An unexpected error occurred';
    return Promise.reject(new Error(message));
  }
);

export const todoApi = {
  getAllTodos: (params: TodoFilters) => api.get<TodoResponse>('/todos', { params }),
  createTodo: (todo: TodoFormData) => api.post<{ data: Todo }>('/todos', todo),
  updateTodo: (id: number, todo: Partial<TodoFormData>) => 
    api.put<{ data: Todo }>(`/todos/${id}`, todo),
  deleteTodo: (id: number) => api.delete(`/todos/${id}`),
  toggleTodo: (id: number) => api.patch(`/todos/${id}/toggle`),
};
```

### Error Handling Strategy
I implemented comprehensive error handling at multiple levels:

1. **API Level**: Axios interceptors for consistent error formatting
2. **Context Level**: Error state management with user-friendly messages
3. **Component Level**: Error boundaries and fallback UIs
4. **User Level**: Toast notifications with Ant Design message component

## Form Validation

### Validation Rules
I implemented comprehensive form validation using Ant Design's form validation:

```typescript
// Todo form validation rules
export const todoValidationRules = {
  title: [
    { required: true, message: 'Title is required' },
    { min: 1, max: 255, message: 'Title must be between 1-255 characters' },
    { whitespace: true, message: 'Title cannot be empty' },
  ],
  description: [
    { max: 1000, message: 'Description must be less than 1000 characters' },
  ],
  priority: [
    {
      type: 'enum' as const,
      enum: ['low', 'medium', 'high'],
      message: 'Please select a valid priority',
    },
  ],
  due_date: [
    {
      validator: (_, value) => {
        if (!value) return Promise.resolve();
        if (dayjs(value).isBefore(dayjs(), 'day')) {
          return Promise.reject('Due date cannot be in the past');
        }
        return Promise.resolve();
      },
    },
  ],
  category_id: [
    { type: 'number' as const, message: 'Please select a category' },
  ],
};
```

## Styling & Theming

### Theme Configuration
I used Ant Design's theme customization for consistent branding:

```typescript
const theme = {
  token: {
    colorPrimary: '#1890ff',
    colorSuccess: '#52c41a',
    colorWarning: '#faad14',
    colorError: '#ff4d4f',
    borderRadius: 6,
    wireframe: false,
  },
  components: {
    Layout: {
      siderBg: '#fff',
      headerBg: '#fff',
    },
    Menu: {
      itemSelectedBg: '#e6f7ff',
    },
  },
};
```

### Priority Color System
```typescript
export const priorityColors = {
  low: '#52c41a',    // Green
  medium: '#faad14', // Orange
  high: '#ff4d4f',   // Red
} as const;

export const getPriorityColor = (priority: Priority): string => {
  return priorityColors[priority] || priorityColors.medium;
};
```

## Responsive Design Details

### Breakpoint Strategy
```typescript
const breakpoints = {
  xs: 0,      // Mobile phones
  sm: 576,    // Large phones, small tablets
  md: 768,    // Tablets
  lg: 992,    // Small laptops
  xl: 1200,   // Desktops
  xxl: 1600,  // Large desktops
};
```

### Mobile Optimizations
- **Touch-friendly**: 44px minimum touch targets
- **Gesture Support**: Swipe to close drawer, pull-to-refresh
- **Simplified Navigation**: Priority-based feature visibility
- **Full-width Forms**: Optimized form layouts for small screens
- **Large Action Buttons**: Better accessibility on mobile devices

## Performance Optimizations

### React Optimizations
```typescript
// Memoized components to prevent unnecessary re-renders
const TodoCard = React.memo(({ todo, onEdit, onDelete, onToggleComplete }) => {
  // Component logic
});

// Memoized expensive calculations  
const filteredTodos = useMemo(() => {
  return todos.filter(todo => {
    if (filters.search && !todo.title.toLowerCase().includes(filters.search.toLowerCase())) {
      return false;
    }
    if (filters.completed !== undefined && todo.completed !== filters.completed) {
      return false;
    }
    // Additional filtering logic
    return true;
  });
}, [todos, filters]);

// Debounced search to reduce API calls
const debouncedSearch = useMemo(
  () => debounce((searchTerm: string) => {
    updateFilters({ search: searchTerm });
  }, 300),
  [updateFilters]
);
```

### Bundle Optimization
- **Code Splitting**: Lazy loading of non-critical components
- **Tree Shaking**: Removing unused code from bundles
- **Asset Optimization**: Optimized images and fonts
- **Caching Strategy**: Service worker for offline capability (future enhancement)

## Development Workflow

### Available Scripts
```bash
npm start          # Start development server
npm run build      # Build for production  
npm run lint       # Lint code
npm run format     # Format code with Prettier
npm run type-check # TypeScript type checking
```

### Code Quality
- **ESLint**: Consistent code style and error prevention
- **Prettier**: Automatic code formatting
- **TypeScript**: Compile-time error catching

### Development Best Practices
- **Component-First**: Build reusable, testable components
- **Type Safety**: Comprehensive TypeScript usage
- **Error Boundaries**: Graceful error handling
- **Performance**: Optimization for large datasets
- **Accessibility**: WCAG compliance and screen reader support

## Key Features Summary

This frontend application demonstrates:

- **Modern React**: Latest React 19 features and patterns
- **Type Safety**: Comprehensive TypeScript implementation
- **State Management**: Proper Context API usage with useReducer
- **Responsive Design**: Mobile-first approach with Ant Design
- **Performance**: Optimized rendering and API calls
- **User Experience**: Intuitive interface with excellent feedback
- **Code Quality**: Clean, maintainable, and well-structured code
- **Error Handling**: Robust error management and recovery
- **Accessibility**: Keyboard navigation and screen reader support

The application provides a professional, production-ready user interface that effectively demonstrates modern frontend development skills and best practices.