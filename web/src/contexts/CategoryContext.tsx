import React, { createContext, useContext, useReducer, useCallback, ReactNode } from 'react';
import { Category, CategoryFormData } from '../types';
import { categoryApi } from '../services/api';
import { message } from 'antd';

interface CategoryState {
  categories: Category[];
  loading: boolean;
  error: string | null;
}

type CategoryAction =
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_CATEGORIES'; payload: Category[] }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'ADD_CATEGORY'; payload: Category }
  | { type: 'UPDATE_CATEGORY'; payload: Category }
  | { type: 'DELETE_CATEGORY'; payload: number };

const initialState: CategoryState = {
  categories: [],
  loading: false,
  error: null,
};

const categoryReducer = (state: CategoryState, action: CategoryAction): CategoryState => {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, loading: action.payload };
    case 'SET_CATEGORIES':
      return {
        ...state,
        categories: action.payload,
        loading: false,
        error: null,
      };
    case 'SET_ERROR':
      return { ...state, error: action.payload, loading: false };
    case 'ADD_CATEGORY':
      return { ...state, categories: [...state.categories, action.payload] };
    case 'UPDATE_CATEGORY':
      return {
        ...state,
        categories: state.categories.map(category =>
          category.id === action.payload.id ? action.payload : category
        ),
      };
    case 'DELETE_CATEGORY':
      return {
        ...state,
        categories: state.categories.filter(category => category.id !== action.payload),
      };
    default:
      return state;
  }
};

interface CategoryContextType {
  state: CategoryState;
  fetchCategories: () => Promise<void>;
  createCategory: (categoryData: CategoryFormData) => Promise<boolean>;
  updateCategory: (id: number, categoryData: Partial<CategoryFormData>) => Promise<boolean>;
  deleteCategory: (id: number) => Promise<boolean>;
}

const CategoryContext = createContext<CategoryContextType | undefined>(undefined);

export const useCategories = () => {
  const context = useContext(CategoryContext);
  if (!context) {
    throw new Error('useCategories must be used within a CategoryProvider');
  }
  return context;
};

interface CategoryProviderProps {
  children: ReactNode;
}

export const CategoryProvider: React.FC<CategoryProviderProps> = ({ children }) => {
  const [state, dispatch] = useReducer(categoryReducer, initialState);

  const fetchCategories = useCallback(async () => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const response = await categoryApi.getAllCategories();
      dispatch({ type: 'SET_CATEGORIES', payload: response.data });
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch categories';
      dispatch({ type: 'SET_ERROR', payload: errorMessage });
      message.error(errorMessage);
    }
  }, []);

  const createCategory = useCallback(async (categoryData: CategoryFormData): Promise<boolean> => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const response = await categoryApi.createCategory(categoryData);
      dispatch({ type: 'ADD_CATEGORY', payload: response.data });
      message.success('Category created successfully!');
      return true;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create category';
      dispatch({ type: 'SET_ERROR', payload: errorMessage });
      message.error(errorMessage);
      return false;
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  }, []);

  const updateCategory = useCallback(async (id: number, categoryData: Partial<CategoryFormData>): Promise<boolean> => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const response = await categoryApi.updateCategory(id, categoryData);
      dispatch({ type: 'UPDATE_CATEGORY', payload: response.data });
      message.success('Category updated successfully!');
      return true;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update category';
      dispatch({ type: 'SET_ERROR', payload: errorMessage });
      message.error(errorMessage);
      return false;
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  }, []);

  const deleteCategory = useCallback(async (id: number): Promise<boolean> => {
    try {
      await categoryApi.deleteCategory(id);
      dispatch({ type: 'DELETE_CATEGORY', payload: id });
      message.success('Category deleted successfully!');
      return true;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to delete category';
      message.error(errorMessage);
      return false;
    }
  }, []);

  const value: CategoryContextType = {
    state,
    fetchCategories,
    createCategory,
    updateCategory,
    deleteCategory,
  };

  return (
    <CategoryContext.Provider value={value}>
      {children}
    </CategoryContext.Provider>
  );
};