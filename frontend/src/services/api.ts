// API 服務層 - 處理所有後端 API 調用

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// API 回應介面
interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  timestamp?: string;
}

// 玩家介面定義
export interface Player {
  id: number;
  player_id: string;
  username: string;
  email: string;
  real_name: string;
  phone: string;
  language: string;
  timezone: string;
  status: 'active' | 'inactive' | 'suspended' | 'deleted';
  verification_level: 'none' | 'email' | 'phone' | 'identity';
  risk_level: 'low' | 'medium' | 'high' | 'blacklist';
  vip_level: number;
  last_login_at: string;
  login_count: number;
  total_deposit: number;
  total_withdraw: number;
  total_bet: number;
  total_win: number;
  created_at: string;
  updated_at: string;
  balance: number;
}

// 分頁介面
export interface PaginationParams {
  page?: number;
  limit?: number;
  sort?: string;
  order?: 'asc' | 'desc';
}

// 搜尋參數介面
export interface PlayerSearchParams extends PaginationParams {
  search?: string;
  status?: string;
  startDate?: string;
  endDate?: string;
  minBalance?: number;
  maxBalance?: number;
  risk_level?: string;
  vip_level?: number;
}

// 分頁回應介面
export interface PaginatedResponse<T> {
  players: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
    has_next: boolean;
    has_previous: boolean;
  };
  filters: {
    search: string;
    status: string;
    start_date: string;
    end_date: string;
    risk_level: string;
    verification_level: string;
  };
}

// HTTP 客戶端類別
class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseURL}${endpoint}`;
    
    const defaultHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    // 如果有 token，加入 Authorization header
    const token = localStorage.getItem('auth_token');
    if (token) {
      defaultHeaders['Authorization'] = `Bearer ${token}`;
    }

    const config: RequestInit = {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  async get<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, data?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async put<T>(endpoint: string, data?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

// API 客戶端實例
const apiClient = new ApiClient(API_BASE_URL);

// 玩家相關 API 方法
export const playerApi = {
  // 獲取玩家列表
  async getPlayers(params: PlayerSearchParams = {}): Promise<ApiResponse<PaginatedResponse<Player>>> {
    const queryParams = new URLSearchParams();
    
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        queryParams.append(key, value.toString());
      }
    });

    const endpoint = `/players/${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return apiClient.get<PaginatedResponse<Player>>(endpoint);
  },

  // 獲取單一玩家詳細資訊
  async getPlayer(id: string): Promise<ApiResponse<Player>> {
    return apiClient.get<Player>(`/players/${id}`);
  },

  // 獲取玩家遊戲歷史
  async getPlayerGameHistory(id: string, params: PaginationParams = {}): Promise<ApiResponse<unknown>> {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        queryParams.append(key, value.toString());
      }
    });

    const endpoint = `/players/${id}/game-history${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return apiClient.get(endpoint);
  },

  // 獲取玩家餘額
  async getPlayerBalance(id: string): Promise<ApiResponse<unknown>> {
    return apiClient.get(`/players/${id}/balance`);
  },

  // 更新玩家狀態
  async updatePlayerStatus(id: string, status: string): Promise<ApiResponse<unknown>> {
    return apiClient.put(`/players/${id}/status`, { status });
  },

  // 設定玩家限制
  async setPlayerLimits(id: string, limits: unknown): Promise<ApiResponse<unknown>> {
    return apiClient.put(`/players/${id}/limits`, limits);
  },

  // 玩家點數調整
  async adjustPlayerBalance(id: string, adjustment: { amount: number; type: 'add' | 'subtract'; reason: string }): Promise<ApiResponse<unknown>> {
    return apiClient.post(`/players/${id}/balance/adjust`, adjustment);
  },

  // 獲取玩家行為分析
  async getPlayerBehaviorAnalysis(id: string): Promise<ApiResponse<unknown>> {
    return apiClient.post(`/players/${id}/behavior-analysis`);
  },

  // 獲取玩家遊戲偏好統計
  async getPlayerGamePreference(id: string): Promise<ApiResponse<unknown>> {
    return apiClient.post(`/players/${id}/game-preference`);
  },

  // 獲取玩家消費習慣分析
  async getPlayerSpendingHabits(id: string): Promise<ApiResponse<unknown>> {
    return apiClient.post(`/players/${id}/spending-habits`);
  },

  // 計算玩家價值評分
  async calculatePlayerValueScore(id: string, config?: unknown): Promise<ApiResponse<unknown>> {
    return apiClient.post(`/players/${id}/value-score`, config);
  },
};

// 身份驗證相關 API
export const authApi = {
  async login(credentials: { username: string; password: string }): Promise<ApiResponse<{ token: string; user: unknown }>> {
    return apiClient.post('/auth/login', credentials);
  },

  async logout(): Promise<ApiResponse<unknown>> {
    return apiClient.post('/auth/logout');
  },

  async refreshToken(): Promise<ApiResponse<{ token: string }>> {
    return apiClient.post('/auth/refresh');
  },
};

export default apiClient; 