'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Player, PlayerSearchParams, playerApi } from '../services/api';

interface PlayerListProps {
  onPlayerSelect?: (player: Player) => void;
}

const PlayerList: React.FC<PlayerListProps> = ({ onPlayerSelect }) => {
  const router = useRouter();
  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [searchParams, setSearchParams] = useState<PlayerSearchParams>({
    page: 1,
    limit: 20,
    sort: 'created_at',
    order: 'desc'
  });
  const [totalPages, setTotalPages] = useState(1);
  const [totalCount, setTotalCount] = useState(0);

  // 搜尋表單狀態
  const [searchForm, setSearchForm] = useState({
    search: '',
    status: '',
    startDate: '',
    endDate: '',
    minBalance: '',
    maxBalance: ''
  });

  // 載入玩家列表
  const loadPlayers = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const params: PlayerSearchParams = {
        ...searchParams,
        search: searchForm.search || undefined,
        status: searchForm.status || undefined,
        startDate: searchForm.startDate || undefined,
        endDate: searchForm.endDate || undefined,
        minBalance: searchForm.minBalance ? parseFloat(searchForm.minBalance) : undefined,
        maxBalance: searchForm.maxBalance ? parseFloat(searchForm.maxBalance) : undefined,
      };

      const response = await playerApi.getPlayers(params);
      setPlayers(response.data.data);
      setTotalPages(response.data.pagination?.totalPages || 1);
      setTotalCount(response.data.pagination?.totalCount || 0);
    } catch (err) {
      setError(err instanceof Error ? err.message : '載入玩家列表失敗');
    } finally {
      setLoading(false);
    }
  };

  // 初始載入和參數變更時重新載入
  useEffect(() => {
    loadPlayers();
  }, [searchParams]);

  // 處理搜尋
  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchParams(prev => ({ ...prev, page: 1 }));
  };

  // 處理排序
  const handleSort = (field: string) => {
    setSearchParams(prev => ({
      ...prev,
      sort: field,
      order: prev.sort === field && prev.order === 'asc' ? 'desc' : 'asc',
      page: 1
    }));
  };

  // 處理分頁
  const handlePageChange = (page: number) => {
    setSearchParams(prev => ({ ...prev, page }));
  };

  // 狀態顯示樣式
  const getStatusBadge = (status: string) => {
    const styles = {
      active: 'bg-green-100 text-green-800',
      inactive: 'bg-gray-100 text-gray-800',
      suspended: 'bg-yellow-100 text-yellow-800',
      banned: 'bg-red-100 text-red-800'
    };
    return styles[status as keyof typeof styles] || 'bg-gray-100 text-gray-800';
  };

  // 風險等級顯示樣式
  const getRiskBadge = (riskLevel: string) => {
    const styles = {
      low: 'bg-green-100 text-green-800',
      medium: 'bg-yellow-100 text-yellow-800',
      high: 'bg-red-100 text-red-800'
    };
    return styles[riskLevel as keyof typeof styles] || 'bg-gray-100 text-gray-800';
  };

  // 格式化日期
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('zh-TW', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  // 格式化金額
  const formatAmount = (amount: number) => {
    return new Intl.NumberFormat('zh-TW', {
      style: 'currency',
      currency: 'TWD',
      minimumFractionDigits: 0
    }).format(amount);
  };

  return (
    <div className="bg-white rounded-lg shadow-md">
      {/* 標題 */}
      <div className="px-6 py-4 border-b border-gray-200">
        <h2 className="text-xl font-semibold text-gray-900">玩家管理</h2>
        <p className="text-sm text-gray-600 mt-1">
          共 {totalCount} 位玩家
        </p>
      </div>

      {/* 搜尋表單 */}
      <div className="px-6 py-4 bg-gray-50 border-b border-gray-200">
        <form onSubmit={handleSearch} className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {/* 關鍵字搜尋 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                搜尋關鍵字
              </label>
              <input
                type="text"
                placeholder="用戶名、郵箱或姓名"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={searchForm.search}
                onChange={(e) => setSearchForm(prev => ({ ...prev, search: e.target.value }))}
              />
            </div>

            {/* 狀態篩選 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                玩家狀態
              </label>
              <select
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={searchForm.status}
                onChange={(e) => setSearchForm(prev => ({ ...prev, status: e.target.value }))}
              >
                <option value="">全部狀態</option>
                <option value="active">啟用</option>
                <option value="inactive">停用</option>
                <option value="suspended">暫停</option>
                <option value="banned">封鎖</option>
              </select>
            </div>

            {/* 註冊日期範圍 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                註冊日期（起）
              </label>
              <input
                type="date"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={searchForm.startDate}
                onChange={(e) => setSearchForm(prev => ({ ...prev, startDate: e.target.value }))}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                註冊日期（迄）
              </label>
              <input
                type="date"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={searchForm.endDate}
                onChange={(e) => setSearchForm(prev => ({ ...prev, endDate: e.target.value }))}
              />
            </div>

            {/* 餘額範圍 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                最低餘額
              </label>
              <input
                type="number"
                placeholder="0"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={searchForm.minBalance}
                onChange={(e) => setSearchForm(prev => ({ ...prev, minBalance: e.target.value }))}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                最高餘額
              </label>
              <input
                type="number"
                placeholder="無限制"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={searchForm.maxBalance}
                onChange={(e) => setSearchForm(prev => ({ ...prev, maxBalance: e.target.value }))}
              />
            </div>
          </div>

          <div className="flex justify-end space-x-2">
            <button
              type="button"
              onClick={() => {
                setSearchForm({
                  search: '',
                  status: '',
                  startDate: '',
                  endDate: '',
                  minBalance: '',
                  maxBalance: ''
                });
                setSearchParams(prev => ({ ...prev, page: 1 }));
              }}
              className="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              清除
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              搜尋
            </button>
          </div>
        </form>
      </div>

      {/* 錯誤訊息 */}
      {error && (
        <div className="px-6 py-4 bg-red-50 border-l-4 border-red-400">
          <p className="text-red-700">{error}</p>
        </div>
      )}

      {/* 載入中 */}
      {loading && (
        <div className="px-6 py-8 text-center">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <p className="mt-2 text-gray-600">載入中...</p>
        </div>
      )}

      {/* 玩家列表 */}
      {!loading && (
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                  onClick={() => handleSort('username')}
                >
                  用戶名
                  {searchParams.sort === 'username' && (
                    <span className="ml-1">
                      {searchParams.order === 'asc' ? '↑' : '↓'}
                    </span>
                  )}
                </th>
                <th
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                  onClick={() => handleSort('real_name')}
                >
                  姓名
                  {searchParams.sort === 'real_name' && (
                    <span className="ml-1">
                      {searchParams.order === 'asc' ? '↑' : '↓'}
                    </span>
                  )}
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  郵箱
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  狀態
                </th>
                <th
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                  onClick={() => handleSort('balance')}
                >
                  餘額
                  {searchParams.sort === 'balance' && (
                    <span className="ml-1">
                      {searchParams.order === 'asc' ? '↑' : '↓'}
                    </span>
                  )}
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  風險等級
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  VIP等級
                </th>
                <th
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                  onClick={() => handleSort('created_at')}
                >
                  註冊時間
                  {searchParams.sort === 'created_at' && (
                    <span className="ml-1">
                      {searchParams.order === 'asc' ? '↑' : '↓'}
                    </span>
                  )}
                </th>
                <th
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                  onClick={() => handleSort('last_login_at')}
                >
                  最後登入
                  {searchParams.sort === 'last_login_at' && (
                    <span className="ml-1">
                      {searchParams.order === 'asc' ? '↑' : '↓'}
                    </span>
                  )}
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  操作
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {players.map((player) => (
                <tr key={player.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">
                      {player.username}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900">{player.real_name}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900">{player.email}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusBadge(player.status)}`}>
                      {player.status === 'active' ? '啟用' : 
                       player.status === 'inactive' ? '停用' : 
                       player.status === 'suspended' ? '暫停' : '封鎖'}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">
                      {formatAmount(player.balance)}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getRiskBadge(player.risk_level)}`}>
                      {player.risk_level === 'low' ? '低風險' : 
                       player.risk_level === 'medium' ? '中風險' : '高風險'}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900">VIP {player.vip_level}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900">
                      {formatDate(player.created_at)}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900">
                      {formatDate(player.last_login_at)}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button
                      onClick={() => router.push(`/players/${player.id}`)}
                      className="text-blue-600 hover:text-blue-900 mr-3"
                    >
                      查看
                    </button>
                    <button className="text-green-600 hover:text-green-900 mr-3">
                      編輯
                    </button>
                    <button className="text-red-600 hover:text-red-900">
                      管理
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* 無資料提示 */}
      {!loading && players.length === 0 && (
        <div className="px-6 py-8 text-center text-gray-500">
          沒有找到符合條件的玩家
        </div>
      )}

      {/* 分頁 */}
      {!loading && totalPages > 1 && (
        <div className="px-6 py-4 bg-gray-50 border-t border-gray-200">
          <div className="flex items-center justify-between">
            <div className="text-sm text-gray-700">
              顯示第 {((searchParams.page || 1) - 1) * (searchParams.limit || 20) + 1} - {Math.min((searchParams.page || 1) * (searchParams.limit || 20), totalCount)} 筆，共 {totalCount} 筆
            </div>
            <div className="flex space-x-1">
              <button
                onClick={() => handlePageChange((searchParams.page || 1) - 1)}
                disabled={(searchParams.page || 1) <= 1}
                className="px-3 py-2 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                上一頁
              </button>
              
              {/* 頁碼按鈕 */}
              {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                const page = Math.max(1, Math.min(totalPages - 4, (searchParams.page || 1) - 2)) + i;
                return (
                  <button
                    key={page}
                    onClick={() => handlePageChange(page)}
                    className={`px-3 py-2 text-sm border rounded-md ${
                      page === (searchParams.page || 1)
                        ? 'bg-blue-600 text-white border-blue-600'
                        : 'border-gray-300 hover:bg-gray-50'
                    }`}
                  >
                    {page}
                  </button>
                );
              })}
              
              <button
                onClick={() => handlePageChange((searchParams.page || 1) + 1)}
                disabled={(searchParams.page || 1) >= totalPages}
                className="px-3 py-2 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                下一頁
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default PlayerList; 