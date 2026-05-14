import React from 'react';
import { Settings } from 'lucide-react';

export default function Rooms() {
  const rooms = [
    { id: 'R-001', name: 'Biển Tân Thủ 1', type: 'Tân Thủ', bet: '10 Vàng/Đạn', players: 4, maxPlayers: 4, status: 'Đang chơi', baseRtp: 90 },
    { id: 'R-002', name: 'Biển Đại Dương 1', type: 'Nâng Cao', bet: '100 Vàng/Đạn', players: 2, maxPlayers: 4, status: 'Đang chờ', baseRtp: 85 },
    { id: 'R-003', name: 'Vịnh Thử Thách', type: 'Cao Thủ', bet: '1000 Vàng/Đạn', players: 4, maxPlayers: 4, status: 'Đang chơi', baseRtp: 80 },
    { id: 'R-004', name: 'Đảo Kho Báu (VIP)', type: 'VIP', bet: '10000 Vàng/Đạn', players: 1, maxPlayers: 4, status: 'Đang chờ', baseRtp: 88 },
    { id: 'R-005', name: 'Bão Táp Biển Sâu', type: 'Boss', bet: '5000 Vàng/Đạn', players: 0, maxPlayers: 4, status: 'Trống', baseRtp: 75 },
  ];

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center p-4 bg-[#111827] border border-slate-800 rounded-lg">
        <p className="text-slate-400 text-xs">Quản lý và giám sát các phòng bắn cá đang hoạt động. Điều chỉnh tỷ lệ ăn thua (RTP - Return to Player).</p>
        <button className="px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white rounded text-[10px] font-bold uppercase transition-colors">
          + Tạo Phòng Mới
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {rooms.map((room) => (
          <div key={room.id} className="bg-slate-800/40 border border-slate-700 rounded-lg overflow-hidden shadow-sm hover:border-slate-500 transition-all flex flex-col">
            <div className="p-4 border-b border-slate-800 flex justify-between items-start bg-[#111827]">
              <div>
                <h4 className="text-sm font-bold text-[#e2e8f0]">{room.name}</h4>
                <div className="text-[10px] text-slate-500 mt-0.5">ID: {room.id}</div>
              </div>
              <span className={`px-2 py-0.5 text-[9px] font-bold rounded border uppercase ${
                room.status === 'Đang chơi' ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 
                room.status === 'Đang chờ' ? 'bg-yellow-500/10 text-yellow-400 border-yellow-500/20' : 
                'bg-slate-500/10 text-slate-400 border-slate-500/20'
              }`}>
                {room.status}
              </span>
            </div>
            
            <div className="p-4 flex-1">
              <div className="space-y-2 mb-4">
                <div className="flex justify-between items-center text-xs">
                  <span className="text-slate-500 font-medium">Loại Phòng</span>
                  <span className="text-[#e2e8f0]">{room.type}</span>
                </div>
                <div className="flex justify-between items-center text-xs">
                  <span className="text-slate-500 font-medium">Mức Cược</span>
                  <span className="font-mono text-yellow-400">{room.bet}</span>
                </div>
                <div className="flex justify-between items-center text-xs bg-[#0d1425] p-2 rounded border border-slate-800 group">
                  <span className="text-slate-400 font-bold uppercase text-[10px]">Tỷ lệ trả thưởng (RTP)</span>
                  <div className="flex items-center gap-2">
                     <span className={`font-mono font-bold ${room.baseRtp > 100 ? 'text-red-400' : (room.baseRtp < 100 ? 'text-emerald-400' : 'text-slate-400')}`}>{room.baseRtp}%</span>
                     <button className="opacity-0 group-hover:opacity-100 text-slate-500 hover:text-blue-400 transition-colors">
                       <Settings className="w-3 h-3" />
                     </button>
                  </div>
                </div>
                <div className="flex justify-between items-center text-xs pt-2">
                  <span className="text-slate-500 font-medium">Người Chơi</span>
                  <div className="flex items-center gap-1.5 object-cover">
                    <span className="font-mono text-blue-400">{room.players}/{room.maxPlayers}</span>
                  </div>
                </div>
              </div>

              {/* Progress bar for players */}
              <div className="w-full h-1.5 bg-slate-800 rounded-full overflow-hidden">
                <div 
                  className={`h-full max-w-full ${
                    room.players === room.maxPlayers ? 'bg-emerald-500' : 'bg-blue-500'
                  }`} 
                  style={{ width: `${(room.players / room.maxPlayers) * 100}%` }}
                ></div>
              </div>
            </div>

            <div className="p-3 border-t border-slate-800 bg-[#0d1425] flex gap-2">
              <button className="flex-1 text-[10px] bg-blue-600/10 hover:bg-blue-600/20 border border-blue-500/20 py-1.5 rounded text-blue-400 font-bold uppercase transition-colors">
                Xem trực tiếp
              </button>
              <button className="flex-1 text-[10px] bg-slate-800 hover:bg-slate-700 border border-slate-700 py-1.5 rounded text-slate-300 font-bold uppercase transition-colors">
                Cài đặt Phòng
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
