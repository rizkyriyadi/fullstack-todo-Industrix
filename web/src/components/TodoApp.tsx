import React from 'react';
import { ConfigProvider, theme } from 'antd';
import { TodoProvider } from '../contexts/TodoContext';
import { CategoryProvider } from '../contexts/CategoryContext';
import { AppLayout } from './layout/AppLayout';
import { TodoList } from './todos/TodoList';
import { TodoPagination } from './common/TodoPagination';

export const TodoApp: React.FC = () => {
  return (
    <ConfigProvider
      theme={{
        algorithm: theme.defaultAlgorithm,
        token: {
          colorPrimary: '#1890ff',
          borderRadius: 8,
        },
      }}
    >
      <CategoryProvider>
        <TodoProvider>
          <AppLayout>
            <div style={{ maxWidth: 1200, margin: '0 auto' }}>
              <TodoList />
              <TodoPagination />
            </div>
          </AppLayout>
        </TodoProvider>
      </CategoryProvider>
    </ConfigProvider>
  );
};