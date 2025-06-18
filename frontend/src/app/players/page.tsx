'use client';

import React from 'react';
import PlayerList from '../../components/PlayerList';
import { Player } from '../../services/api';

export default function PlayersPage() {
  const handlePlayerSelect = (player: Player) => {
    console.log('Selected player:', player);
    // 這裡可以導航到玩家詳細頁面或打開詳細資訊模態框
    // 例如：router.push(`/players/${player.id}`);
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="container mx-auto px-4 py-8">
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">玩家管理系統</h1>
          <p className="text-gray-600 mt-2">
            管理和監控所有註冊玩家的資訊、狀態與活動
          </p>
        </div>
        
        <PlayerList onPlayerSelect={handlePlayerSelect} />
      </div>
    </div>
  );
} 