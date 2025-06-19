'use client';

import React, { useState, useEffect } from 'react';
import { Player, playerApi } from '../services/api';

interface PlayerPointsManagementProps {
  player: Player;
  onBalanceUpdate?: (newBalance: number) => void;
}

interface BalanceAdjustment {
  amount: number;
  type: 'add' | 'subtract';
  reason: string;
}

interface BalanceHistory {
  id: number;
  amount: number;
  type: string;
  reason: string;
  balance_before: number;
  balance_after: number;
  created_at: string;
  admin_user: string;
}

const PlayerPointsManagement: React.FC<PlayerPointsManagementProps> = ({
  player,
  onBalanceUpdate,
}) => {
  const [currentBalance, setCurrentBalance] = useState(player.balance);
  const [adjustmentForm, setAdjustmentForm] = useState<BalanceAdjustment>({
    amount: 0,
    type: 'add',
    reason: '',
  });
  const [balanceHistory, setBalanceHistory] = useState<BalanceHistory[]>([]);
  const [loading, setLoading] = useState(false);
  const [historyLoading, setHistoryLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // 格式化金額
  const formatAmount = (amount: number) => {
    return new Intl.NumberFormat('zh-TW', {
      style: 'currency',
      currency: 'TWD',
      minimumFractionDigits: 0
    }).format(amount);
  };

  // 格式化日期
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('zh-TW');
  };

  // 載入餘額歷史
  const loadBalanceHistory = async () => {
    try {
      setHistoryLoading(true);
      // 這裡應該調用實際的歷史記錄 API
      // const response = await playerApi.getPlayerBalanceHistory(player.id);
      // setBalanceHistory(response.data);
      
      // 模擬資料
      setBalanceHistory([
        {
          id: 1,
          amount: 1000,
          type: 'deposit',
          reason: '充值',
          balance_before: 0,
          balance_after: 1000,
          created_at: new Date().toISOString(),
          admin_user: 'admin'
        }
      ]);
    } catch (err) {
      console.error('載入餘額歷史失敗:', err);
    } finally {
      setHistoryLoading(false);
    }
  };

  // 調整餘額
  const handleAdjustBalance = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (adjustmentForm.amount <= 0) {
      setError('金額必須大於0');
      return;
    }

    if (!adjustmentForm.reason.trim()) {
      setError('請填寫調整原因');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      setSuccess(null);

      const response = await playerApi.adjustPlayerBalance(player.player_id, adjustmentForm);
      
      if (response.success) {
        const newBalance = adjustmentForm.type === 'add' 
          ? currentBalance + adjustmentForm.amount
          : currentBalance - adjustmentForm.amount;
        
        setCurrentBalance(newBalance);
        onBalanceUpdate?.(newBalance);
        setSuccess('餘額調整成功');
        
        // 重置表單
        setAdjustmentForm({
          amount: 0,
          type: 'add',
          reason: '',
        });

        // 重新載入歷史記錄
        loadBalanceHistory();
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '調整餘額失敗');
    } finally {
      setLoading(false);
    }
  };

  // 預設調整原因
  const commonReasons = [
    '系統獎勵',
    '活動贈送',
    '客服補償', 
    '錯誤扣款補回',
    '違規扣款',
    '系統錯誤調整',
    '推薦獎勵',
    '首儲獎勵',
  ];

  useEffect(() => {
    loadBalanceHistory();
  }, []);

  return (
    <div className="space-y-6">
      {/* 當前餘額顯示 */}
      <div className="bg-gradient-to-r from-green-600 to-green-800 rounded-lg shadow-md text-white p-6">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-green-100 mb-2">當前餘額</h2>
          <p className="text-3xl font-bold">{formatAmount(currentBalance)}</p>
          <p className="text-green-100 text-sm mt-2">玩家: {player.real_name} (@{player.username})</p>
        </div>
      </div>

      {/* 餘額調整表單 */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">餘額調整</h3>
        
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

        <form onSubmit={handleAdjustBalance} className="space-y-4">
          {/* 調整類型 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              調整類型
            </label>
            <div className="grid grid-cols-2 gap-4">
              <button
                type="button"
                onClick={() => setAdjustmentForm(prev => ({ ...prev, type: 'add' }))}
                className={`py-3 px-4 rounded-md text-center font-medium transition-colors ${
                  adjustmentForm.type === 'add'
                    ? 'bg-green-600 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                ➕ 增加點數
              </button>
              <button
                type="button"
                onClick={() => setAdjustmentForm(prev => ({ ...prev, type: 'subtract' }))}
                className={`py-3 px-4 rounded-md text-center font-medium transition-colors ${
                  adjustmentForm.type === 'subtract'
                    ? 'bg-red-600 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                ➖ 扣除點數
              </button>
            </div>
          </div>

          {/* 調整金額 */}
          <div>
            <label htmlFor="amount" className="block text-sm font-medium text-gray-700 mb-2">
              調整金額
            </label>
            <div className="relative">
              <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">$</span>
              <input
                type="number"
                id="amount"
                min="1"
                step="1"
                value={adjustmentForm.amount || ''}
                onChange={(e) => setAdjustmentForm(prev => ({ 
                  ...prev, 
                  amount: parseInt(e.target.value) || 0 
                }))}
                className="w-full pl-8 pr-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="請輸入調整金額"
                required
              />
            </div>
          </div>

          {/* 調整原因 */}
          <div>
            <label htmlFor="reason" className="block text-sm font-medium text-gray-700 mb-2">
              調整原因
            </label>
            <select
              value={adjustmentForm.reason}
              onChange={(e) => setAdjustmentForm(prev => ({ ...prev, reason: e.target.value }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent mb-2"
            >
              <option value="">請選擇調整原因</option>
              {commonReasons.map((reason) => (
                <option key={reason} value={reason}>
                  {reason}
                </option>
              ))}
              <option value="custom">其他原因</option>
            </select>
            
            {(adjustmentForm.reason === 'custom' || !commonReasons.includes(adjustmentForm.reason)) && (
              <input
                type="text"
                placeholder="請輸入自訂原因"
                value={adjustmentForm.reason === 'custom' ? '' : adjustmentForm.reason}
                onChange={(e) => setAdjustmentForm(prev => ({ ...prev, reason: e.target.value }))}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              />
            )}
          </div>

          {/* 預覽 */}
          {adjustmentForm.amount > 0 && (
            <div className="bg-gray-50 rounded-md p-4">
              <h4 className="font-medium text-gray-900 mb-2">調整預覽</h4>
              <div className="text-sm text-gray-600 space-y-1">
                <p>目前餘額: {formatAmount(currentBalance)}</p>
                <p className={adjustmentForm.type === 'add' ? 'text-green-600' : 'text-red-600'}>
                  {adjustmentForm.type === 'add' ? '增加' : '扣除'}: {formatAmount(adjustmentForm.amount)}
                </p>
                <p className="font-medium">
                  調整後餘額: {formatAmount(
                    adjustmentForm.type === 'add' 
                      ? currentBalance + adjustmentForm.amount
                      : currentBalance - adjustmentForm.amount
                  )}
                </p>
              </div>
            </div>
          )}

          {/* 提交按鈕 */}
          <button
            type="submit"
            disabled={loading || adjustmentForm.amount <= 0 || !adjustmentForm.reason.trim()}
            className={`w-full py-3 px-4 rounded-md font-medium text-white transition-colors ${
              adjustmentForm.type === 'add'
                ? 'bg-green-600 hover:bg-green-700 disabled:bg-gray-400'
                : 'bg-red-600 hover:bg-red-700 disabled:bg-gray-400'
            }`}
          >
            {loading ? '處理中...' : (adjustmentForm.type === 'add' ? '確認增加點數' : '確認扣除點數')}
          </button>
        </form>
      </div>

      {/* 餘額歷史記錄 */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-semibold text-gray-900">餘額變動歷史</h3>
          <button
            onClick={loadBalanceHistory}
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
        ) : balanceHistory.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    時間
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    類型
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    金額
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    調整前
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    調整後
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
                {balanceHistory.map((record) => (
                  <tr key={record.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {formatDate(record.created_at)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        record.type === 'add' || record.type === 'deposit' 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-red-100 text-red-800'
                      }`}>
                        {record.type === 'add' ? '增加' : 
                         record.type === 'subtract' ? '扣除' : 
                         record.type === 'deposit' ? '充值' : '提現'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <span className={record.amount > 0 ? 'text-green-600' : 'text-red-600'}>
                        {record.amount > 0 ? '+' : ''}{formatAmount(record.amount)}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {formatAmount(record.balance_before)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {formatAmount(record.balance_after)}
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
            暫無餘額變動記錄
          </div>
        )}
      </div>
    </div>
  );
};

export default PlayerPointsManagement; 