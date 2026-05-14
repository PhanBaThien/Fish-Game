import React from 'react';
import { Fish, DollarSign } from 'lucide-react';

export default function GamePlay() {
  return (
    <div className="w-full h-[600px] border-2 border-blue-500/30 bg-[#0a0f1d] rounded-2xl overflow-hidden relative shadow-2xl shadow-blue-900/20">
      <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1544551763-46a013bb70d5?q=80&w=2670&auto=format&fit=crop')] bg-cover bg-center opacity-30"></div>
      
      {/* HUD overlays for game view */}
      <div className="absolute top-4 left-4 right-4 flex justify-between z-10">
        <div className="bg-slate-900/80 border border-slate-700 px-4 py-2 rounded-full flex items-center gap-3 backdrop-blur-sm">
          <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-xs font-bold ring-2 ring-blue-300">AD</div>
          <div>
            <div className="text-[10px] text-blue-300 uppercase font-bold tracking-wider">Admin Player</div>
            <div className="text-yellow-400 font-mono text-sm font-bold flex items-center gap-1">
              <DollarSign className="w-4 h-4" /> 9,999,999
            </div>
          </div>
        </div>
        
        <div className="bg-slate-900/80 border border-slate-700 px-4 py-2 rounded-full flex gap-4 backdrop-blur-sm shadow-lg">
          <div className="flex flex-col items-center justify-center">
            <span className="text-[9px] text-slate-400 uppercase tracking-widest font-bold">Phòng</span>
            <span className="text-xs text-white font-bold">VIP 1</span>
          </div>
          <div className="w-px h-full bg-slate-700"></div>
          <div className="flex flex-col items-center justify-center">
            <span className="text-[9px] text-slate-400 uppercase tracking-widest font-bold">Cược</span>
            <span className="text-xs text-yellow-400 font-bold">10,000</span>
          </div>
        </div>
      </div>

      <div className="absolute inset-0 flex items-center justify-center z-0">
          <div className="text-center animate-pulse">
            <Fish className="w-24 h-24 text-blue-400/50 mx-auto mb-4" />
            <h3 className="text-2xl font-bold text-white/50 uppercase tracking-widest">Khu Vực Bắn Cá</h3>
            <p className="text-blue-300/50 text-sm mt-2">Đang kết nối WebSocket Server...</p>
          </div>
      </div>
      
      <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex items-end gap-16 z-10">
        {/* Cannons */}
        <div className="flex flex-col items-center gap-2">
          <div className="bg-slate-900/80 px-3 py-1 rounded text-yellow-400 font-mono text-xs font-bold border border-slate-700">100,200</div>
          <div className="w-16 h-20 bg-gradient-to-t from-slate-800 to-slate-600 rounded-lg flex items-center justify-center border-2 border-slate-500 relative">
            <div className="w-4 h-12 bg-blue-500 rounded-full absolute -top-4 shadow-[0_0_15px_rgba(59,130,246,0.5)]"></div>
          </div>
        </div>
        
        <div className="flex flex-col items-center gap-2 opacity-50">
          <div className="bg-slate-900/80 px-3 py-1 rounded text-yellow-400 font-mono text-xs font-bold border border-slate-700">Trống</div>
          <div className="w-16 h-20 bg-gradient-to-t from-slate-800 to-slate-700 rounded-lg flex items-center justify-center border-2 border-slate-700">
             <span className="text-slate-500 text-[10px] uppercase font-bold">Ngồi</span>
          </div>
        </div>
        
        <div className="flex flex-col items-center gap-2">
          <div className="bg-slate-900/80 px-3 py-1 rounded text-yellow-400 font-mono text-xs font-bold border border-slate-700">9,999,999</div>
          <div className="w-16 h-20 bg-gradient-to-t from-orange-800 to-yellow-600 rounded-lg flex items-center justify-center border-2 border-yellow-500 relative shadow-[0_0_20px_rgba(234,179,8,0.3)]">
            <div className="w-6 h-14 bg-gradient-to-t from-yellow-500 to-white rounded-full absolute -top-6 shadow-[0_0_20px_rgba(255,255,255,0.8)]"></div>
          </div>
        </div>
        
        <div className="flex flex-col items-center gap-2 opacity-50">
          <div className="bg-slate-900/80 px-3 py-1 rounded text-yellow-400 font-mono text-xs font-bold border border-slate-700">Trống</div>
          <div className="w-16 h-20 bg-gradient-to-t from-slate-800 to-slate-700 rounded-lg flex items-center justify-center border-2 border-slate-700">
             <span className="text-slate-500 text-[10px] uppercase font-bold">Ngồi</span>
          </div>
        </div>
      </div>
    </div>
  );
}
