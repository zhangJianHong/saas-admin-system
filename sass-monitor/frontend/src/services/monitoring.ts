import { apiRequest } from './api';
import {
  Overview,
  DatabaseStatus,
  Organization,
  OrganizationMetrics,
  AlertRule,
  CreateAlertRequest,
  ResourceMetric,
  MetricsHistory,
  SystemHealth,
  MonitoringLog,
  SearchParams
} from '../types';

export class MonitoringService {
  /**
   * 获取仪表板概览
   */
  static async getOverview(): Promise<Overview> {
    return apiRequest.get<Overview>('/dashboard/overview');
  }

  /**
   * 获取组织列表
   */
  static async getOrganizations(params: SearchParams = {}): Promise<{
    organizations: Organization[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
  }> {
    const response = await apiRequest.get<{
      organizations: Organization[];
      total: number;
      page: number;
      page_size: number;
      total_pages: number;
    }>('/dashboard/organizations', { params });
    return response;
  }

  /**
   * 获取组织指标详情
   */
  static async getOrganizationMetrics(orgId: string): Promise<OrganizationMetrics> {
    return apiRequest.get<OrganizationMetrics>(`/dashboard/organizations/${orgId}/metrics`);
  }

  /**
   * 获取数据库状态
   */
  static async getDatabaseStatus(): Promise<DatabaseStatus> {
    return apiRequest.get<DatabaseStatus>('/dashboard/database-status');
  }

  /**
   * 获取监控指标
   */
  static async getMetrics(params: SearchParams = {}): Promise<{
    metrics: ResourceMetric[];
    total: number;
    page: number;
    page_size: number;
  }> {
    const response = await apiRequest.get<{
      metrics: ResourceMetric[];
      total: number;
      page: number;
      page_size: number;
    }>('/monitoring/metrics', { params });
    return response;
  }

  /**
   * 获取指标历史数据
   */
  static async getMetricsHistory(params: {
    database_type?: string;
    database_name?: string;
    metric_type?: string;
    organization_id?: string;
    hours?: number;
    start_time?: string;
    end_time?: string;
  } = {}): Promise<MetricsHistory> {
    return apiRequest.get<MetricsHistory>('/monitoring/metrics/history', params);
  }

  /**
   * 获取所有组织（用于筛选）
   */
  static async getOrganizationsForFilter(search?: string): Promise<Organization[]> {
    return apiRequest.get<Organization[]>('/monitoring/organizations', { search });
  }

  /**
   * 获取组织资源使用情况
   */
  static async getOrganizationUsage(orgId: string, hours: number = 24): Promise<{
    organization_id: string;
    period: {
      start_time: number;
      end_time: number;
      hours: number;
    };
    usage_by_database: Record<string, Record<string, Array<{
      timestamp: number;
      metric_name: string;
      metric_value: number;
      unit: string;
    }>>>;
    usage_by_type: Record<string, Array<{
      timestamp: number;
      database_name: string;
      metric_name: string;
      metric_value: number;
      unit: string;
    }>>;
  }> {
    return apiRequest.get(`/monitoring/organizations/${orgId}/usage`, { hours });
  }

  /**
   * 获取数据库详细信息
   */
  static async getDatabaseInfo(type?: string): Promise<any> {
    return apiRequest.get('/monitoring/databases', { type });
  }

  /**
   * 获取告警列表
   */
  static async getAlerts(params: SearchParams = {}): Promise<{
    alerts: AlertRule[];
    total: number;
    page: number;
    page_size: number;
  }> {
    const response = await apiRequest.get<{
      alerts: AlertRule[];
      total: number;
      page: number;
      page_size: number;
    }>('/monitoring/alerts', { params });
    return response;
  }

  /**
   * 创建告警规则
   */
  static async createAlert(alertData: CreateAlertRequest): Promise<AlertRule> {
    return apiRequest.post<AlertRule>('/monitoring/alerts', alertData);
  }

  /**
   * 更新告警规则
   */
  static async updateAlert(alertId: string, alertData: Partial<CreateAlertRequest>): Promise<AlertRule> {
    return apiRequest.put<AlertRule>(`/monitoring/alerts/${alertId}`, alertData);
  }

  /**
   * 删除告警规则
   */
  static async deleteAlert(alertId: string): Promise<void> {
    return apiRequest.delete(`/monitoring/alerts/${alertId}`);
  }

  /**
   * 获取系统健康状态
   */
  static async getSystemHealth(): Promise<{
    health_records: SystemHealth[];
  }> {
    return apiRequest.get('/system/health');
  }

  /**
   * 获取系统日志
   */
  static async getSystemLogs(params: SearchParams = {}): Promise<{
    logs: MonitoringLog[];
    total: number;
    page: number;
    page_size: number;
  }> {
    const response = await apiRequest.get<{
      logs: MonitoringLog[];
      total: number;
      page: number;
      page_size: number;
    }>('/system/logs', { params });
    return response;
  }

  /**
   * 获取系统配置
   */
  static async getSystemConfigs(): Promise<{
    configs: Record<string, string>;
  }> {
    return apiRequest.get('/system/configs');
  }

  /**
   * 更新系统配置
   */
  static async updateSystemConfigs(configs: Record<string, any>): Promise<void> {
    return apiRequest.put('/system/configs', configs);
  }

  /**
   * 获取实时指标数据
   */
  static async getRealTimeMetrics(databaseType: string, databaseName?: string): Promise<{
    timestamp: number;
    metrics: Array<{
      name: string;
      value: number;
      unit: string;
    }>;
  }> {
    return apiRequest.get('/monitoring/metrics/realtime', {
      database_type: databaseType,
      database_name: databaseName,
    });
  }

  /**
   * 获取性能指标对比
   */
  static async getPerformanceComparison(params: {
    database_type: string;
    start_date: string;
    end_date: string;
    compare_start_date?: string;
    compare_end_date?: string;
  }): Promise<{
    current_period: Array<{
      date: string;
      cpu_usage: number;
      memory_usage: number;
      disk_usage: number;
      response_time: number;
    }>;
    comparison_period?: Array<{
      date: string;
      cpu_usage: number;
      memory_usage: number;
      disk_usage: number;
      response_time: number;
    }>;
    summary: {
      avg_cpu: number;
      avg_memory: number;
      avg_disk: number;
      avg_response_time: number;
    };
  }> {
    return apiRequest.get('/monitoring/performance/comparison', params);
  }

  /**
   * 获取趋势分析数据
   */
  static async getTrendAnalysis(params: {
    metric_name: string;
    database_type: string;
    database_name?: string;
    period: '7d' | '30d' | '90d' | '1y';
  }): Promise<{
    trend_data: Array<{
      date: string;
      value: number;
    }>;
    trend_direction: 'up' | 'down' | 'stable';
    trend_percentage: number;
    forecast: Array<{
      date: string;
      predicted_value: number;
      confidence_interval: [number, number];
    }>;
  }> {
    return apiRequest.get('/monitoring/trends/analysis', params);
  }

  /**
   * 导出监控数据
   */
  static async exportMonitoringData(params: {
    format: 'csv' | 'excel' | 'json';
    start_date: string;
    end_date: string;
    database_type?: string;
    organization_id?: string;
  }): Promise<Blob> {
    const response = await fetch(`/api/v1/monitoring/export`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
      body: JSON.stringify(params),
    });

    if (!response.ok) {
      throw new Error('Export failed');
    }

    return response.blob();
  }

  /**
   * 获取告警统计
   */
  static async getAlertStatistics(params: {
    start_date?: string;
    end_date?: string;
    severity?: string;
  } = {}): Promise<{
    total_alerts: number;
    alerts_by_severity: Record<string, number>;
    alerts_by_type: Record<string, number>;
    alerts_trend: Array<{
      date: string;
      count: number;
    }>;
    top_alerts: Array<{
      rule_name: string;
      count: number;
      severity: string;
    }>;
  }> {
    return apiRequest.get('/monitoring/alerts/statistics', params);
  }
}

export default MonitoringService;