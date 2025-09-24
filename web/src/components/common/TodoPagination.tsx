import React from 'react';
import { Pagination, Typography, Space, Select } from 'antd';
import { useTodos } from '../../contexts/TodoContext';

const { Text } = Typography;
const { Option } = Select;

export const TodoPagination: React.FC = () => {
  const { state, setFilters } = useTodos();

  if (!state.pagination || state.pagination.total === 0) {
    return null;
  }

  const { pagination, filters } = state;

  const handlePageChange = (page: number, pageSize?: number) => {
    setFilters({
      page,
      per_page: pageSize || filters.per_page,
    });
  };

  const handlePageSizeChange = (pageSize: number) => {
    setFilters({
      page: 1,
      per_page: pageSize,
    });
  };

  const startItem = (pagination.current_page - 1) * pagination.per_page + 1;
  const endItem = Math.min(pagination.current_page * pagination.per_page, pagination.total);

  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginTop: 24,
        flexWrap: 'wrap',
        gap: 16,
      }}
    >
      <Space>
        <Text type="secondary">
          Showing {startItem} to {endItem} of {pagination.total} todos
        </Text>
      </Space>

      <Space>
        <Text type="secondary" style={{ fontSize: '13px' }}>
          Items per page:
        </Text>
        <Select
          size="small"
          value={pagination.per_page}
          onChange={handlePageSizeChange}
          style={{ width: 70 }}
        >
          <Option value={5}>5</Option>
          <Option value={10}>10</Option>
          <Option value={20}>20</Option>
          <Option value={50}>50</Option>
        </Select>

        <Pagination
          current={pagination.current_page}
          total={pagination.total}
          pageSize={pagination.per_page}
          onChange={handlePageChange}
          showSizeChanger={false}
          showQuickJumper
          size="small"
          responsive
        />
      </Space>
    </div>
  );
};