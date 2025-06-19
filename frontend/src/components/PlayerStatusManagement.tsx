'use client';

import React, { useState, useEffect } from 'react';
import { Player, playerApi } from '../services/api';

interface PlayerStatusManagementProps {
  player: Player;
  onPlayerUpdate?: (player: Player) => void;
}

interface PlayerLimits {
  daily_bet_limit?: number;
  session_time_limit?: number;
  deposit_limit?: number;
  loss_limit?: number;
  bet_amount_limit?: number;
}

interface StatusHistory {
  id: number;
  old_status: string;
  new_status: string;
  reason: string;
  created_at: string;
  admin_user: string;
}

const PlayerStatusManagement: React.FC<PlayerStatusManagementProps> = ({
  player,
  onPlayerUpdate,
}) => {
  const [currentPlayer, setCurrentPlayer] = useState<Player>(player);
  const [selectedStatus, setSelectedStatus] = useState(player.status);
  const [statusReason, setStatusReason] = useState('');
  const [playerLimits, setPlayerLimits] = useState<PlayerLimits>({});
  const [statusHistory, setStatusHistory] = useState<StatusHistory[]>([]);
  const [loading, setLoading] = useState(false);
  const [historyLoading, setHistoryLoading] = useState(false);
  const [limitsLoading, setLimitsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // 狀態選項
  const statusOptions = [
    { value: 'active', label: '啟用', color: 'green', description: '玩家可以正常使用所有功能' },
    { value: 'inactive', label: '停用', color: 'gray', description: '玩家無法登入或進行任何操作' },
    { value: 'suspended', label: '暫停', color: 'yellow', description: '玩家暫時被限制使用部分功能' },
    { value: 'deleted', label: '刪除', color: 'red', description: '玩家帳戶已被刪除（不可恢復）' },
  ];

  // 風險等級選項
  const riskOptions = [
    { value: 'low', label: '低風險', color: 'green' },
    { value: 'medium', label: '中風險', color: 'yellow' },
    { value: 'high', label: '高風險', color: 'red' },
    { value: 'blacklist', label: '黑名單', color: 'red' },
  ];

  // 常見狀態變更原因
  const commonReasons = [
    '正常操作',
    '違反服務條款',
    '可疑活動檢測',
    '客服請求',
    '系統維護',
    '風險控制',
    '帳戶審核',
    '自我排除',
  ];

  // 格式化日期
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('zh-TW');
  };



  // 獲取狀態樣式
  const getStatusStyle = (status: string) => {
    const option = statusOptions.find(opt => opt.value === status);
    const colors = {
      green: 'bg-green-100 text-green-800',
      gray: 'bg-gray-100 text-gray-800',
      yellow: 'bg-yellow-100 text-yellow-800',
      red: 'bg-red-100 text-red-800',
    };
    return colors[option?.color as keyof typeof colors] || 'bg-gray-100 text-gray-800';
  };

  // 載入狀態歷史
  const loadStatusHistory = async () => {
    try {
      setHistoryLoading(true);
      // 模擬資料
      setStatusHistory([
        {
          id: 1,
          old_status: 'inactive',
          new_status: 'active',
          reason: '完成身份驗證',
          created_at: new Date().toISOString(),
          admin_user: 'admin'
        }
      ]);
    } catch (err) {
      console.error('載入狀態歷史失敗:', err);
    } finally {
      setHistoryLoading(false);
    }
  };

  // 更新玩家狀態
  const handleUpdateStatus = async () => {
    if (selectedStatus === currentPlayer.status) {
      setError('狀態未發生變更');
      return;
    }

    if (!statusReason.trim()) {
      setError('請填寫狀態變更原因');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      setSuccess(null);

      const response = await playerApi.updatePlayerStatus(player.player_id, selectedStatus);
      
      if (response.success) {
        const updatedPlayer = { ...currentPlayer, status: selectedStatus as Player['status'] };
        setCurrentPlayer(updatedPlayer);
        onPlayerUpdate?.(updatedPlayer);
        setSuccess('玩家狀態更新成功');
        setStatusReason('');
        loadStatusHistory();
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '更新狀態失敗');
    } finally {
      setLoading(false);
    }
  };

  // 設定玩家限制
  const handleSetLimits = async () => {
    try {
      setLimitsLoading(true);
      setError(null);
      setSuccess(null);

      const response = await playerApi.setPlayerLimits(player.player_id, playerLimits);
      
      if (response.success) {
        setSuccess('玩家限制設定成功');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '設定限制失敗');
    } finally {
      setLimitsLoading(false);
    }
  };

  useEffect(() => {
    loadStatusHistory();
  }, []);

  return (
    <div className="space-y-6">
      {/* 當前狀態顯示 */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">玩家狀態管理</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* 基本狀態資訊 */}
          <div>
            <h3 className="text-lg font-medium text-gray-900 mb-3">當前狀態</h3>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-gray-600">帳戶狀態:</span>
                <span className={`inline-flex px-3 py-1 text-sm font-semibold rounded-full ${getStatusStyle(currentPlayer.status)}`}>
                  {statusOptions.find(opt => opt.value === currentPlayer.status)?.label}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-gray-600">風險等級:</span>
                <span className={`inline-flex px-3 py-1 text-sm font-semibold rounded-full ${getStatusStyle(currentPlayer.risk_level)}`}>
                  {riskOptions.find(opt => opt.value === currentPlayer.risk_level)?.label}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-gray-600">VIP等級:</span>
                <span className="font-medium">VIP {currentPlayer.vip_level}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-gray-600">驗證等級:</span>
                <span className="font-medium">
                  {currentPlayer.verification_level === 'none' ? '未驗證' : 
                   currentPlayer.verification_level === 'email' ? '郵箱驗證' : 
                   currentPlayer.verification_level === 'phone' ? '電話驗證' : '身份驗證'}
                </span>
              </div>
            </div>
          </div>

          {/* 狀態統計 */}
          <div>
            <h3 className="text-lg font-medium text-gray-900 mb-3">帳戶統計</h3>
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="text-gray-600">登入次數:</span>
                <span className="font-medium">{currentPlayer.login_count} 次</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">最後登入:</span>
                <span className="font-medium">{formatDate(currentPlayer.last_login_at)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">註冊時間:</span>
                <span className="font-medium">{formatDate(currentPlayer.created_at)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">更新時間:</span>
                <span className="font-medium">{formatDate(currentPlayer.updated_at)}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 狀態變更 */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">變更狀態</h3>
        
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-md p-4 mb-4">
            <p className="text-red-700">{error}</p>
          </div>
        )}

        {success && (
          <div className="bg-green-50 border border-green-200 rounded-md p-4 mb-4">
            <p className="text-green-700">{success}</p>
          </div>
        )}

        <div className="space-y-4">
          {/* 狀態選擇 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              選擇新狀態
            </label>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {statusOptions.map((option) => (
                <div
                  key={option.value}
                  onClick={() => setSelectedStatus(option.value as Player['status'])}
                  className={`p-3 border rounded-md cursor-pointer transition-colors ${
                    selectedStatus === option.value
                      ? 'border-blue-500 bg-blue-50'
                      : 'border-gray-300 hover:border-gray-400'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusStyle(option.value)}`}>
                      {option.label}
                    </span>
                    <input
                      type="radio"
                      name="status"
                      value={option.value}
                      checked={selectedStatus === option.value}
                      onChange={() => setSelectedStatus(option.value as Player['status'])}
                      className="h-4 w-4 text-blue-600"
                    />
                  </div>
                  <p className="text-sm text-gray-600 mt-1">{option.description}</p>
                </div>
              ))}
            </div>
          </div>

          {/* 變更原因 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              變更原因
            </label>
            <select
              value={statusReason}
              onChange={(e) => setStatusReason(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent mb-2"
            >
              <option value="">請選擇變更原因</option>
              {commonReasons.map((reason) => (
                <option key={reason} value={reason}>
                  {reason}
                </option>
              ))}
              <option value="custom">其他原因</option>
            </select>
            
            {(statusReason === 'custom' || !commonReasons.includes(statusReason)) && (
              <input
                type="text"
                placeholder="請輸入自訂原因"
                value={statusReason === 'custom' ? '' : statusReason}
                onChange={(e) => setStatusReason(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            )}
          </div>

          {/* 變更預覽 */}
          {selectedStatus !== currentPlayer.status && (
            <div className="bg-yellow-50 border border-yellow-200 rounded-md p-4">
              <h4 className="font-medium text-yellow-900 mb-2">狀態變更預覽</h4>
              <div className="text-sm">
                <p className="mb-1">
                  <span className="text-gray-600">從:</span>
                  <span className={`ml-2 inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusStyle(currentPlayer.status)}`}>
                    {statusOptions.find(opt => opt.value === currentPlayer.status)?.label}
                  </span>
                </p>
                <p>
                  <span className="text-gray-600">變更為:</span>
                  <span className={`ml-2 inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusStyle(selectedStatus)}`}>
                    {statusOptions.find(opt => opt.value === selectedStatus)?.label}
                  </span>
                </p>
              </div>
            </div>
          )}

          {/* 提交按鈕 */}
          <button
            onClick={handleUpdateStatus}
            disabled={loading || selectedStatus === currentPlayer.status || !statusReason.trim()}
            className="w-full bg-blue-600 text-white py-3 px-4 rounded-md hover:bg-blue-700 transition-colors disabled:bg-gray-400"
          >
            {loading ? '更新中...' : '確認變更狀態'}
          </button>
        </div>
      </div>

      {/* 玩家限制設定 */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">玩家限制設定</h3>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* 每日投注限額 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              每日投注限額 (TWD)
            </label>
            <input
              type="number"
              min="0"
              value={playerLimits.daily_bet_limit || ''}
              onChange={(e) => setPlayerLimits(prev => ({ 
                ...prev, 
                daily_bet_limit: parseInt(e.target.value) || undefined 
              }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="設定每日投注上限"
            />
          </div>

          {/* 單次投注限額 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              單次投注限額 (TWD)
            </label>
            <input
              type="number"
              min="0"
              value={playerLimits.bet_amount_limit || ''}
              onChange={(e) => setPlayerLimits(prev => ({ 
                ...prev, 
                bet_amount_limit: parseInt(e.target.value) || undefined 
              }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="設定單次投注上限"
            />
          </div>

          {/* 每日充值限額 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              每日充值限額 (TWD)
            </label>
            <input
              type="number"
              min="0"
              value={playerLimits.deposit_limit || ''}
              onChange={(e) => setPlayerLimits(prev => ({ 
                ...prev, 
                deposit_limit: parseInt(e.target.value) || undefined 
              }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="設定每日充值上限"
            />
          </div>

          {/* 每日損失限額 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              每日損失限額 (TWD)
            </label>
            <input
              type="number"
              min="0"
              value={playerLimits.loss_limit || ''}
              onChange={(e) => setPlayerLimits(prev => ({ 
                ...prev, 
                loss_limit: parseInt(e.target.value) || undefined 
              }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="設定每日損失上限"
            />
          </div>

          {/* 遊戲時間限制 */}
          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              每日遊戲時間限制 (分鐘)
            </label>
            <input
              type="number"
              min="0"
              value={playerLimits.session_time_limit || ''}
              onChange={(e) => setPlayerLimits(prev => ({ 
                ...prev, 
                session_time_limit: parseInt(e.target.value) || undefined 
              }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="設定每日遊戲時間上限（分鐘）"
            />
          </div>
        </div>

        <button
          onClick={handleSetLimits}
          disabled={limitsLoading}
          className="mt-4 w-full bg-orange-600 text-white py-3 px-4 rounded-md hover:bg-orange-700 transition-colors disabled:bg-gray-400"
        >
          {limitsLoading ? '設定中...' : '儲存限制設定'}
        </button>
      </div>

      {/* 狀態變更歷史 */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-semibold text-gray-900">狀態變更歷史</h3>
          <button
            onClick={loadStatusHistory}
            disabled={historyLoading}
            className="bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50"
          >
            {historyLoading ? '載入中...' : '重新載入'}
          </button>
        </div>

        {historyLoading ? (
          <div className="flex justify-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
        ) : statusHistory.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    時間
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    原狀態
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    新狀態
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    原因
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    操作員
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {statusHistory.map((record) => (
                  <tr key={record.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {formatDate(record.created_at)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusStyle(record.old_status)}`}>
                        {statusOptions.find(opt => opt.value === record.old_status)?.label}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusStyle(record.new_status)}`}>
                        {statusOptions.find(opt => opt.value === record.new_status)?.label}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900">
                      {record.reason}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {record.admin_user}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="text-center py-8 text-gray-500">
            暫無狀態變更記錄
          </div>
        )}
      </div>
    </div>
  );
};

export default PlayerStatusManagement; 