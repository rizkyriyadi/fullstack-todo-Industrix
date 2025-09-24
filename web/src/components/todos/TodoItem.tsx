import React from 'react';
import {
  Card,
  Checkbox,
  Tag,
  Typography,
  Space,
  Button,
  Dropdown,
  Tooltip,
  theme,
} from 'antd';
import {
  EditOutlined,
  DeleteOutlined,
  MoreOutlined,
  CalendarOutlined,
  FlagOutlined,
} from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { Todo } from '../../types';
import {
  getPriorityColor,
  getPriorityLabel,
  formatDate,
  isOverdue,
  isDueSoon,
  getStatusColor,
} from '../../utils/helpers';

const { Text, Paragraph } = Typography;

interface TodoItemProps {
  todo: Todo;
  onToggleComplete: (id: number) => void;
  onEdit: (todo: Todo) => void;
  onDelete: (id: number) => void;
}

export const TodoItem: React.FC<TodoItemProps> = ({
  todo,
  onToggleComplete,
  onEdit,
  onDelete,
}) => {
  const { token } = theme.useToken();

  const menuItems: MenuProps['items'] = [
    {
      key: 'edit',
      icon: <EditOutlined />,
      label: 'Edit',
      onClick: () => onEdit(todo),
    },
    {
      key: 'delete',
      icon: <DeleteOutlined />,
      label: 'Delete',
      danger: true,
      onClick: () => onDelete(todo.id),
    },
  ];

  const getDueDateStatus = () => {
    if (!todo.due_date) return null;
    
    if (isOverdue(todo.due_date) && !todo.completed) {
      return { color: '#ff4d4f', text: 'Overdue' };
    }
    
    if (isDueSoon(todo.due_date) && !todo.completed) {
      return { color: '#faad14', text: 'Due Soon' };
    }
    
    return { color: '#1890ff', text: formatDate(todo.due_date) };
  };

  const dueDateStatus = getDueDateStatus();

  return (
    <Card
      size="small"
      hoverable
      style={{
        marginBottom: 12,
        opacity: todo.completed ? 0.7 : 1,
        borderLeft: `4px solid ${getStatusColor(todo.completed, todo.due_date)}`,
      }}
      bodyStyle={{ padding: '12px 16px' }}
    >
      <div style={{ display: 'flex', alignItems: 'flex-start', gap: 12 }}>
        <Checkbox
          checked={todo.completed}
          onChange={() => onToggleComplete(todo.id)}
          style={{ marginTop: 2 }}
        />
        
        <div style={{ flex: 1, minWidth: 0 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 4 }}>
            <Text
              strong
              style={{
                textDecoration: todo.completed ? 'line-through' : 'none',
                color: todo.completed ? token.colorTextDisabled : token.colorText,
                flex: 1,
              }}
            >
              {todo.title}
            </Text>
            
            <Space size={4}>
              <Tooltip title={`Priority: ${getPriorityLabel(todo.priority)}`}>
                <FlagOutlined
                  style={{ color: getPriorityColor(todo.priority) }}
                />
              </Tooltip>
              
              {todo.category && (
                <Tag
                  color={todo.category.color}
                  style={{ margin: 0, fontSize: '11px' }}
                >
                  {todo.category.name}
                </Tag>
              )}
            </Space>
          </div>

          {todo.description && (
            <Paragraph
              style={{
                margin: '4px 0 8px 0',
                color: todo.completed ? token.colorTextDisabled : token.colorTextSecondary,
                fontSize: '13px',
              }}
              ellipsis={{ rows: 2, expandable: true, symbol: 'more' }}
            >
              {todo.description}
            </Paragraph>
          )}

          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <Space size={12}>
              {dueDateStatus && (
                <Space size={4}>
                  <CalendarOutlined style={{ color: dueDateStatus.color, fontSize: '12px' }} />
                  <Text style={{ color: dueDateStatus.color, fontSize: '12px' }}>
                    {dueDateStatus.text}
                  </Text>
                </Space>
              )}
            </Space>

            <Text type="secondary" style={{ fontSize: '11px' }}>
              {formatDate(todo.created_at)}
            </Text>
          </div>
        </div>

        <Dropdown menu={{ items: menuItems }} trigger={['click']} placement="bottomRight">
          <Button
            type="text"
            size="small"
            icon={<MoreOutlined />}
            style={{ flexShrink: 0 }}
          />
        </Dropdown>
      </div>
    </Card>
  );
};