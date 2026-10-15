import { ApiResponse, OpptyJobDetailVO, OpptyHistoryItemVO, UserProfileVO } from '../types';

const BASE_PATH = ((import.meta as any).env?.BASE_URL as string) || '/';
const API_BASE = BASE_PATH.endsWith('/') ? `${BASE_PATH}api/v1` : `${BASE_PATH}/api/v1`;

/**
 * 通用网络请求封装
 */
async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem('oppty_token');
  const headers = new Headers(options.headers || {});

  if (!headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  const url = path.startsWith('http') ? path : `${API_BASE}${path}`;
  const response = await fetch(url, {
    ...options,
    headers,
  });

  const resJson: ApiResponse<T> = await response.json().catch(() => {
    throw new Error('网络请求解析异常，请稍后重试');
  });

  if (!response.ok || (resJson.code !== 200 && resJson.code !== 202)) {
    throw new Error(resJson.message || `请求失败 (HTTP ${response.status})`);
  }

  return resJson.data;
}

/**
 * 认证 API
 */
export const authApi = {
  // 发送邮箱验证码
  sendCode: (email: string) => {
    return request<{ cooldownSeconds: number }>('/auth/send-code', {
      method: 'POST',
      body: JSON.stringify({ email }),
    });
  },

  // 校验验证码并登录
  verifyCode: (email: string, code: string) => {
    return request<{ token: string; user: UserProfileVO }>('/auth/verify-code', {
      method: 'POST',
      body: JSON.stringify({ email, code }),
    });
  },

  // 退出登录
  logout: () => {
    return request<void>('/auth/logout', { method: 'POST' });
  },
};

/**
 * 商机挖掘与任务 API
 */
export const opptyApi = {
  // 提交商机分析任务
  createJob: (params: {
    keywords: string[];
    dataSource: string;
    deepCrawl: boolean;
    maxVideos: number;
    maxCommentsPerVideo: number;
  }) => {
    return request<{ jobId: string; status: string; estimatedCost: number }>('/analyze', {
      method: 'POST',
      body: JSON.stringify(params),
    });
  },

  // 轮询或获取任务详情与深度报告
  getJobDetail: (jobId: string) => {
    return request<OpptyJobDetailVO>(`/jobs/${jobId}`);
  },

  // 实时测算成本
  calculateCost: (params: { maxVideos: number; maxComments: number; deepCrawl: boolean }) => {
    return request<{ estimatedCost: number }>('/calculate-cost', {
      method: 'POST',
      body: JSON.stringify(params),
    });
  },

  // 获取商机历史列表
  getHistory: (pageNo = 1, pageSize = 10, keyword = '') => {
    const query = new URLSearchParams({
      pageNo: String(pageNo),
      pageSize: String(pageSize),
      keyword,
    });
    return request<{ total: number; pageNo: number; pageSize: number; records: OpptyHistoryItemVO[] }>(
      `/history?${query.toString()}`
    );
  },

  // 删除历史记录
  deleteHistory: (jobId: string) => {
    return request<void>(`/history/${jobId}`, {
      method: 'DELETE',
    });
  },

  // 公开免登录获取商机透视报告
  getShareJobDetail: (jobId: string) => {
    return request<OpptyJobDetailVO>(`/share/jobs/${jobId}`);
  },
};
