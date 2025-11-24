import axios, { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import { message } from 'antd';

// API基础配置
const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://127.0.0.1:8080/api/v1';
const REQUEST_TIMEOUT = 30000; // 30秒

// 创建axios实例
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: REQUEST_TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 添加认证token
    const token = localStorage.getItem('access_token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // 添加请求时间戳
    config.metadata = { startTime: new Date() };

    return config;
  },
  (error) => {
    console.error('Request interceptor error:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    // 计算请求耗时
    const endTime = new Date();
    const startTime = response.config.metadata?.startTime?.getTime();
    const duration = startTime ? endTime.getTime() - startTime : 0;

    console.log(`API请求耗时: ${duration}ms - ${response.config.method?.toUpperCase()} ${response.config.url}`);

    return response;
  },
  (error) => {
    // 统一错误处理
    handleApiError(error);
    return Promise.reject(error);
  }
);

// 错误处理函数
const handleApiError = (error: any) => {
  const { response, request } = error;

  if (response) {
    // 服务器响应了错误状态码
    const { status, data } = response;

    switch (status) {
      case 401:
        // 未授权，清除token并跳转到登录页
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user_info');

        // 使用React Router导航而不是window.location
        if (window.location.pathname !== '/login') {
          window.location.href = '/login';
        }
        message.error('登录已过期，请重新登录');
        break;

      case 403:
        message.error('没有权限访问此资源');
        break;

      case 404:
        message.error('请求的资源不存在');
        break;

      case 429:
        message.error('请求过于频繁，请稍后再试');
        break;

      case 500:
        message.error('服务器内部错误，请稍后再试');
        break;

      default:
        message.error(data?.error || data?.message || '请求失败');
    }
  } else if (request) {
    // 请求已发出，但没有收到响应
    message.error('网络连接失败，请检查网络设置');
  } else {
    // 请求配置出错
    message.error('请求配置错误');
  }

  console.error('API Error:', error);
};

// 通用请求方法
export const apiRequest = {
  // GET请求
  get: async <T = any>(url: string, params?: any, config?: InternalAxiosRequestConfig): Promise<T> => {
    const response = await apiClient.get<T>(url, { ...config, params });
    return response.data;
  },

  // POST请求
  post: async <T = any>(url: string, data?: any, config?: InternalAxiosRequestConfig): Promise<T> => {
    const response = await apiClient.post<T>(url, data, config);
    return response.data;
  },

  // PUT请求
  put: async <T = any>(url: string, data?: any, config?: InternalAxiosRequestConfig): Promise<T> => {
    const response = await apiClient.put<T>(url, data, config);
    return response.data;
  },

  // DELETE请求
  delete: async <T = any>(url: string, config?: InternalAxiosRequestConfig): Promise<T> => {
    const response = await apiClient.delete<T>(url, config);
    return response.data;
  },

  // PATCH请求
  patch: async <T = any>(url: string, data?: any, config?: InternalAxiosRequestConfig): Promise<T> => {
    const response = await apiClient.patch<T>(url, data, config);
    return response.data;
  },
};

// 分页请求方法
export const paginatedRequest = async <T = any>(
  url: string,
  params: any = {},
  config?: InternalAxiosRequestConfig
): Promise<{
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}> => {
  const response = await apiClient.get(url, { ...config, params });
  return response.data;
};

// 文件上传方法
export const uploadFile = async (
  url: string,
  file: File,
  onProgress?: (progress: number) => void,
  config?: InternalAxiosRequestConfig
): Promise<any> => {
  const formData = new FormData();
  formData.append('file', file);

  return apiClient.post(url, formData, {
    ...config,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    onUploadProgress: (progressEvent) => {
      if (onProgress && progressEvent.total) {
        const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total);
        onProgress(progress);
      }
    },
  });
};

// 批量请求方法
export const batchRequest = async <T = any>(
  requests: Array<{ url: string; method?: string; data?: any }>
): Promise<T[]> => {
  const promises = requests.map(({ url, method = 'GET', data }) => {
    switch (method.toUpperCase()) {
      case 'POST':
        return apiRequest.post(url, data);
      case 'PUT':
        return apiRequest.put(url, data);
      case 'DELETE':
        return apiRequest.delete(url);
      default:
        return apiRequest.get(url);
    }
  });

  return Promise.all(promises);
};

// 重试请求方法
export const retryRequest = async <T = any>(
  requestFn: () => Promise<T>,
  maxRetries: number = 3,
  delay: number = 1000
): Promise<T> => {
  let lastError: any;

  for (let i = 0; i <= maxRetries; i++) {
    try {
      return await requestFn();
    } catch (error) {
      lastError = error;

      if (i < maxRetries) {
        console.warn(`Request failed, retrying... (${i + 1}/${maxRetries})`);
        await new Promise(resolve => setTimeout(resolve, delay * Math.pow(2, i)));
      }
    }
  }

  throw lastError;
};

// 取消请求方法
export const createCancelToken = () => {
  return axios.CancelToken.source();
};

// 请求状态检查
export const isRequestCanceled = (error: any): boolean => {
  return axios.isCancel(error);
};

// 设置认证token
export const setAuthToken = (token: string) => {
  localStorage.setItem('access_token', token);
  apiClient.defaults.headers.common['Authorization'] = `Bearer ${token}`;
};

// 清除认证token
export const clearAuthToken = () => {
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
  delete apiClient.defaults.headers.common['Authorization'];
};

// 刷新token
export const refreshAuthToken = async (): Promise<string> => {
  const refreshToken = localStorage.getItem('refresh_token');
  if (!refreshToken) {
    throw new Error('No refresh token available');
  }

  try {
    const response = await apiClient.post<{ token: string }>('/auth/refresh', {}, {
      headers: {
        'Refresh-Token': refreshToken,
      } as any,
    });

    const newToken = response.data.token;
    setAuthToken(newToken);

    return newToken;
  } catch (error) {
    clearAuthToken();
    throw error;
  }
};

// 添加请求元数据类型
declare module 'axios' {
  interface InternalAxiosRequestConfig {
    metadata?: {
      startTime?: Date;
    };
  }
}

export default apiClient;