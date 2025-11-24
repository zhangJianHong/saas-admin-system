import { apiRequest } from './api';
import { Overview, DatabaseStatus, Organization, OrganizationMetrics } from '../types';

// 仪表板API
export const dashboardApi = {
  // 获取概览数据
  getOverview: async (): Promise<Overview> => {
    return apiRequest.get('/dashboard/overview');
  },

  // 获取组织列表
  getOrganizations: async (): Promise<Organization[]> => {
    return apiRequest.get('/dashboard/organizations');
  },

  // 获取组织指标详情
  getOrganizationMetrics: async (organizationId: string): Promise<OrganizationMetrics> => {
    return apiRequest.get(`/dashboard/organizations/${organizationId}/metrics`);
  },

  // 获取数据库状态
  getDatabaseStatus: async (): Promise<DatabaseStatus> => {
    return apiRequest.get('/dashboard/database-status');
  },

  // 获取实时指标
  getRealtimeMetrics: async (params?: {
    organization_id?: string;
    metric_type?: string;
    hours?: number;
  }) => {
    return apiRequest.get('/monitoring/metrics', params);
  },

  // 获取指标历史数据
  getMetricsHistory: async (params: {
    organization_id?: string;
    metric_type: string;
    start_time: number;
    end_time: number;
  }) => {
    return apiRequest.get('/monitoring/metrics/history', params);
  },

  // 获取组织使用情况
  getOrganizationUsage: async (organizationId: string, hours: number = 24) => {
    return apiRequest.get(`/monitoring/organizations/${organizationId}/usage`, {
      hours,
    });
  },

  // 获取组织概览数据（新增的API）
  getOrganizationOverview: async () => {
    return apiRequest.get('/monitoring/organizations/overview');
  },
};