import axios from 'axios';
import { Todo, Category, ApiResponse, TodoFormData, CategoryFormData, TodoFilters } from '../types';

const API_BASE_URL = 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Todo API functions
export const todoApi = {
  // Get all todos with optional filters
  getTodos: async (filters?: TodoFilters): Promise<ApiResponse<Todo[]>> => {
    const params = new URLSearchParams();
    
    if (filters?.search) params.append('search', filters.search);
    if (filters?.completed !== undefined) params.append('completed', filters.completed.toString());
    if (filters?.category_id) params.append('category_id', filters.category_id.toString());
    if (filters?.priority) params.append('priority', filters.priority);
    if (filters?.page) params.append('page', filters.page.toString());
    if (filters?.per_page) params.append('per_page', filters.per_page.toString());

    const response = await api.get(`/todos?${params.toString()}`);
    return response.data;
  },

  // Get a single todo by ID
  getTodo: async (id: number): Promise<ApiResponse<Todo>> => {
    const response = await api.get(`/todos/${id}`);
    return response.data;
  },

  // Create a new todo
  createTodo: async (todoData: TodoFormData): Promise<ApiResponse<Todo>> => {
    const response = await api.post('/todos', todoData);
    return response.data;
  },

  // Update an existing todo
  updateTodo: async (id: number, todoData: Partial<TodoFormData>): Promise<ApiResponse<Todo>> => {
    const response = await api.put(`/todos/${id}`, todoData);
    return response.data;
  },

  // Delete a todo
  deleteTodo: async (id: number): Promise<ApiResponse<null>> => {
    const response = await api.delete(`/todos/${id}`);
    return response.data;
  },

  // Toggle todo completion status
  toggleComplete: async (id: number): Promise<ApiResponse<Todo>> => {
    const response = await api.patch(`/todos/${id}/complete`);
    return response.data;
  },
};

// Category API functions
export const categoryApi = {
  // Get all categories
  getCategories: async (): Promise<ApiResponse<Category[]>> => {
    const response = await api.get('/categories');
    return response.data;
  },

  // Get all categories without pagination
  getAllCategories: async (): Promise<ApiResponse<Category[]>> => {
    const response = await api.get('/categories/all');
    return response.data;
  },

  // Get a single category by ID
  getCategory: async (id: number): Promise<ApiResponse<Category>> => {
    const response = await api.get(`/categories/${id}`);
    return response.data;
  },

  // Create a new category
  createCategory: async (categoryData: CategoryFormData): Promise<ApiResponse<Category>> => {
    const response = await api.post('/categories', categoryData);
    return response.data;
  },

  // Update an existing category
  updateCategory: async (id: number, categoryData: Partial<CategoryFormData>): Promise<ApiResponse<Category>> => {
    const response = await api.put(`/categories/${id}`, categoryData);
    return response.data;
  },

  // Delete a category
  deleteCategory: async (id: number): Promise<ApiResponse<null>> => {
    const response = await api.delete(`/categories/${id}`);
    return response.data;
  },
};

// Health check
export const healthCheck = async (): Promise<{ status: string; message: string }> => {
  const response = await api.get('/health');
  return response.data;
};

export default api;