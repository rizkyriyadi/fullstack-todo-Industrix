export type Priority = 'low' | 'medium' | 'high';

export interface Category {
  id: number;
  name: string;
  color: string;
  created_at: string;
  updated_at: string;
}

export interface Todo {
  id: number;
  title: string;
  description: string;
  completed: boolean;
  priority: Priority;
  due_date?: string;
  category_id?: number;
  created_at: string;
  updated_at: string;
  category?: Category;
}

export interface PaginationInfo {
  current_page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
  pagination?: PaginationInfo;
}

export interface TodoFormData {
  title: string;
  description: string;
  priority: Priority;
  due_date?: string;
  category_id?: number;
}

export interface CategoryFormData {
  name: string;
  color: string;
}

export interface TodoFilters {
  search?: string;
  completed?: boolean;
  category_id?: number;
  priority?: Priority;
  page?: number;
  per_page?: number;
}