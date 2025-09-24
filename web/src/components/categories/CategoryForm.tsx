import React, { useEffect } from 'react';
import {
  Modal,
  Form,
  Input,
  Button,
  ColorPicker,
  Space,
} from 'antd';
import { CategoryFormData } from '../../types';

interface CategoryFormProps {
  visible: boolean;
  onCancel: () => void;
  onSubmit: (data: CategoryFormData) => Promise<boolean>;
  initialData?: Partial<CategoryFormData>;
  loading?: boolean;
}

const defaultColors = [
  '#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6',
  '#EC4899', '#14B8A6', '#F97316', '#84CC16', '#6366F1',
];

export const CategoryForm: React.FC<CategoryFormProps> = ({
  visible,
  onCancel,
  onSubmit,
  initialData,
  loading = false,
}) => {
  const [form] = Form.useForm();

  useEffect(() => {
    if (visible) {
      if (initialData) {
        form.setFieldsValue(initialData);
      } else {
        form.resetFields();
        form.setFieldsValue({
          color: '#3B82F6', // Default color
        });
      }
    }
  }, [visible, initialData, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const formData: CategoryFormData = {
        name: values.name,
        color: values.color,
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

  return (
    <Modal
      title={initialData ? 'Edit Category' : 'Create New Category'}
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
      width={500}
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={{
          color: '#3B82F6',
        }}
      >
        <Form.Item
          name="name"
          label="Category Name"
          rules={[
            { required: true, message: 'Please enter a category name' },
            { max: 100, message: 'Category name must be less than 100 characters' },
          ]}
        >
          <Input placeholder="Enter category name" />
        </Form.Item>

        <Form.Item
          name="color"
          label="Color"
          rules={[
            { required: true, message: 'Please select a color' },
          ]}
        >
          <ColorPicker
            showText
            presets={[
              {
                label: 'Recommended',
                colors: defaultColors,
              },
            ]}
            format="hex"
          />
        </Form.Item>

        <div style={{ marginTop: 16 }}>
          <div style={{ marginBottom: 8, fontSize: '14px', color: '#666' }}>
            Quick Colors:
          </div>
          <Space wrap>
            {defaultColors.map(color => (
              <div
                key={color}
                style={{
                  width: 24,
                  height: 24,
                  backgroundColor: color,
                  borderRadius: '50%',
                  cursor: 'pointer',
                  border: '2px solid #f0f0f0',
                }}
                onClick={() => form.setFieldsValue({ color })}
              />
            ))}
          </Space>
        </div>
      </Form>
    </Modal>
  );
};