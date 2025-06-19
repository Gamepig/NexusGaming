'use client';

import React, { useState, useEffect } from 'react';
import { Player, playerApi } from '../services/api';

interface PlayerDetailProps {
  playerId: string;
}

const PlayerDetail: React.FC<PlayerDetailProps> = ({ playerId }) => {
  const [player, setPlayer] = useState<Player | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPlayer = async () => {
      try {
        setLoading(true);
        const response = await playerApi.getPlayer(playerId);
        setPlayer(response.data);
      } catch (err) {
        setError(err instanceof Error ? err.message : '載入玩家資料失敗');
      } finally {
        setLoading(false);
      }
    };

    fetchPlayer();
  }, [playerId]);

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
      {/* 玩家基本資訊 */}
      <div className="bg-gradient-to-r from-blue-600 to-blue-800 rounded-lg shadow-md text-white p-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-white">{player.real_name}</h1>
            <p className="text-blue-100">@{player.username}</p>
          </div>
          <div className="text-right">
            <p className="text-2xl font-bold">{formatAmount(player.balance)}</p>
            <p className="text-blue-100">目前餘額</p>
          </div>
        </div>
      </div>

      {/* 詳細資訊 */}
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
              <span className="text-gray-600">狀態:</span>
              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusBadge(player.status)}`}>
                {player.status === 'active' ? '啟用' : 
                 player.status === 'inactive' ? '停用' : 
                 player.status === 'suspended' ? '暫停' : '刪除'}
              </span>
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

        {/* 帳戶資訊 */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">帳戶資訊</h2>
          <div className="space-y-3">
            <div className="flex justify-between">
              <span className="text-gray-600">VIP等級:</span>
              <span className="font-medium">VIP {player.vip_level}</span>
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
              <span className="text-gray-600">驗證等級:</span>
              <span className="font-medium">
                {player.verification_level === 'none' ? '未驗證' : 
                 player.verification_level === 'email' ? '郵箱驗證' : 
                 player.verification_level === 'phone' ? '電話驗證' : '身份驗證'}
              </span>
            </div>
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

        {/* 財務統計 */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">財務統計</h2>
          <div className="space-y-3">
            <div className="flex justify-between">
              <span className="text-gray-600">目前餘額:</span>
              <span className="font-medium text-green-600">{formatAmount(player.balance)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">總充值:</span>
              <span className="font-medium">{formatAmount(player.total_deposit)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">總提現:</span>
              <span className="font-medium">{formatAmount(player.total_withdraw)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">總投注:</span>
              <span className="font-medium">{formatAmount(player.total_bet)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">總獲勝:</span>
              <span className="font-medium text-green-600">{formatAmount(player.total_win)}</span>
            </div>
            <div className="flex justify-between border-t pt-3">
              <span className="text-gray-600 font-semibold">盈虧:</span>
              <span className={`font-bold ${player.total_win - player.total_bet >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                {formatAmount(player.total_win - player.total_bet)}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PlayerDetail; 