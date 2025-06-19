'use client';

import React, { useState } from 'react';
import PlayerList from '../../components/PlayerList';
import PlayerDetail from '../../components/PlayerDetail';
import PlayerPointsManagement from '../../components/PlayerPointsManagement';
import PlayerStatusManagement from '../../components/PlayerStatusManagement';
import { Player } from '../../services/api';

// 模擬玩家資料
const mockPlayer: Player = {
  id: 1,
  player_id: 'P001',
  username: 'testuser123',
  email: 'test@example.com',
  real_name: '測試玩家',
  phone: '0912345678',
  language: 'zh-TW',
  timezone: 'Asia/Taipei',
  status: 'active',
  verification_level: 'email',
  risk_level: 'low',
  vip_level: 1,
  last_login_at: new Date().toISOString(),
  login_count: 42,
  total_deposit: 50000,
  total_withdraw: 25000,
  total_bet: 100000,
  total_win: 85000,
  created_at: '2024-01-01T00:00:00Z',
  updated_at: new Date().toISOString(),
  balance: 12500,
};

export default function PlayerTestPage() {
  const [activeComponent, setActiveComponent] = useState('list');
  const [selectedPlayer, setSelectedPlayer] = useState<Player>(mockPlayer);

  const components = [
    { id: 'list', name: '玩家列表', component: 'PlayerList' },
    { id: 'detail', name: '玩家詳細', component: 'PlayerDetail' },
    { id: 'points', name: '點數管理', component: 'PlayerPointsManagement' },
    { id: 'status', name: '狀態管理', component: 'PlayerStatusManagement' },
  ];

  const handlePlayerUpdate = (updatedPlayer: Player) => {
    setSelectedPlayer(updatedPlayer);
    console.log('玩家資料已更新:', updatedPlayer);
  };

  const handleBalanceUpdate = (newBalance: number) => {
    setSelectedPlayer(prev => ({ ...prev, balance: newBalance }));
    console.log('餘額已更新:', newBalance);
  };

  const renderComponent = () => {
    switch (activeComponent) {
      case 'list':
        return (
          <div className="bg-white rounded-lg shadow-lg p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">玩家列表組件測試</h2>
            <div className="border border-gray-200 rounded-lg overflow-hidden">
              <PlayerList />
            </div>
          </div>
        );
      
      case 'detail':
        return (
          <div className="bg-white rounded-lg shadow-lg p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">玩家詳細組件測試</h2>
            <PlayerDetail 
              playerId={selectedPlayer.player_id} 
              onPlayerUpdate={handlePlayerUpdate}
            />
          </div>
        );
      
      case 'points':
        return (
          <div className="bg-white rounded-lg shadow-lg p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">點數管理組件測試</h2>
            <PlayerPointsManagement 
              player={selectedPlayer}
              onBalanceUpdate={handleBalanceUpdate}
            />
          </div>
        );
      
      case 'status':
        return (
          <div className="bg-white rounded-lg shadow-lg p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">狀態管理組件測試</h2>
            <PlayerStatusManagement 
              player={selectedPlayer}
              onPlayerUpdate={handlePlayerUpdate}
            />
          </div>
        );
      
      default:
        return null;
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 導航標題 */}
      <div className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold text-gray-900">
              🎮 NexusGaming - 玩家管理組件測試
            </h1>
            <div className="text-sm text-gray-500">
              當前測試玩家: {selectedPlayer.real_name} (@{selectedPlayer.username})
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
          {/* 側邊欄導航 */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg shadow-sm p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">組件測試導航</h3>
              <nav className="space-y-2">
                {components.map((comp) => (
                  <button
                    key={comp.id}
                    onClick={() => setActiveComponent(comp.id)}
                    className={`w-full text-left px-4 py-3 rounded-md text-sm font-medium transition-colors ${
                      activeComponent === comp.id
                        ? 'bg-blue-100 text-blue-700 border border-blue-200'
                        : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                    }`}
                  >
                    {comp.name}
                  </button>
                ))}
              </nav>

              {/* 玩家資訊卡片 */}
              <div className="mt-6 p-4 bg-gray-50 rounded-lg">
                <h4 className="text-sm font-semibold text-gray-700 mb-2">測試玩家資訊</h4>
                <div className="text-xs text-gray-600 space-y-1">
                  <p><span className="font-medium">ID:</span> {selectedPlayer.player_id}</p>
                  <p><span className="font-medium">餘額:</span> ${selectedPlayer.balance.toLocaleString()}</p>
                  <p><span className="font-medium">狀態:</span> 
                    <span className={`ml-1 px-2 py-0.5 rounded text-xs font-medium ${
                      selectedPlayer.status === 'active' ? 'bg-green-100 text-green-800' :
                      selectedPlayer.status === 'inactive' ? 'bg-gray-100 text-gray-800' :
                      selectedPlayer.status === 'suspended' ? 'bg-yellow-100 text-yellow-800' : 
                      'bg-red-100 text-red-800'
                    }`}>
                      {selectedPlayer.status === 'active' ? '啟用' : 
                       selectedPlayer.status === 'inactive' ? '停用' : 
                       selectedPlayer.status === 'suspended' ? '暫停' : '刪除'}
                    </span>
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* 主要內容區域 */}
          <div className="lg:col-span-3">
            {renderComponent()}
          </div>
        </div>
      </div>

      {/* 底部說明 */}
      <div className="bg-white border-t">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="text-center">
            <h3 className="text-lg font-semibold text-gray-900 mb-2">組件測試說明</h3>
            <div className="text-sm text-gray-600 space-y-1">
              <p>• 使用左側導航切換不同的組件進行測試</p>
              <p>• 當前使用模擬數據，實際部署時會連接真實API</p>
              <p>• 所有組件都支援響應式設計，可在不同設備上測試</p>
              <p>• 檢查瀏覽器控制台可查看組件的狀態變更日誌</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
} 