'use client';

import React, { useState, useEffect } from 'react';
import { playerApi, Player } from '../services/api';

interface PlayerDetailProps {
  playerId: string;
}

interface PlayerBalance {
  balance: number;
  frozen_balance: number;
  total_deposits: number;
  total_withdrawals: number;
  currency: string;
}

interface GameHistory {
  id: string;
  game_name: string;
  bet_amount: number;
  win_amount: number;
  created_at: string;
  status: string;
}

const PlayerDetail: React.FC<PlayerDetailProps> = ({ playerId }) => {
  const [player, setPlayer] = useState<Player | null>(null);
  const [balance, setBalance] = useState<PlayerBalance | null>(null);
  const [gameHistory, setGameHistory] = useState<GameHistory[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState('basic');

  useEffect(() => {
    fetchPlayerData();
  }, [playerId]);

  const fetchPlayerData = async () => {
    try {
      setLoading(true);
      setError(null);

      // 並行獲取玩家資訊、餘額和遊戲歷史
      const [playerResponse, balanceResponse, historyResponse] = await Promise.all([
        playerApi.getPlayer(playerId),
        playerApi.getPlayerBalance(playerId),
        playerApi.getPlayerGameHistory(playerId, { limit: 10 })
      ]);

      if (playerResponse.success) {
        setPlayer(playerResponse.data);
      } else {
        throw new Error(playerResponse.message || '獲取玩家資訊失敗');
      }

      if (balanceResponse.success) {
        setBalance(balanceResponse.data);
      }

      if (historyResponse.success) {
        setGameHistory(historyResponse.data?.data || []);
      }

    } catch (err) {
      console.error('獲取玩家資料失敗:', err);
      setError(err instanceof Error ? err.message : '獲取玩家資料時發生錯誤');
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'active':
        return 'bg-green-100 text-green-800';
      case 'inactive':
        return 'bg-yellow-100 text-yellow-800';
      case 'suspended':
        return 'bg-orange-100 text-orange-800';
      case 'banned':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getRiskLevelClass = (riskLevel: string) => {
    switch (riskLevel) {
      case 'low':
        return 'bg-green-100 text-green-800';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800';
      case 'high':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const formatDateTime = (dateString: string) => {
    return new Date(dateString).toLocaleString('zh-TW');
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('zh-TW', {
      style: 'currency',
      currency: 'TWD',
      minimumFractionDigits: 0,
    }).format(amount);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">載入玩家資料中...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-red-500 text-6xl mb-4">⚠️</div>
          <h2 className="text-xl font-semibold text-gray-900 mb-2">載入失敗</h2>
          <p className="text-gray-600 mb-4">{error}</p>
          <button
            onClick={fetchPlayerData}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            重新載入
          </button>
        </div>
      </div>
    );
  }

  if (!player) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-gray-900 mb-2">玩家未找到</h2>
          <p className="text-gray-600">找不到指定的玩家資料</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* 玩家基本資訊卡片 */}
        <div className="bg-white shadow-lg rounded-lg overflow-hidden mb-6">
          <div className="px-6 py-4 bg-gradient-to-r from-blue-600 to-blue-700">
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-2xl font-bold text-white">{player.full_name}</h1>
                <p className="text-blue-100">@{player.username}</p>
              </div>
              <div className="text-right">
                <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${getStatusBadgeClass(player.status)}`}>
                  {player.status === 'active' ? '啟用' : 
                   player.status === 'inactive' ? '未啟用' :
                   player.status === 'suspended' ? '暫停' : '封禁'}
                </span>
              </div>
            </div>
          </div>

          {/* 分頁導航 */}
          <div className="border-b border-gray-200">
            <nav className="-mb-px flex space-x-8 px-6">
              {[
                { id: 'basic', name: '基本資料' },
                { id: 'balance', name: '帳戶餘額' },
                { id: 'history', name: '遊戲歷史' },
                { id: 'analysis', name: '行為分析' }
              ].map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`py-4 px-1 border-b-2 font-medium text-sm ${
                    activeTab === tab.id
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  {tab.name}
                </button>
              ))}
            </nav>
          </div>

          {/* 分頁內容 */}
          <div className="p-6">
            {activeTab === 'basic' && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <div className="space-y-4">
                  <h3 className="text-lg font-medium text-gray-900">個人資料</h3>
                  <div className="space-y-3">
                    <div>
                      <dt className="text-sm font-medium text-gray-500">玩家 ID</dt>
                      <dd className="mt-1 text-sm text-gray-900">{player.id}</dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">使用者名稱</dt>
                      <dd className="mt-1 text-sm text-gray-900">{player.username}</dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">真實姓名</dt>
                      <dd className="mt-1 text-sm text-gray-900">{player.full_name}</dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">電子郵件</dt>
                      <dd className="mt-1 text-sm text-gray-900">{player.email}</dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">電話號碼</dt>
                      <dd className="mt-1 text-sm text-gray-900">{player.phone || '未提供'}</dd>
                    </div>
                  </div>
                </div>

                <div className="space-y-4">
                  <h3 className="text-lg font-medium text-gray-900">帳戶狀態</h3>
                  <div className="space-y-3">
                    <div>
                      <dt className="text-sm font-medium text-gray-500">帳戶狀態</dt>
                      <dd className="mt-1">
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusBadgeClass(player.status)}`}>
                          {player.status === 'active' ? '啟用' : 
                           player.status === 'inactive' ? '未啟用' :
                           player.status === 'suspended' ? '暫停' : '封禁'}
                        </span>
                      </dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">風險等級</dt>
                      <dd className="mt-1">
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getRiskLevelClass(player.risk_level)}`}>
                          {player.risk_level === 'low' ? '低風險' :
                           player.risk_level === 'medium' ? '中風險' : '高風險'}
                        </span>
                      </dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">VIP 等級</dt>
                      <dd className="mt-1 text-sm text-gray-900">等級 {player.vip_level}</dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">註冊時間</dt>
                      <dd className="mt-1 text-sm text-gray-900">{formatDateTime(player.registration_date)}</dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">最後登入</dt>
                      <dd className="mt-1 text-sm text-gray-900">{formatDateTime(player.last_login)}</dd>
                    </div>
                  </div>
                </div>

                <div className="space-y-4">
                  <h3 className="text-lg font-medium text-gray-900">遊戲統計</h3>
                  <div className="space-y-3">
                    <div>
                      <dt className="text-sm font-medium text-gray-500">遊戲場次</dt>
                      <dd className="mt-1 text-sm text-gray-900">{player.game_sessions?.toLocaleString() || 0} 次</dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">總儲值金額</dt>
                      <dd className="mt-1 text-sm text-gray-900">{formatCurrency(player.total_deposits || 0)}</dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">總提領金額</dt>
                      <dd className="mt-1 text-sm text-gray-900">{formatCurrency(player.total_withdrawals || 0)}</dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">目前餘額</dt>
                      <dd className="mt-1 text-sm font-semibold text-green-600">{formatCurrency(player.balance || 0)}</dd>
                    </div>
                    {player.tags && player.tags.length > 0 && (
                      <div>
                        <dt className="text-sm font-medium text-gray-500">標籤</dt>
                        <dd className="mt-1">
                          <div className="flex flex-wrap gap-1">
                            {player.tags.map((tag, index) => (
                              <span key={index} className="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-blue-100 text-blue-800">
                                {tag}
                              </span>
                            ))}
                          </div>
                        </dd>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'balance' && balance && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <div className="bg-gradient-to-r from-green-500 to-green-600 rounded-lg p-6 text-white">
                  <h3 className="text-lg font-medium">可用餘額</h3>
                  <p className="text-3xl font-bold">{formatCurrency(balance.balance)}</p>
                </div>
                <div className="bg-gradient-to-r from-yellow-500 to-yellow-600 rounded-lg p-6 text-white">
                  <h3 className="text-lg font-medium">凍結餘額</h3>
                  <p className="text-3xl font-bold">{formatCurrency(balance.frozen_balance)}</p>
                </div>
                <div className="bg-gradient-to-r from-blue-500 to-blue-600 rounded-lg p-6 text-white">
                  <h3 className="text-lg font-medium">總儲值</h3>
                  <p className="text-3xl font-bold">{formatCurrency(balance.total_deposits)}</p>
                </div>
                <div className="bg-gradient-to-r from-red-500 to-red-600 rounded-lg p-6 text-white">
                  <h3 className="text-lg font-medium">總提領</h3>
                  <p className="text-3xl font-bold">{formatCurrency(balance.total_withdrawals)}</p>
                </div>
              </div>
            )}

            {activeTab === 'history' && (
              <div>
                <h3 className="text-lg font-medium text-gray-900 mb-4">近期遊戲歷史</h3>
                {gameHistory.length > 0 ? (
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">遊戲名稱</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">下注金額</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">獲勝金額</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">時間</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">狀態</th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {gameHistory.map((game) => (
                          <tr key={game.id}>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{game.game_name}</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{formatCurrency(game.bet_amount)}</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{formatCurrency(game.win_amount)}</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDateTime(game.created_at)}</td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                                game.status === 'completed' ? 'bg-green-100 text-green-800' :
                                game.status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                                'bg-red-100 text-red-800'
                              }`}>
                                {game.status === 'completed' ? '已完成' :
                                 game.status === 'pending' ? '進行中' : '已取消'}
                              </span>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <p className="text-gray-500">暂无遊戲歷史記錄</p>
                  </div>
                )}
              </div>
            )}

            {activeTab === 'analysis' && (
              <div className="space-y-6">
                <h3 className="text-lg font-medium text-gray-900">行為分析</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="bg-gray-50 rounded-lg p-6">
                    <h4 className="text-md font-medium text-gray-900 mb-4">分析工具</h4>
                    <div className="space-y-3">
                      <button
                        onClick={() => playerApi.getPlayerBehaviorAnalysis(playerId)}
                        className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                      >
                        行為模式分析
                      </button>
                      <button
                        onClick={() => playerApi.getPlayerGamePreference(playerId)}
                        className="w-full px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700"
                      >
                        遊戲偏好分析
                      </button>
                      <button
                        onClick={() => playerApi.getPlayerSpendingHabits(playerId)}
                        className="w-full px-4 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700"
                      >
                        消費習慣分析
                      </button>
                      <button
                        onClick={() => playerApi.calculatePlayerValueScore(playerId)}
                        className="w-full px-4 py-2 bg-orange-600 text-white rounded-md hover:bg-orange-700"
                      >
                        價值評分計算
                      </button>
                    </div>
                  </div>
                  <div className="bg-gray-50 rounded-lg p-6">
                    <h4 className="text-md font-medium text-gray-900 mb-4">分析結果</h4>
                    <p className="text-gray-600">點擊左側按鈕以執行相應的分析</p>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default PlayerDetail; 