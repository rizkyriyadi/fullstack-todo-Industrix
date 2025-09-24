import React, { useState, useEffect } from 'react';
import {
  Modal,
  Form,
  Input,
  Select,
  DatePicker,
  Button,
  Space,
  message,
} from 'antd';
import { TodoFormData, Category, Priority } from '../../types';
import { categoryApi } from '../../services/api';
import { getPriorityColor, getPriorityLabel } from '../../utils/helpers';
import dayjs from 'dayjs';

const { TextArea } = Input;
const { Option } = Select;

interface TodoFormProps {
  visible: boolean;
  onCancel: () => void;
  onSubmit: (data: TodoFormData) => Promise<boolean>;
  initialData?: Partial<TodoFormData>;
  loading?: boolean;
}

export const TodoForm: React.FC<TodoFormProps> = ({
  visible,
  onCancel,
  onSubmit,
  initialData,
  loading = false,
}) => {
  const [form] = Form.useForm();
  const [categories, setCategories] = useState<Category[]>([]);
  const [categoriesLoading, setCategoriesLoading] = useState(false);

  useEffect(() => {
    if (visible) {
      fetchCategories();
      if (initialData) {
        form.setFieldsValue({
          ...initialData,
          due_date: initialData.due_date ? dayjs(initialData.due_date) : null,
        });
      } else {
        form.resetFields();
      }
    }
  }, [visible, initialData, form]);

  const fetchCategories = async () => {
    setCategoriesLoading(true);
    try {
      const response = await categoryApi.getAllCategories();
      setCategories(response.data);
    } catch (error) {
      message.error('Failed to load categories');
    } finally {
      setCategoriesLoading(false);
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const formData: TodoFormData = {
        title: values.title,
        description: values.description || '',
        priority: values.priority || 'medium',
        category_id: values.category_id,
        due_date: values.due_date ? values.due_date.toISOString() : undefined,
      };

      const success = await onSubmit(formData);
      if (success) {
        form.resetFields();
        onCancel();
      }
    } catch (error) {
      // Form validation errors are handled by Ant Design
    }
  };

  const priorities: Priority[] = ['low', 'medium', 'high'];

  return (
    <Modal
      title={initialData ? 'Edit Todo' : 'Create New Todo'}
      open={visible}
      onCancel={onCancel}
      footer={[
        <Button key="cancel" onClick={onCancel}>
          Cancel
        </Button>,
        <Button
          key="submit"
          type="primary"
          loading={loading}
          onClick={handleSubmit}
        >
          {initialData ? 'Update' : 'Create'}
        </Button>,
      ]}
      destroyOnClose
      width={600}
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={{
          priority: 'medium',
        }}
      >
        <Form.Item
          name="title"
          label="Title"
          rules={[
            { required: true, message: 'Please enter a title' },
            { max: 255, message: 'Title must be less than 255 characters' },
          ]}
        >
          <Input placeholder="Enter todo title" />
        </Form.Item>

        <Form.Item
          name="description"
          label="Description"
        >
          <TextArea
            placeholder="Enter todo description"
            rows={3}
            maxLength={1000}
            showCount
          />
        </Form.Item>

        <Space style={{ width: '100%' }} size="large">
          <Form.Item
            name="priority"
            label="Priority"
            style={{ flex: 1 }}
          >
            <Select placeholder="Select priority">
              {priorities.map(priority => (
                <Option key={priority} value={priority}>
                  <Space>
                    <div
                      style={{
                        width: 8,
                        height: 8,
                        borderRadius: '50%',
                        backgroundColor: getPriorityColor(priority),
                      }}
                    />
                    {getPriorityLabel(priority)}
                  </Space>
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="category_id"
            label="Category"
            style={{ flex: 1 }}
          >
            <Select
              placeholder="Select category"
              loading={categoriesLoading}
              allowClear
            >
              {categories.map(category => (
                <Option key={category.id} value={category.id}>
                  <Space>
                    <div
                      style={{
                        width: 8,
                        height: 8,
                        borderRadius: '50%',
                        backgroundColor: category.color,
                      }}
                    />
                    {category.name}
                  </Space>
                </Option>
              ))}
            </Select>
          </Form.Item>
        </Space>

        <Form.Item
          name="due_date"
          label="Due Date"
        >
          <DatePicker
            style={{ width: '100%' }}
            showTime
            placeholder="Select due date"
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};