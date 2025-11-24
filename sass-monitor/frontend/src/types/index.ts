// 系统管理员用户类型
export interface User {
  id: string;
  username: string;
  email: string;
  full_name: string;
  role: string;
  last_login_at?: string;
}

// Light Admin用户类型
export interface LightAdminUser {
  id: string;
  username: string;
  nickname?: string;
  email?: string;
  avatar_url: string;
  email_verified?: boolean;
  oauth_provider: string;
  clerk_user_id: string;
  organization_count: number;
  workspace_count: number;
  subscription_count: number;
  created_at: string;
  updated_at?: string;
}

// 用户所属组织
export interface UserOrganization {
  organization_id: string;
  organization_name: string;
  description?: string;
  joined_at: string;
  created_at: string;
}

// 用户所属工作空间
export interface UserWorkspace {
  workspace_id: string;
  workspace_name: string;
  organization_id: string;
  organization_name: string;
  user_status: string;
  created_at: string;
}

// 用户订阅
export interface UserSubscription {
  subscription_id: string;
  plan_id: string;
  plan_name: string;
  organization_id: string;
  organization_name: string;
  status: string;
  billing_cycle: string;
  start_date: string;
  end_date?: string;
  trial_days_used?: number;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  refresh_token: string;
  expires_at: string;
  user: User;
}

// 监控相关类型
export interface DatabaseStatus {
  postgresql: {
    status: string;
    connections: number;
    max_connections: number;
    database_size: number;
    response_time: number;
  };
  clickhouse: Record<string, {
    status: string;
    database_size: number;
    table_count: number;
    row_count: number;
    response_time: number;
  }>;
  redis: {
    status: string;
    used_memory: number;
    max_memory: number;
    connected_clients: number;
    response_time: number;
  };
}

export interface Overview {
  total_organizations: number;
  total_users: number;
  total_subscriptions: number;
  database_status: Record<string, string>;
  system_health: Record<string, string>;
  recent_metrics: RecentMetric[];
}

export interface RecentMetric {
  database_name: string;
  metric_type: string;
  metric_value: number;
  unit: string;
  timestamp: number;
}

// 组织相关类型
export interface Organization {
  id: string;
  name: string;
  owner_id: string;
  description?: string;
  user_count: number;
  subscription_count: number;
  active_subscription_count: number;
  workspace_count: number;
  storage_usage: number;
  // 订阅到期相关
  subscription_status: 'active' | 'expiring_soon' | 'expired' | 'none';
  subscription_end_date?: string;
  days_until_expiration?: number;
  created_at: string;
  updated_at?: string;
}

// 组织订阅详情
export interface OrganizationSubscription {
  id: string;
  user_id: string;
  username: string;
  user_email?: string;
  plan_id: string;
  plan_name: string;
  plan_pricing: number;
  status: string;
  billing_cycle: string;
  start_date: string;
  end_date?: string;
  days_until_expiry?: number;
  payment_method?: string;
  last_billed_at?: string;
  trial_days_used?: number;
  created_at: string;
}

export interface OrganizationMetrics {
  organization: Organization;
  user_stats: {
    total_users: number;
    active_users: number;
  };
  sub_stats: {
    total_subscriptions: number;
    active_subscriptions: number;
    monthly_revenue: number;
  };
  workspace_count: number;
  resource_usage: {
    storage_usage_mb: number;
    query_count_today: number;
  };
}

// 告警相关类型
export interface AlertRule {
  id: string;
  name: string;
  description?: string;
  rule_type: string;
  target_type: string;
  target_name?: string;
  metric_name: string;
  operator: string;
  threshold: number;
  duration: number;
  severity: string;
  enabled: boolean;
  notification_config: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface CreateAlertRequest {
  name: string;
  description?: string;
  rule_type: string;
  target_type: string;
  target_name?: string;
  metric_name: string;
  operator: string;
  threshold: number;
  duration?: number;
  severity?: string;
  notification_config?: string;
}

// 监控指标相关类型
export interface ResourceMetric {
  id: string;
  organization_id?: string;
  database_type: string;
  database_name: string;
  metric_type: string;
  metric_name: string;
  metric_value: number;
  unit?: string;
  tags?: string;
  collected_at: string;
  created_at: string;
}

export interface MetricsHistory {
  database_type: string;
  database_name: string;
  metric_type: string;
  start_time: number;
  end_time: number;
  data: Array<{
    timestamp: number;
    value: number;
    metric_name: string;
    unit: string;
  }>;
}

// 系统相关类型
export interface SystemHealth {
  id: string;
  component_name: string;
  component_type: string;
  status: string;
  response_time?: number;
  error_message?: string;
  last_checked_at: string;
  created_at: string;
  updated_at: string;
}

export interface MonitoringLog {
  id: string;
  log_level: string;
  source: string;
  component?: string;
  organization_id?: string;
  message: string;
  details?: string;
  created_at: string;
}

// API响应类型
export interface ApiResponse<T = any> {
  data?: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T = any> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

// 主题相关类型
export type ThemeMode = 'light' | 'dark';

export interface ThemeConfig {
  mode: ThemeMode;
  primary_color: string;
  border_radius: number;
}

// 路由相关类型
export interface RouteConfig {
  path: string;
  component: React.ComponentType;
  title: string;
  icon?: React.ReactNode;
  children?: RouteConfig[];
  required?: string[]; // 所需权限
}

// 表格相关类型
export interface TableColumn {
  title: string;
  dataIndex: string;
  key: string;
  width?: number;
  align?: 'left' | 'center' | 'right';
  sorter?: boolean;
  fixed?: 'left' | 'right';
  render?: (value: any, record: any, index: number) => React.ReactNode;
}

// 图表相关类型
export interface ChartData {
  name: string;
  value: number;
  timestamp?: number;
  [key: string]: any;
}

export interface ChartConfig {
  title: string;
  type: 'line' | 'bar' | 'pie' | 'area' | 'gauge';
  data: ChartData[];
  xAxis?: string;
  yAxis?: string;
  height?: number;
  color?: string[];
}

// 通知相关类型
export interface NotificationItem {
  id: string;
  type: 'success' | 'info' | 'warning' | 'error';
  title: string;
  message: string;
  duration?: number;
  timestamp?: number;
}

// 错误相关类型
export interface ApiError {
  code: string;
  message: string;
  details?: any;
}

// 加载状态类型
export interface LoadingState {
  [key: string]: boolean;
}

// 表单相关类型
export interface FormField {
  name: string;
  label: string;
  type: 'input' | 'password' | 'select' | 'textarea' | 'number' | 'date' | 'switch';
  required?: boolean;
  placeholder?: string;
  options?: Array<{ label: string; value: any }>;
  rules?: any[];
}

// 搜索相关类型
export interface SearchParams {
  keyword?: string;
  page?: number;
  page_size?: number;
  sort_field?: string;
  sort_order?: 'asc' | 'desc';
  [key: string]: any;
}

// 订阅计划相关类型
export interface SubscriptionPlan {
  id: string;
  tier_name: string;
  pricing_monthly: number;
  pricing_quarterly: number;
  pricing_yearly: number;
  limits: string; // JSON string
  features?: string; // JSON string
  target_users?: string;
  upgrade_path?: string;
  is_custom?: boolean;
  default_flow_package?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  stripe_price_id_monthly?: string;
  stripe_price_id_quarterly?: string;
  stripe_price_id_yearly?: string;
}

export interface CreateSubscriptionPlanRequest {
  tier_name: string;
  pricing_monthly: number;
  pricing_quarterly: number;
  pricing_yearly: number;
  limits: string; // JSON string
  features?: string; // JSON string
  target_users?: string;
  upgrade_path?: string;
  is_custom?: boolean;
  default_flow_package?: string;
  is_active: boolean;
  stripe_price_id_monthly?: string;
  stripe_price_id_quarterly?: string;
  stripe_price_id_yearly?: string;
}

export interface UpdateSubscriptionPlanRequest {
  tier_name?: string;
  pricing_monthly?: number;
  pricing_quarterly?: number;
  pricing_yearly?: number;
  limits?: string; // JSON string
  features?: string; // JSON string
  target_users?: string;
  upgrade_path?: string;
  is_custom?: boolean;
  default_flow_package?: string;
  is_active?: boolean;
  stripe_price_id_monthly?: string;
  stripe_price_id_quarterly?: string;
  stripe_price_id_yearly?: string;
}