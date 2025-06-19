'use client';

import React, { useState, useEffect } from 'react';
import { Player, playerApi } from '../services/api';

interface PlayerDetailProps {
  playerId: string;
  onPlayerUpdate?: (player: Player) => void;
}

interface AnalysisData {
  behaviorAnalysis?: any;
  gamePreference?: any;
  spendingHabits?: any;
  valueScore?: any;
}

const PlayerDetail: React.FC<PlayerDetailProps> = ({ playerId, onPlayerUpdate }) => {
  const [player, setPlayer] = useState<Player | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [analysisData, setAnalysisData] = useState<AnalysisData>({});
  const [analysisLoading, setAnalysisLoading] = useState(false);
  const [activeTab, setActiveTab] = useState('basic');

  useEffect(() => {
    const fetchPlayer = async () => {
      try {
        setLoading(true);
        const response = await playerApi.getPlayer(playerId);
        setPlayer(response.data);
        onPlayerUpdate?.(response.data);
      } catch (err) {
        setError(err instanceof Error ? err.message : '載入玩家資料失敗');
      } finally {
        setLoading(false);
      }
    };

    fetchPlayer();
  }, [playerId, onPlayerUpdate]);

  // 載入分析資料
  const loadAnalysisData = async () => {
    if (!player || analysisLoading) return;
    
    try {
      setAnalysisLoading(true);
      const [behaviorRes, preferenceRes, spendingRes, valueRes] = await Promise.allSettled([
        playerApi.getPlayerBehaviorAnalysis(playerId),
        playerApi.getPlayerGamePreference(playerId),
        playerApi.getPlayerSpendingHabits(playerId),
        playerApi.calculatePlayerValueScore(playerId),
      ]);

      setAnalysisData({
        behaviorAnalysis: behaviorRes.status === 'fulfilled' ? behaviorRes.value.data : null,
        gamePreference: preferenceRes.status === 'fulfilled' ? preferenceRes.value.data : null,
        spendingHabits: spendingRes.status === 'fulfilled' ? spendingRes.value.data : null,
        valueScore: valueRes.status === 'fulfilled' ? valueRes.value.data : null,
      });
    } catch (err) {
      console.error('載入分析資料失敗:', err);
    } finally {
      setAnalysisLoading(false);
    }
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

  // 狀態顯示樣式
  const getStatusBadge = (status: string) => {
    const styles = {
      active: 'bg-green-100 text-green-800',
      inactive: 'bg-gray-100 text-gray-800',
      suspended: 'bg-yellow-100 text-yellow-800',
      deleted: 'bg-red-100 text-red-800'
    };
    return styles[status as keyof typeof styles] || 'bg-gray-100 text-gray-800';
  };

  // 風險等級顯示樣式
  const getRiskBadge = (riskLevel: string) => {
    const styles = {
      low: 'bg-green-100 text-green-800',
      medium: 'bg-yellow-100 text-yellow-800',
      high: 'bg-red-100 text-red-800',
      blacklist: 'bg-red-100 text-red-800'
    };
    return styles[riskLevel as keyof typeof styles] || 'bg-gray-100 text-gray-800';
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-md p-4">
        <p className="text-red-700">{error}</p>
      </div>
    );
  }

  if (!player) {
    return (
      <div className="bg-yellow-50 border border-yellow-200 rounded-md p-4">
        <p className="text-yellow-700">玩家資料不存在</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* 玩家基本資訊卡片 */}
      <div className="bg-gradient-to-r from-blue-600 to-blue-800 rounded-lg shadow-md text-white p-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-white">{player.real_name}</h1>
            <p className="text-blue-100">@{player.username}</p>
            <p className="text-blue-100 text-sm">{player.email}</p>
          </div>
          <div className="text-right">
            <p className="text-3xl font-bold">{formatAmount(player.balance)}</p>
            <p className="text-blue-100">目前餘額</p>
            <div className="mt-2 flex space-x-2">
              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusBadge(player.status)}`}>
                {player.status === 'active' ? '啟用' : 
                 player.status === 'inactive' ? '停用' : 
                 player.status === 'suspended' ? '暫停' : '刪除'}
              </span>
              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getRiskBadge(player.risk_level)}`}>
                {player.risk_level === 'low' ? '低風險' : 
                 player.risk_level === 'medium' ? '中風險' : 
                 player.risk_level === 'high' ? '高風險' : '黑名單'}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* 分頁標籤 */}
      <div className="border-b border-gray-200">
        <nav className="-mb-px flex space-x-8">
          {[
            { id: 'basic', name: '基本資料' },
            { id: 'account', name: '帳戶資訊' },
            { id: 'financial', name: '財務統計' },
            { id: 'analysis', name: '行為分析' },
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
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
      <div className="space-y-6">
        {activeTab === 'basic' && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* 基本資料 */}
            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">基本資料</h2>
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-gray-600">用戶名:</span>
                  <span className="font-medium">{player.username}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">真實姓名:</span>
                  <span className="font-medium">{player.real_name}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">郵箱:</span>
                  <span className="font-medium">{player.email}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">電話:</span>
                  <span className="font-medium">{player.phone}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">註冊時間:</span>
                  <span className="font-medium">{formatDate(player.created_at)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">最後登入:</span>
                  <span className="font-medium">{formatDate(player.last_login_at)}</span>
                </div>
              </div>
            </div>

            {/* 聯絡資訊 */}
            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">聯絡資訊</h2>
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-gray-600">語言:</span>
                  <span className="font-medium">{player.language || '繁體中文'}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">時區:</span>
                  <span className="font-medium">{player.timezone}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">登入次數:</span>
                  <span className="font-medium">{player.login_count} 次</span>
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'account' && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* 帳戶資訊 */}
            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">帳戶資訊</h2>
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-gray-600">VIP等級:</span>
                  <span className="font-medium">VIP {player.vip_level}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">驗證等級:</span>
                  <span className="font-medium">
                    {player.verification_level === 'none' ? '未驗證' : 
                     player.verification_level === 'email' ? '郵箱驗證' : 
                     player.verification_level === 'phone' ? '電話驗證' : '身份驗證'}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">風險等級:</span>
                  <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getRiskBadge(player.risk_level)}`}>
                    {player.risk_level === 'low' ? '低風險' : 
                     player.risk_level === 'medium' ? '中風險' : 
                     player.risk_level === 'high' ? '高風險' : '黑名單'}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">帳戶狀態:</span>
                  <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusBadge(player.status)}`}>
                    {player.status === 'active' ? '啟用' : 
                     player.status === 'inactive' ? '停用' : 
                     player.status === 'suspended' ? '暫停' : '刪除'}
                  </span>
                </div>
              </div>
            </div>

            {/* 帳戶操作 */}
            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">帳戶操作</h2>
              <div className="space-y-3">
                <button className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition-colors">
                  點數管理
                </button>
                <button className="w-full bg-yellow-600 text-white py-2 px-4 rounded-md hover:bg-yellow-700 transition-colors">
                  狀態管理
                </button>
                <button className="w-full bg-gray-600 text-white py-2 px-4 rounded-md hover:bg-gray-700 transition-colors">
                  設定限制
                </button>
                <button className="w-full bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700 transition-colors">
                  遊戲歷史
                </button>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'financial' && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {/* 財務統計 */}
            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">餘額資訊</h2>
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-gray-600">目前餘額:</span>
                  <span className="font-medium text-green-600">{formatAmount(player.balance)}</span>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">充值統計</h2>
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-gray-600">總充值:</span>
                  <span className="font-medium text-blue-600">{formatAmount(player.total_deposit)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">總提現:</span>
                  <span className="font-medium text-red-600">{formatAmount(player.total_withdraw)}</span>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">遊戲統計</h2>
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-gray-600">總投注:</span>
                  <span className="font-medium text-purple-600">{formatAmount(player.total_bet)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">總贏得:</span>
                  <span className="font-medium text-green-600">{formatAmount(player.total_win)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">淨輸贏:</span>
                  <span className={`font-medium ${player.total_win - player.total_bet >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                    {formatAmount(player.total_win - player.total_bet)}
                  </span>
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'analysis' && (
          <div className="space-y-6">
            {/* 載入分析資料按鈕 */}
            <div className="text-center">
              <button
                onClick={loadAnalysisData}
                disabled={analysisLoading}
                className="bg-blue-600 text-white py-2 px-6 rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50"
              >
                {analysisLoading ? '載入中...' : '載入行為分析'}
              </button>
            </div>

            {/* 分析結果 */}
            {Object.keys(analysisData).length > 0 && (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {analysisData.behaviorAnalysis && (
                  <div className="bg-white rounded-lg shadow-md p-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">行為分析</h3>
                    <pre className="text-sm text-gray-600 overflow-auto">
                      {JSON.stringify(analysisData.behaviorAnalysis, null, 2)}
                    </pre>
                  </div>
                )}

                {analysisData.gamePreference && (
                  <div className="bg-white rounded-lg shadow-md p-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">遊戲偏好</h3>
                    <pre className="text-sm text-gray-600 overflow-auto">
                      {JSON.stringify(analysisData.gamePreference, null, 2)}
                    </pre>
                  </div>
                )}

                {analysisData.spendingHabits && (
                  <div className="bg-white rounded-lg shadow-md p-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">消費習慣</h3>
                    <pre className="text-sm text-gray-600 overflow-auto">
                      {JSON.stringify(analysisData.spendingHabits, null, 2)}
                    </pre>
                  </div>
                )}

                {analysisData.valueScore && (
                  <div className="bg-white rounded-lg shadow-md p-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">價值評分</h3>
                    <pre className="text-sm text-gray-600 overflow-auto">
                      {JSON.stringify(analysisData.valueScore, null, 2)}
                    </pre>
                  </div>
                )}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default PlayerDetail; 