import React from 'react';
import { Filter, Settings } from 'lucide-react';

export default function Players() {
  const players = [
    { id: 'USR001', name: 'HaiKute', gold: 45000, status: 'Active', login: '10 phút trước', winRate: '95%' },
    { id: 'USR002', name: 'SharkHunter', gold: 1200, status: 'Banned', login: '2 ngày trước', winRate: '120%' },
    { id: 'USR003', name: 'CaVangNo1', gold: 900000, status: 'Active', login: 'Vừa xong', winRate: '80%' },
    { id: 'USR004', name: 'TienNgoc', gold: 340, status: 'Active', login: '1 giờ trước', winRate: '100%' },
  ];

  return (
    <div className="bg-[#111827] rounded-lg border border-slate-800 overflow-hidden">
      <div className="p-3 border-b border-slate-800 flex justify-between items-center bg-[#0d1425]">
        <div className="flex space-x-2">
          <button className="px-3 py-1.5 bg-slate-800 border border-slate-700 rounded text-[10px] uppercase font-bold text-slate-300 flex items-center hover:bg-slate-700 transition-colors">
            <Filter className="w-3 h-3 mr-1.5" />
            Lọc
          </button>
        </div>
        <button className="px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white rounded text-[10px] font-bold uppercase transition-colors">
          + Thêm Quà Tặng
        </button>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-[#1f2937] border-b border-slate-800 text-[10px] uppercase tracking-widest text-slate-400">
              <th className="px-4 py-3 font-bold">ID</th>
              <th className="px-4 py-3 font-bold">Tên hiển thị</th>
              <th className="px-4 py-3 font-bold">Số Vàng</th>
              <th className="px-4 py-3 font-bold">Lần cuối đăng nhập</th>
              <th className="px-4 py-3 font-bold">Tỷ lệ thắng (RTP)</th>
              <th className="px-4 py-3 font-bold">Trạng thái</th>
              <th className="px-4 py-3 font-bold">Hành động</th>
            </tr>
          </thead>
          <tbody>
            {players.map((p) => {
               const rateNum = parseInt(p.winRate);
               const colorClass = rateNum > 100 ? 'text-red-400' : (rateNum < 100 ? 'text-emerald-400' : 'text-slate-400');
               return (
              <tr key={p.id} className="border-b border-slate-800/50 hover:bg-slate-800/30 transition-colors">
                <td className="px-4 py-3 text-xs font-mono text-slate-500">{p.id}</td>
                <td className="px-4 py-3 text-xs font-semibold text-[#e2e8f0]">{p.name}</td>
                <td className="px-4 py-3 text-xs text-yellow-400 font-bold">{p.gold.toLocaleString()}</td>
                <td className="px-4 py-3 text-xs text-slate-500">{p.login}</td>
                <td className="px-4 py-3">
                   <div className="flex items-center gap-2 cursor-pointer group">
                     <span className={`text-xs font-mono ${colorClass}`}>{p.winRate}</span>
                     <button className="opacity-0 group-hover:opacity-100 p-1 bg-slate-800 rounded hover:bg-slate-700 transition-all text-blue-400">
                        <Settings className="w-3 h-3" />
                     </button>
                   </div>
                </td>
                <td className="px-4 py-3">
                  <span className={`px-2 py-0.5 text-[9px] font-bold rounded border uppercase ${
                    p.status === 'Active' 
                      ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' 
                      : 'bg-red-500/10 text-red-400 border-red-500/20'
                  }`}>
                    {p.status}
                  </span>
                </td>
                <td className="px-4 py-3 text-[10px] font-bold uppercase space-x-3">
                  <button className="text-blue-400 hover:text-blue-300 transition-colors">Sửa</button>
                  <button className="text-red-400 hover:text-red-300 transition-colors">Khoá</button>
                </td>
              </tr>
            )})}
          </tbody>
        </table>
      </div>
    </div>
  );
}
