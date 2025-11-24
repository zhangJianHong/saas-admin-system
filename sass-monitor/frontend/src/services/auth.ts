import { apiRequest, setAuthToken, clearAuthToken } from './api';
import { LoginRequest, LoginResponse, User } from '../types';

export class AuthService {
  /**
   * 用户登录
   */
  static async login(credentials: LoginRequest): Promise<LoginResponse> {
    try {
      const response = await apiRequest.post<LoginResponse>('/auth/login', credentials);

      // 保存认证信息
      if (response.token) {
        setAuthToken(response.token);
        localStorage.setItem('refresh_token', response.refresh_token);
        localStorage.setItem('user_info', JSON.stringify(response.user));
        localStorage.setItem('token_expires_at', response.expires_at);
      }

      return response;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }

  /**
   * 用户登出
   */
  static async logout(): Promise<void> {
    try {
      await apiRequest.post('/auth/logout');
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      // 无论接口调用是否成功，都清除本地认证信息
      clearAuthToken();
      localStorage.removeItem('user_info');
      localStorage.removeItem('token_expires_at');
    }
  }

  /**
   * 刷新token
   */
  static async refreshToken(): Promise<string> {
    try {
      const response = await apiRequest.post<{ token: string }>('/auth/refresh');

      if (response.token) {
        setAuthToken(response.token);
        // 更新过期时间（假设新的token有效期和原来一样）
        const oldExpiresAt = localStorage.getItem('token_expires_at');
        if (oldExpiresAt) {
          const expiresAt = new Date(oldExpiresAt);
          expiresAt.setHours(expiresAt.getHours() + 24); // 假设有效期24小时
          localStorage.setItem('token_expires_at', expiresAt.toISOString());
        }
      }

      return response.token;
    } catch (error) {
      console.error('Token refresh error:', error);
      clearAuthToken();
      localStorage.removeItem('user_info');
      localStorage.removeItem('token_expires_at');
      throw error;
    }
  }

  /**
   * 获取当前用户信息
   */
  static async getCurrentUser(): Promise<User> {
    try {
      const response = await apiRequest.get<User>('/profile');

      // 更新本地用户信息
      localStorage.setItem('user_info', JSON.stringify(response));

      return response;
    } catch (error) {
      console.error('Get current user error:', error);
      throw error;
    }
  }

  /**
   * 更新用户资料
   */
  static async updateProfile(userData: Partial<User>): Promise<User> {
    try {
      const response = await apiRequest.put<User>('/profile', userData);

      // 更新本地用户信息
      localStorage.setItem('user_info', JSON.stringify(response));

      return response;
    } catch (error) {
      console.error('Update profile error:', error);
      throw error;
    }
  }

  /**
   * 修改密码
   */
  static async changePassword(oldPassword: string, newPassword: string): Promise<void> {
    try {
      await apiRequest.post('/change-password', {
        old_password: oldPassword,
        new_password: newPassword,
      });
    } catch (error) {
      console.error('Change password error:', error);
      throw error;
    }
  }

  /**
   * 检查用户是否已登录
   */
  static isAuthenticated(): boolean {
    const token = localStorage.getItem('access_token');
    const expiresAt = localStorage.getItem('token_expires_at');

    if (!token) {
      return false;
    }

    if (expiresAt) {
      const expirationDate = new Date(expiresAt);
      if (expirationDate <= new Date()) {
        // Token已过期
        clearAuthToken();
        localStorage.removeItem('user_info');
        localStorage.removeItem('token_expires_at');
        return false;
      }
    }

    return true;
  }

  /**
   * 获取存储的用户信息
   */
  static getStoredUser(): User | null {
    try {
      const userInfo = localStorage.getItem('user_info');
      return userInfo ? JSON.parse(userInfo) : null;
    } catch (error) {
      console.error('Parse stored user info error:', error);
      return null;
    }
  }

  /**
   * 检查token是否即将过期（30分钟内）
   */
  static isTokenExpiringSoon(): boolean {
    const expiresAt = localStorage.getItem('token_expires_at');

    if (!expiresAt) {
      return false;
    }

    const expirationDate = new Date(expiresAt);
    const thirtyMinutesFromNow = new Date(Date.now() + 30 * 60 * 1000);

    return expirationDate <= thirtyMinutesFromNow;
  }

  /**
   * 检查用户权限
   */
  static hasRole(role: string): boolean {
    const user = this.getStoredUser();
    return user?.role === role;
  }

  /**
   * 检查用户是否为管理员
   */
  static isAdmin(): boolean {
    const user = this.getStoredUser();
    return user?.role === 'admin' || user?.role === 'super_admin';
  }

  /**
   * 检查用户是否为超级管理员
   */
  static isSuperAdmin(): boolean {
    const user = this.getStoredUser();
    return user?.role === 'super_admin';
  }

  /**
   * 获取认证头部
   */
  static getAuthHeader(): { Authorization: string } | {} {
    const token = localStorage.getItem('access_token');
    return token ? { Authorization: `Bearer ${token}` } : {};
  }

  /**
   * 初始化认证状态
   * 在应用启动时调用，检查用户登录状态
   */
  static async initializeAuth(): Promise<User | null> {
    if (!this.isAuthenticated()) {
      return null;
    }

    try {
      // 尝试获取当前用户信息验证token有效性
      const user = await this.getCurrentUser();
      return user;
    } catch (error) {
      // Token无效，清除认证信息
      clearAuthToken();
      localStorage.removeItem('user_info');
      localStorage.removeItem('token_expires_at');
      return null;
    }
  }
}

export default AuthService;