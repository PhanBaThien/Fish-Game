import React from 'react';
import { Fish } from 'lucide-react';

export default function FishConfig() {
  const fishes = [
    { id: 'F01', name: 'Cá Nhỏ (Xanh)', multiplier: 'x2', baseProb: '50.0%', speed: 'Nhanh', role: 'Thường' },
    { id: 'F02', name: 'Cá Đuối', multiplier: 'x15', baseProb: '6.6%', speed: 'Vừa', role: 'Thường' },
    { id: 'F03', name: 'Cá Mập', multiplier: 'x100', baseProb: '1.0%', speed: 'Chậm', role: 'Vừa' },
    { id: 'B01', name: 'Tiên Cá (Boss)', multiplier: 'x1000', baseProb: '0.1%', speed: 'Rất Chậm', role: 'Boss' },
  ];

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center p-4 bg-[#111827] border border-slate-800 rounded-lg">
        <p className="text-slate-400 text-xs">Cấu hình Hệ số nhân (Multiplier) và Tỷ lệ bắt gốc. Xác suất thực tế phụ thuộc RTP của vòng chơi.</p>
        <button className="px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white rounded text-[10px] font-bold uppercase transition-colors">
          + Thêm Cá Mới
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {fishes.map((f) => (
          <div key={f.id} className="bg-slate-800/40 border border-slate-700 rounded-lg overflow-hidden shadow-sm hover:border-slate-500 transition-all flex flex-col">
            <div className="h-28 bg-[#111827] flex items-center justify-center border-b border-slate-800 relative">
               <Fish className={`w-12 h-12 ${f.role === 'Boss' ? 'text-purple-400' : 'text-blue-400'}`} />
               <span className={`absolute top-2 right-2 px-1.5 py-0.5 text-[8px] font-bold rounded uppercase border ${
                 f.role === 'Boss' ? 'bg-purple-500/10 text-purple-400 border-purple-500/20' : 'bg-blue-500/10 text-blue-400 border-blue-500/20'
               }`}>
                  {f.role}
               </span>
            </div>
            <div className="p-3">
              <div className="text-center mb-3">
                <h4 className="text-xs font-bold text-[#e2e8f0]">{f.name}</h4>
                <div className="text-[9px] text-slate-500">ID: {f.id}</div>
              </div>
              <div className="space-y-1.5 text-xs">
                <div className="flex justify-between bg-[#0d1425] p-1.5 rounded border border-slate-800">
                  <span className="text-slate-500 text-[10px] uppercase font-bold tracking-tight">Hệ số nhân</span>
                  <span className="font-mono text-yellow-400 font-bold">{f.multiplier}</span>
                </div>
                <div className="flex justify-between bg-[#0d1425] p-1.5 rounded border border-slate-800">
                  <span className="text-slate-500 text-[10px] uppercase font-bold tracking-tight">Tỷ lệ nổ gốc</span>
                  <span className="font-mono text-emerald-400">{f.baseProb}</span>
                </div>
                <div className="flex justify-between bg-[#0d1425] p-1.5 rounded border border-slate-800">
                  <span className="text-slate-500 text-[10px] uppercase font-bold tracking-tight">Tốc độ</span>
                  <span className="font-mono text-blue-400">{f.speed}</span>
                </div>
              </div>
              <div className="mt-3 pt-3 border-t border-slate-800 flex justify-center gap-2">
                <button className="flex-1 text-[10px] bg-slate-800 hover:bg-slate-700 border border-slate-700 py-1.5 rounded text-slate-300 font-bold uppercase transition-colors">Sửa cấu hình</button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
