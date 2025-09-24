import React, { createContext, useContext, useReducer, useCallback, ReactNode } from 'react';
import { Todo, PaginationInfo, TodoFilters, TodoFormData } from '../types';
import { todoApi } from '../services/api';
import { message } from 'antd';

interface TodoState {
  todos: Todo[];
  loading: boolean;
  error: string | null;
  pagination: PaginationInfo | null;
  filters: TodoFilters;
}

type TodoAction =
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_TODOS'; payload: { todos: Todo[]; pagination?: PaginationInfo } }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'ADD_TODO'; payload: Todo }
  | { type: 'UPDATE_TODO'; payload: Todo }
  | { type: 'DELETE_TODO'; payload: number }
  | { type: 'SET_FILTERS'; payload: Partial<TodoFilters> }
  | { type: 'CLEAR_FILTERS' };

const initialState: TodoState = {
  todos: [],
  loading: false,
  error: null,
  pagination: null,
  filters: {
    page: 1,
    per_page: 10,
  },
};

const todoReducer = (state: TodoState, action: TodoAction): TodoState => {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, loading: action.payload };
    case 'SET_TODOS':
      return {
        ...state,
        todos: action.payload.todos,
        pagination: action.payload.pagination || null,
        loading: false,
        error: null,
      };
    case 'SET_ERROR':
      return { ...state, error: action.payload, loading: false };
    case 'ADD_TODO':
      return { ...state, todos: [action.payload, ...state.todos] };
    case 'UPDATE_TODO':
      if (!action.payload || !action.payload.id) {
        console.error('UPDATE_TODO action payload is invalid:', action.payload);
        return state;
      }
      return {
        ...state,
        todos: state.todos.map(todo =>
          todo.id === action.payload.id ? action.payload : todo
        ),
      };
    case 'DELETE_TODO':
      return {
        ...state,
        todos: state.todos.filter(todo => todo.id !== action.payload),
      };
    case 'SET_FILTERS':
      return {
        ...state,
        filters: { ...state.filters, ...action.payload },
      };
    case 'CLEAR_FILTERS':
      return {
        ...state,
        filters: { page: 1, per_page: 10 },
      };
    default:
      return state;
  }
};

interface TodoContextType {
  state: TodoState;
  fetchTodos: () => Promise<void>;
  createTodo: (todoData: TodoFormData) => Promise<boolean>;
  updateTodo: (id: number, todoData: Partial<TodoFormData>) => Promise<boolean>;
  deleteTodo: (id: number) => Promise<boolean>;
  toggleComplete: (id: number) => Promise<boolean>;
  setFilters: (filters: Partial<TodoFilters>) => void;
  clearFilters: () => void;
}

const TodoContext = createContext<TodoContextType | undefined>(undefined);

export const useTodos = () => {
  const context = useContext(TodoContext);
  if (!context) {
    throw new Error('useTodos must be used within a TodoProvider');
  }
  return context;
};

interface TodoProviderProps {
  children: ReactNode;
}

export const TodoProvider: React.FC<TodoProviderProps> = ({ children }) => {
  const [state, dispatch] = useReducer(todoReducer, initialState);

  const fetchTodos = useCallback(async () => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const response = await todoApi.getTodos(state.filters);
      dispatch({
        type: 'SET_TODOS',
        payload: {
          todos: response.data,
          pagination: response.pagination,
        },
      });
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch todos';
      dispatch({ type: 'SET_ERROR', payload: errorMessage });
      message.error(errorMessage);
    }
  }, [state.filters]);

  const createTodo = useCallback(async (todoData: TodoFormData): Promise<boolean> => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const response = await todoApi.createTodo(todoData);
      dispatch({ type: 'ADD_TODO', payload: response.data });
      message.success('Todo created successfully!');
      return true;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create todo';
      dispatch({ type: 'SET_ERROR', payload: errorMessage });
      message.error(errorMessage);
      return false;
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  }, []);

  const updateTodo = useCallback(async (id: number, todoData: Partial<TodoFormData>): Promise<boolean> => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const response = await todoApi.updateTodo(id, todoData);
      dispatch({ type: 'UPDATE_TODO', payload: response.data });
      message.success('Todo updated successfully!');
      return true;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update todo';
      dispatch({ type: 'SET_ERROR', payload: errorMessage });
      message.error(errorMessage);
      return false;
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  }, []);

  const deleteTodo = useCallback(async (id: number): Promise<boolean> => {
    try {
      await todoApi.deleteTodo(id);
      dispatch({ type: 'DELETE_TODO', payload: id });
      message.success('Todo deleted successfully!');
      return true;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to delete todo';
      message.error(errorMessage);
      return false;
    }
  }, []);

  const toggleComplete = useCallback(async (id: number): Promise<boolean> => {
    try {
      // Optimistically update the local state first
      const todoToUpdate = state.todos.find(todo => todo.id === id);
      if (!todoToUpdate) {
        message.error('Todo not found');
        return false;
      }

      // Update local state optimistically
      const updatedTodo = { ...todoToUpdate, completed: !todoToUpdate.completed };
      dispatch({ type: 'UPDATE_TODO', payload: updatedTodo });

      // Call the API
      const response = await todoApi.toggleComplete(id);
      console.log('Toggle complete response:', response);
      
      if (response.success) {
        // If the API returns data, use it; otherwise keep the optimistic update
        if (response.data) {
          dispatch({ type: 'UPDATE_TODO', payload: response.data });
        }
        return true;
      } else {
        // Revert the optimistic update on API failure
        dispatch({ type: 'UPDATE_TODO', payload: todoToUpdate });
        message.error('Failed to update todo status');
        return false;
      }
    } catch (error) {
      console.error('Toggle complete error:', error);
      
      // Revert the optimistic update on error
      const originalTodo = state.todos.find(todo => todo.id === id);
      if (originalTodo) {
        dispatch({ type: 'UPDATE_TODO', payload: originalTodo });
      }
      
      const errorMessage = error instanceof Error ? error.message : 'Failed to toggle todo status';
      message.error(errorMessage);
      return false;
    }
  }, [state.todos]);

  const setFilters = useCallback((filters: Partial<TodoFilters>) => {
    dispatch({ type: 'SET_FILTERS', payload: filters });
  }, []);

  const clearFilters = useCallback(() => {
    dispatch({ type: 'CLEAR_FILTERS' });
  }, []);

  const value: TodoContextType = {
    state,
    fetchTodos,
    createTodo,
    updateTodo,
    deleteTodo,
    toggleComplete,
    setFilters,
    clearFilters,
  };

  return (
    <TodoContext.Provider value={value}>
      {children}
    </TodoContext.Provider>
  );
};