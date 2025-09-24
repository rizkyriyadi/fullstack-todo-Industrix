import React, { useState, useEffect } from 'react';
import {
  Layout,
  Menu,
  Button,
  Space,
  Input,
  Select,
  Drawer,
  Typography,
  Badge,
  theme,
} from 'antd';
import {
  PlusOutlined,
  SearchOutlined,
  FilterOutlined,
  MenuOutlined,
  TagOutlined,
} from '@ant-design/icons';
import { TodoForm } from '../todos/TodoForm';
import { CategoryForm } from '../categories/CategoryForm';
import { useTodos } from '../../contexts/TodoContext';
import { useCategories } from '../../contexts/CategoryContext';
import { TodoFilters, Priority } from '../../types';
import { getPriorityColor, getPriorityLabel, debounce } from '../../utils/helpers';

const { Header, Content, Sider } = Layout;
const { Title } = Typography;
const { Option } = Select;

interface AppLayoutProps {
  children: React.ReactNode;
}

export const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  const { token } = theme.useToken();
  const { state: todoState, setFilters, clearFilters, createTodo, fetchTodos } = useTodos();
  const { state: categoryState, fetchCategories, createCategory } = useCategories();

  const [collapsed, setCollapsed] = useState(false);
  const [mobileMenuVisible, setMobileMenuVisible] = useState(false);
  const [todoFormVisible, setTodoFormVisible] = useState(false);
  const [categoryFormVisible, setCategoryFormVisible] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

  // Debounced search
  const debouncedSetFilters = debounce(setFilters, 300);

  useEffect(() => {
    const handleResize = () => {
      const mobile = window.innerWidth < 768;
      setIsMobile(mobile);
      if (mobile) {
        setCollapsed(true);
      }
    };

    handleResize();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  useEffect(() => {
    fetchTodos();
    fetchCategories();
  }, [fetchTodos, fetchCategories]);

  useEffect(() => {
    fetchTodos();
  }, [todoState.filters, fetchTodos]);

  const handleSearch = (value: string) => {
    debouncedSetFilters({ search: value || undefined, page: 1 });
  };

  const handleFilterChange = (key: keyof TodoFilters, value: any) => {
    setFilters({ [key]: value, page: 1 });
  };

  const handleClearFilters = () => {
    clearFilters();
  };

  const completedCount = todoState.todos.filter(todo => todo.completed).length;
  const pendingCount = todoState.todos.filter(todo => !todo.completed).length;

  const siderContent = (
    <div style={{ 
      height: '100%', 
      display: 'flex', 
      flexDirection: 'column',
      overflow: 'hidden'
    }}>
      {/* Header */}
      <div style={{ 
        padding: collapsed ? '16px 8px' : '16px',
        borderBottom: `1px solid ${token.colorBorderSecondary}`,
        flexShrink: 0
      }}>
        {!collapsed && (
          <Title level={4} style={{ margin: 0, color: token.colorPrimary }}>
            Todo App
          </Title>
        )}
      </div>

      {/* Action Buttons */}
      <div style={{ 
        padding: collapsed ? '16px 8px' : '16px',
        flexShrink: 0
      }}>
        <Space direction="vertical" style={{ width: '100%' }} size="small">
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => setTodoFormVisible(true)}
            block={!collapsed}
            style={collapsed ? { width: '100%', padding: '8px' } : {}}
          >
            {!collapsed && 'New Todo'}
          </Button>
          <Button
            icon={<TagOutlined />}
            onClick={() => setCategoryFormVisible(true)}
            block={!collapsed}
            style={collapsed ? { width: '100%', padding: '8px' } : {}}
          >
            {!collapsed && 'New Category'}
          </Button>
        </Space>
      </div>

      {/* Menu */}
      <div style={{ flex: 1, overflow: 'auto' }}>
        <Menu
          mode="inline"
          selectedKeys={[]}
          style={{ border: 'none' }}
          inlineCollapsed={collapsed}
          items={[
            {
              key: 'all',
              icon: <Badge count={todoState.todos.length} size="small" />,
              label: 'All Todos',
              onClick: () => handleFilterChange('completed', undefined),
            },
            {
              key: 'pending',
              icon: <Badge count={pendingCount} size="small" color={token.colorWarning} />,
              label: 'Pending',
              onClick: () => handleFilterChange('completed', false),
            },
            {
              key: 'completed',
              icon: <Badge count={completedCount} size="small" color={token.colorSuccess} />,
              label: 'Completed',
              onClick: () => handleFilterChange('completed', true),
            },
          ]}
        />

        {/* Categories */}
        {categoryState.categories.length > 0 && !collapsed && (
          <div style={{ 
            padding: '16px',
            borderTop: `1px solid ${token.colorBorderSecondary}`,
            marginTop: 'auto'
          }}>
            <Title level={5} style={{ margin: '0 0 12px 0', fontSize: '14px' }}>
              Categories
            </Title>
            <div style={{ maxHeight: '200px', overflow: 'auto' }}>
              {categoryState.categories.map(category => (
                <div
                  key={category.id}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 8,
                    padding: '6px 8px',
                    cursor: 'pointer',
                    borderRadius: 4,
                    marginBottom: 4,
                    transition: 'background-color 0.2s',
                  }}
                  className="category-item"
                  onClick={() => handleFilterChange('category_id', category.id)}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = token.colorFillTertiary;
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = 'transparent';
                  }}
                >
                  <div
                    style={{
                      width: 8,
                      height: 8,
                      borderRadius: '50%',
                      backgroundColor: category.color,
                      flexShrink: 0,
                    }}
                  />
                  <span style={{ 
                    fontSize: '13px',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap'
                  }}>
                    {category.name}
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );

  return (
    <Layout style={{ minHeight: '100vh' }}>
      {!isMobile ? (
        <Sider
          collapsible
          collapsed={collapsed}
          onCollapse={setCollapsed}
          width={280}
          theme="light"
          style={{
            boxShadow: '2px 0 8px rgba(0,0,0,0.1)',
          }}
        >
          {siderContent}
        </Sider>
      ) : (
        <Drawer
          title={
            <Title level={4} style={{ margin: 0, color: token.colorPrimary }}>
              Todo App
            </Title>
          }
          placement="left"
          open={mobileMenuVisible}
          onClose={() => setMobileMenuVisible(false)}
          width={300}
          styles={{
            body: { padding: 0 }
          }}
        >
          <div style={{ 
            height: '100%', 
            display: 'flex', 
            flexDirection: 'column',
            overflow: 'hidden'
          }}>
            {/* Action Buttons */}
            <div style={{ padding: '16px', flexShrink: 0 }}>
              <Space direction="vertical" style={{ width: '100%' }} size="small">
                <Button
                  type="primary"
                  icon={<PlusOutlined />}
                  onClick={() => {
                    setTodoFormVisible(true);
                    setMobileMenuVisible(false);
                  }}
                  block
                  size="large"
                >
                  New Todo
                </Button>
                <Button
                  icon={<TagOutlined />}
                  onClick={() => {
                    setCategoryFormVisible(true);
                    setMobileMenuVisible(false);
                  }}
                  block
                  size="large"
                >
                  New Category
                </Button>
              </Space>
            </div>

            {/* Menu */}
            <div style={{ flex: 1, overflow: 'auto' }}>
              <Menu
                mode="inline"
                selectedKeys={[]}
                style={{ border: 'none' }}
                items={[
                  {
                    key: 'all',
                    icon: <Badge count={todoState.todos.length} size="small" />,
                    label: 'All Todos',
                    onClick: () => {
                      handleFilterChange('completed', undefined);
                      setMobileMenuVisible(false);
                    },
                  },
                  {
                    key: 'pending',
                    icon: <Badge count={pendingCount} size="small" color={token.colorWarning} />,
                    label: 'Pending',
                    onClick: () => {
                      handleFilterChange('completed', false);
                      setMobileMenuVisible(false);
                    },
                  },
                  {
                    key: 'completed',
                    icon: <Badge count={completedCount} size="small" color={token.colorSuccess} />,
                    label: 'Completed',
                    onClick: () => {
                      handleFilterChange('completed', true);
                      setMobileMenuVisible(false);
                    },
                  },
                ]}
              />

              {/* Categories */}
              {categoryState.categories.length > 0 && (
                <div style={{ 
                  padding: '16px',
                  borderTop: `1px solid ${token.colorBorderSecondary}`
                }}>
                  <Title level={5} style={{ margin: '0 0 12px 0', fontSize: '14px' }}>
                    Categories
                  </Title>
                  <div>
                    {categoryState.categories.map(category => (
                      <div
                        key={category.id}
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: 12,
                          padding: '8px 12px',
                          cursor: 'pointer',
                          borderRadius: 6,
                          marginBottom: 6,
                          transition: 'background-color 0.2s',
                        }}
                        onClick={() => {
                          handleFilterChange('category_id', category.id);
                          setMobileMenuVisible(false);
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.backgroundColor = token.colorFillTertiary;
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.backgroundColor = 'transparent';
                        }}
                      >
                        <div
                          style={{
                            width: 12,
                            height: 12,
                            borderRadius: '50%',
                            backgroundColor: category.color,
                            flexShrink: 0,
                          }}
                        />
                        <span style={{ 
                          fontSize: '14px',
                          overflow: 'hidden',
                          textOverflow: 'ellipsis',
                          whiteSpace: 'nowrap'
                        }}>
                          {category.name}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </Drawer>
      )}

      <Layout>
        <Header
          style={{
            background: token.colorBgContainer,
            padding: '0 16px',
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
            display: 'flex',
            alignItems: 'center',
            gap: 16,
          }}
        >
          {isMobile && (
            <Button
              type="text"
              icon={<MenuOutlined />}
              onClick={() => setMobileMenuVisible(true)}
            />
          )}

          <div style={{ flex: 1, display: 'flex', alignItems: 'center', gap: 16 }}>
            <Input
              placeholder="Search todos..."
              prefix={<SearchOutlined />}
              style={{ maxWidth: 300 }}
              onChange={(e) => handleSearch(e.target.value)}
              allowClear
            />

            <Select
              placeholder="Priority"
              style={{ width: 120 }}
              allowClear
              onChange={(value) => handleFilterChange('priority', value)}
            >
              {(['low', 'medium', 'high'] as Priority[]).map(priority => (
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

            <Button
              icon={<FilterOutlined />}
              onClick={handleClearFilters}
              type="text"
            >
              Clear Filters
            </Button>
          </div>
        </Header>

        <Content
          style={{
            padding: '24px',
            background: token.colorBgLayout,
            overflow: 'auto',
          }}
        >
          {children}
        </Content>
      </Layout>

      <TodoForm
        visible={todoFormVisible}
        onCancel={() => setTodoFormVisible(false)}
        onSubmit={createTodo}
        loading={todoState.loading}
      />

      <CategoryForm
        visible={categoryFormVisible}
        onCancel={() => setCategoryFormVisible(false)}
        onSubmit={createCategory}
        loading={categoryState.loading}
      />
    </Layout>
  );
};