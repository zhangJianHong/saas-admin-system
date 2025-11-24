import { apiRequest, paginatedRequest } from './api';
import {
  LightAdminUser,
  UserOrganization,
  UserWorkspace,
  UserSubscription,
  PaginatedResponse
} from '../types';

export class UserService {
  /**
   * 获取用户列表
   */
  static async getUsers(params: {
    page: number;
    page_size: number;
    search?: string;
  }): Promise<PaginatedResponse<LightAdminUser>> {
    return paginatedRequest<LightAdminUser>('/users', params);
  }

  /**
   * 获取用户详情
   */
  static async getUserById(userId: string): Promise<LightAdminUser> {
    return apiRequest.get<LightAdminUser>(`/users/${userId}`);
  }

  /**
   * 获取用户所属组织列表
   */
  static async getUserOrganizations(
    userId: string,
    params: { page: number; page_size: number }
  ): Promise<PaginatedResponse<UserOrganization>> {
    return paginatedRequest<UserOrganization>(
      `/users/${userId}/organizations`,
      params
    );
  }

  /**
   * 获取用户所属工作空间列表
   */
  static async getUserWorkspaces(
    userId: string,
    params: { page: number; page_size: number }
  ): Promise<PaginatedResponse<UserWorkspace>> {
    return paginatedRequest<UserWorkspace>(
      `/users/${userId}/workspaces`,
      params
    );
  }

  /**
   * 获取用户订阅列表
   */
  static async getUserSubscriptions(
    userId: string,
    params: { page: number; page_size: number }
  ): Promise<PaginatedResponse<UserSubscription>> {
    return paginatedRequest<UserSubscription>(
      `/users/${userId}/subscriptions`,
      params
    );
  }
}

export default UserService;
