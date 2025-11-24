import { apiRequest, paginatedRequest } from './api';
import {
  SubscriptionPlan,
  CreateSubscriptionPlanRequest,
  UpdateSubscriptionPlanRequest,
  PaginatedResponse
} from '../types';

export class SubscriptionPlanService {
  /**
   * 获取订阅计划列表
   */
  static async getSubscriptionPlans(params: {
    page?: number;
    page_size?: number;
    search?: string;
  } = {}): Promise<PaginatedResponse<SubscriptionPlan>> {
    const response = await paginatedRequest<SubscriptionPlan>(
      '/subscription-plans',
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
   * 获取订阅计划详情
   */
  static async getSubscriptionPlanById(id: string): Promise<SubscriptionPlan> {
    return apiRequest.get<SubscriptionPlan>(`/subscription-plans/${id}`);
  }

  /**
   * 创建订阅计划
   */
  static async createSubscriptionPlan(data: CreateSubscriptionPlanRequest): Promise<SubscriptionPlan> {
    return apiRequest.post<SubscriptionPlan>('/subscription-plans', data);
  }

  /**
   * 更新订阅计划
   */
  static async updateSubscriptionPlan(
    id: string,
    data: UpdateSubscriptionPlanRequest
  ): Promise<SubscriptionPlan> {
    return apiRequest.put<SubscriptionPlan>(`/subscription-plans/${id}`, data);
  }

  /**
   * 删除订阅计划
   */
  static async deleteSubscriptionPlan(id: string): Promise<void> {
    return apiRequest.delete(`/subscription-plans/${id}`);
  }

  /**
   * 获取活跃的订阅计划
   */
  static async getActiveSubscriptionPlans(): Promise<SubscriptionPlan[]> {
    const response = await apiRequest.get<{ plans: SubscriptionPlan[] }>('/subscription-plans/active');
    return response.plans;
  }

  /**
   * 按价格范围获取订阅计划
   */
  static async getSubscriptionPlansByPricingRange(params: {
    min_price?: number;
    max_price?: number;
  }): Promise<SubscriptionPlan[]> {
    const response = await apiRequest.get<{ plans: SubscriptionPlan[] }>(
      '/subscription-plans/search',
      params
    );
    return response.plans;
  }
}

export default SubscriptionPlanService;