// 渠道业务类型常量定义
export const BUSINESS_TYPE_OPTIONS = [
  { value: 1, color: 'blue', label: '对话', icon: '💬' },
  { value: 2, color: 'green', label: '应用', icon: '🔧' },
  { value: 3, color: 'purple', label: '工作流', icon: '⚡' },
];

// 业务类型映射
export const BUSINESS_TYPE_MAP = {
  1: { label: '对话', color: 'blue', icon: '💬' },
  2: { label: '应用', color: 'green', icon: '🔧' },
  3: { label: '工作流', color: 'purple', icon: '⚡' },
};

// 获取业务类型标签
export const getBusinessTypeLabel = (type) => {
  return BUSINESS_TYPE_MAP[type]?.label || '未知';
};

// 获取业务类型颜色
export const getBusinessTypeColor = (type) => {
  return BUSINESS_TYPE_MAP[type]?.color || 'grey';
};

// 获取业务类型图标
export const getBusinessTypeIcon = (type) => {
  return BUSINESS_TYPE_MAP[type]?.icon || '❓';
};
