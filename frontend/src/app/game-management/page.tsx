'use client';

import React, { useState } from 'react';
import { Gamepad2, Plus, Settings, BarChart3, Shield, DollarSign, Users, Play, Pause } from 'lucide-react';

// 模擬遊戲資料
const mockGames = [
  {
    id: 1,
    name: 'Texas Hold\'em',
    description: '經典德州撲克遊戲',
    type: 'poker',
    status: 'active',
    players: 156,
    revenue: 125600,
    winRate: 52.3,
    minBet: 10,
    maxBet: 5000,
    odds: 2.1
  },
  {
    id: 2,
    name: '百家樂',
    description: '傳統百家樂遊戲',
    type: 'baccarat',
    status: 'active',
    players: 89,
    revenue: 95400,
    winRate: 48.7,
    minBet: 50,
    maxBet: 10000,
    odds: 1.95
  },
  {
    id: 3,
    name: 'Five Card Stud',
    description: '五張牌梭哈遊戲',
    type: 'stud',
    status: 'maintenance',
    players: 0,
    revenue: 0,
    winRate: 0,
    minBet: 20,
    maxBet: 2000,
    odds: 2.5
  }
];

export default function GameManagementPage() {
  const [games, setGames] = useState(mockGames);

  const handleGameStatusToggle = (gameId: number) => {
    setGames(games.map(game => 
      game.id === gameId 
        ? { ...game, status: game.status === 'active' ? 'inactive' : 'active' }
        : game
    ));
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'text-green-600 bg-green-100';
      case 'inactive': return 'text-red-600 bg-red-100';
      case 'maintenance': return 'text-yellow-600 bg-yellow-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('zh-TW', {
      style: 'currency',
      currency: 'TWD',
      minimumFractionDigits: 0
    }).format(amount);
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* 頁面標題 */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
                <Gamepad2 className="h-8 w-8 text-blue-600" />
                遊戲管理系統
              </h1>
              <p className="text-gray-600 mt-2">管理遊戲設定、監控狀態與調整參數</p>
            </div>
            <button className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2">
              <Plus className="h-4 w-4" />
              新增遊戲
            </button>
          </div>
        </div>

        {/* 統計卡片 */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm">活躍遊戲</p>
                <p className="text-2xl font-bold text-gray-900">{games.filter(g => g.status === 'active').length}</p>
              </div>
              <Play className="h-8 w-8 text-green-600" />
            </div>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm">總玩家數</p>
                <p className="text-2xl font-bold text-gray-900">{games.reduce((sum, g) => sum + g.players, 0)}</p>
              </div>
              <Users className="h-8 w-8 text-blue-600" />
            </div>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm">今日營收</p>
                <p className="text-2xl font-bold text-gray-900">{formatCurrency(games.reduce((sum, g) => sum + g.revenue, 0))}</p>
              </div>
              <DollarSign className="h-8 w-8 text-green-600" />
            </div>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm">平均勝率</p>
                <p className="text-2xl font-bold text-gray-900">{(games.reduce((sum, g) => sum + g.winRate, 0) / games.length).toFixed(1)}%</p>
              </div>
              <BarChart3 className="h-8 w-8 text-purple-600" />
            </div>
          </div>
        </div>

        {/* 遊戲列表 */}
        <div className="bg-white rounded-lg shadow-sm">
          <div className="p-6 border-b border-gray-200">
            <h2 className="text-xl font-semibold text-gray-900">遊戲列表</h2>
          </div>
          
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">遊戲</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">狀態</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">玩家數</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">營收</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">勝率</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">下注範圍</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {games.map((game) => (
                  <tr key={game.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div>
                        <div className="text-sm font-medium text-gray-900">{game.name}</div>
                        <div className="text-sm text-gray-500">{game.description}</div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(game.status)}`}>
                        {game.status === 'active' ? '運行中' : game.status === 'inactive' ? '停用' : '維護中'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {game.players}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {formatCurrency(game.revenue)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {game.winRate}%
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {formatCurrency(game.minBet)} - {formatCurrency(game.maxBet)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <div className="flex items-center gap-2">
                        <button
                          onClick={() => handleGameStatusToggle(game.id)}
                          className={`p-2 rounded-lg transition-colors ${
                            game.status === 'active' 
                              ? 'text-red-600 hover:bg-red-100' 
                              : 'text-green-600 hover:bg-green-100'
                          }`}
                          title={game.status === 'active' ? '停用遊戲' : '啟用遊戲'}
                        >
                          {game.status === 'active' ? <Pause className="h-4 w-4" /> : <Play className="h-4 w-4" />}
                        </button>
                        <button
                          onClick={() => alert(`配置遊戲: ${game.name}`)}
                          className="p-2 text-blue-600 hover:bg-blue-100 rounded-lg transition-colors"
                          title="遊戲設定"
                        >
                          <Settings className="h-4 w-4" />
                        </button>
                        <button className="p-2 text-purple-600 hover:bg-purple-100 rounded-lg transition-colors" title="風險控制">
                          <Shield className="h-4 w-4" />
                        </button>
                        <button className="p-2 text-green-600 hover:bg-green-100 rounded-lg transition-colors" title="賠率管理">
                          <DollarSign className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* 開發提示 */}
        <div className="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-blue-900 mb-2">🚀 開發狀態</h3>
          <p className="text-blue-800 mb-4">
            根據 tasks-prd-game-management-backend.md，遊戲管理系統正在開發中。當前頁面為功能預覽。
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
            <div>
              <h4 className="font-semibold text-blue-900 mb-2">即將實現的功能:</h4>
              <ul className="space-y-1 text-blue-800">
                <li>• 遊戲新增/編輯管理</li>
                <li>• 遊戲狀態控制 (啟用/停用/維護)</li>
                <li>• 遊戲參數調整 (難度、AI強度)</li>
                <li>• 賠率管理系統</li>
                <li>• 下注限制設定</li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold text-blue-900 mb-2">監控功能:</h4>
              <ul className="space-y-1 text-blue-800">
                <li>• 即時遊戲狀態監控</li>
                <li>• 玩家數量統計</li>
                <li>• 投注金額即時統計</li>
                <li>• 營收狀況計算</li>
                <li>• 異常警報系統</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
} 