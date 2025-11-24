import { apiRequest, paginatedRequest } from './api';
import { Organization, OrganizationSubscription, PaginatedResponse } from '../types';

export interface OrganizationSearchParams {
  page?: number;
  page_size?: number;
  search?: string;
}

export class OrganizationService {
  /**
   * 获取组织列表（只读）
   */
  static async getOrganizations(params: OrganizationSearchParams = {}): Promise<PaginatedResponse<Organization>> {
    const response = await paginatedRequest<Organization>(
      '/organizations',
      params
    );
    return {
      data: response.data,
      total: response.total,
      page: response.page,
      page_size: response.page_size,
      total_pages: response.total_pages,
    };
  }

  /**
   * 根据ID获取组织详情（只读）
   */
  static async getOrganizationById(id: string): Promise<Organization> {
    return apiRequest.get<Organization>(`/organizations/${id}`);
  }

  /**
   * 获取组织的用户列表（只读）
   */
  static async getOrganizationUsers(organizationId: string): Promise<any[]> {
    const response = await apiRequest.get<{ users: any[] }>(`/organizations/${organizationId}/users`);
    return response.users;
  }

  /**
   * 获取组织的工作空间列表（只读）
   */
  static async getOrganizationWorkspaces(organizationId: string): Promise<any[]> {
    const response = await apiRequest.get<{ workspaces: any[] }>(`/organizations/${organizationId}/workspaces`);
    return response.workspaces;
  }

  /**
   * 获取组织的订阅信息（只读）
   */
  static async getOrganizationSubscriptions(
    organizationId: string,
    params: { page: number; page_size: number }
  ): Promise<PaginatedResponse<OrganizationSubscription>> {
    return paginatedRequest<OrganizationSubscription>(
      `/organizations/${organizationId}/subscriptions`,
      params
    );
  }

  /**
   * 发送订阅到期提醒邮件（预留功能）
   */
  static async sendExpiryReminder(organizationId: string): Promise<{ message: string }> {
    return apiRequest.post<{ message: string }>(
      `/organizations/${organizationId}/send-expiry-reminder`,
      {}
    );
  }
}

export default OrganizationService;