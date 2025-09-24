import React, { useState } from 'react';
import { Empty, Spin, Modal } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { Todo, TodoFormData } from '../../types';
import { TodoItem } from './TodoItem';
import { TodoForm } from './TodoForm';
import { useTodos } from '../../contexts/TodoContext';

const { confirm } = Modal;

export const TodoList: React.FC = () => {
  const { state, toggleComplete, deleteTodo, updateTodo } = useTodos();
  const [editingTodo, setEditingTodo] = useState<Todo | null>(null);

  const handleToggleComplete = async (id: number) => {
    await toggleComplete(id);
  };

  const handleEdit = (todo: Todo) => {
    setEditingTodo(todo);
  };

  const handleDelete = (id: number) => {
    const todo = state.todos.find((t: Todo) => t.id === id);
    
    confirm({
      title: 'Delete Todo',
      icon: <ExclamationCircleOutlined />,
      content: `Are you sure you want to delete "${todo?.title}"?`,
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      onOk: async () => {
        await deleteTodo(id);
      },
    });
  };

  const handleUpdateTodo = async (formData: TodoFormData): Promise<boolean> => {
    if (!editingTodo) return false;
    const success = await updateTodo(editingTodo.id, formData);
    if (success) {
      setEditingTodo(null);
    }
    return success;
  };

  if (state.loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '40px 0' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (state.todos.length === 0) {
    return (
      <Empty
        description="No todos found"
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        style={{ padding: '40px 0' }}
      />
    );
  }

  return (
    <div>
      {state.todos.map(todo => (
        <TodoItem
          key={todo.id}
          todo={todo}
          onToggleComplete={handleToggleComplete}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      ))}

      {editingTodo && (
        <TodoForm
          visible={!!editingTodo}
          onCancel={() => setEditingTodo(null)}
          onSubmit={handleUpdateTodo}
          initialData={{
            title: editingTodo.title,
            description: editingTodo.description,
            priority: editingTodo.priority,
            category_id: editingTodo.category_id,
            due_date: editingTodo.due_date,
          }}
          loading={state.loading}
        />
      )}
    </div>
  );
};