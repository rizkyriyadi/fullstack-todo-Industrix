import { Priority } from '../types';
import dayjs from 'dayjs';

export const priorityColors = {
  low: '#52c41a',
  medium: '#faad14',
  high: '#ff4d4f',
};

export const priorityLabels = {
  low: 'Low',
  medium: 'Medium',
  high: 'High',
};

export const getPriorityColor = (priority: Priority): string => {
  return priorityColors[priority] || priorityColors.medium;
};

export const getPriorityLabel = (priority: Priority): string => {
  return priorityLabels[priority] || priorityLabels.medium;
};

export const formatDate = (dateString: string): string => {
  return dayjs(dateString).format('MMM DD, YYYY');
};

export const formatDateTime = (dateString: string): string => {
  return dayjs(dateString).format('MMM DD, YYYY HH:mm');
};

export const isOverdue = (dueDateString?: string): boolean => {
  if (!dueDateString) return false;
  return dayjs(dueDateString).isBefore(dayjs(), 'day');
};

export const isDueSoon = (dueDateString?: string): boolean => {
  if (!dueDateString) return false;
  const dueDate = dayjs(dueDateString);
  const today = dayjs();
  return dueDate.isAfter(today) && dueDate.diff(today, 'day') <= 3;
};

export const getStatusColor = (completed: boolean, dueDate?: string): string => {
  if (completed) return '#52c41a';
  if (isOverdue(dueDate)) return '#ff4d4f';
  if (isDueSoon(dueDate)) return '#faad14';
  return '#1890ff';
};

export const truncateText = (text: string, maxLength: number = 50): string => {
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength) + '...';
};

export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  wait: number
): ((...args: Parameters<T>) => void) => {
  let timeout: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
};