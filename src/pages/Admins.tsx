import React, { useState } from 'react';
import { Filter, Settings, Shield, ShieldAlert, Edit, Trash2, Plus } from 'lucide-react';

export default function AdminsTab() {
  const [admins] = useState([
    { id: 'ADM001', username: 'superadmin', email: 'admin@fishgame.com', role: 'super_admin', createdAt: '2023-01-10 08:00', status: 'Active' },
    { id: 'ADM002', username: 'operator1', email: 'op1@fishgame.com', role: 'admin', createdAt: '2023-05-15 14:30', status: 'Active' },
    { id: 'ADM003', username: 'mod_minh', email: 'minh.mod@fishgame.com', role: 'moderator', createdAt: '2023-11-20 09:15', status: 'Inactive' },
    { id: 'ADM004', username: 'operator2', email: 'op2@fishgame.com', role: 'admin', createdAt: '2024-02-01 16:45', status: 'Active' },
  ]);

  return (
    <div className="bg-[#111827] rounded-lg border border-slate-800 overflow-hidden">
      <div className="p-3 border-b border-slate-800 flex justify-between items-center bg-[#0d1425]">
        <div className="flex space-x-2">
          <button className="px-3 py-1.5 bg-slate-800 border border-slate-700 rounded text-[10px] uppercase font-bold text-slate-300 flex items-center hover:bg-slate-700 transition-colors">
            <Filter className="w-3 h-3 mr-1.5" />
            Lọc
          </button>
        </div>
        <button className="px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white rounded text-[10px] font-bold uppercase transition-colors flex items-center">
          <Plus className="w-3 h-3 mr-1" /> Thêm Admin
        </button>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-[#1f2937] border-b border-slate-800 text-[10px] uppercase tracking-widest text-slate-400">
              <th className="px-4 py-3 font-bold">ID</th>
              <th className="px-4 py-3 font-bold">Username</th>
              <th className="px-4 py-3 font-bold">Email</th>
              <th className="px-4 py-3 font-bold">Vai trò (Role)</th>
              <th className="px-4 py-3 font-bold">Trạng thái</th>
              <th className="px-4 py-3 font-bold">Ngày tạo</th>
              <th className="px-4 py-3 font-bold text-right">Hành động</th>
            </tr>
          </thead>
          <tbody>
            {admins.map((admin) => (
              <tr key={admin.id} className="border-b border-slate-800/50 hover:bg-slate-800/30 transition-colors">
                <td className="px-4 py-3 text-xs font-mono text-slate-500">{admin.id}</td>
                <td className="px-4 py-3 text-xs font-semibold text-[#e2e8f0]">{admin.username}</td>
                <td className="px-4 py-3 text-xs text-slate-400">{admin.email}</td>
                <td className="px-4 py-3">
                  <div className="flex items-center gap-1.5">
                    {admin.role === 'super_admin' ? (
                      <ShieldAlert className="w-3.5 h-3.5 text-purple-400" />
                    ) : admin.role === 'admin' ? (
                      <Shield className="w-3.5 h-3.5 text-blue-400" />
                    ) : (
                      <Settings className="w-3.5 h-3.5 text-slate-400" />
                    )}
                    <span className={`text-[10px] font-bold uppercase ${
                      admin.role === 'super_admin' ? 'text-purple-400' : 
                      admin.role === 'admin' ? 'text-blue-400' : 'text-slate-400'
                    }`}>
                      {admin.role.replace('_', ' ')}
                    </span>
                  </div>
                </td>
                <td className="px-4 py-3">
                  <span className={`px-2 py-0.5 text-[9px] font-bold rounded border uppercase ${
                    admin.status === 'Active' 
                      ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' 
                      : 'bg-slate-500/10 text-slate-400 border-slate-500/20'
                  }`}>
                    {admin.status}
                  </span>
                </td>
                <td className="px-4 py-3 text-xs text-slate-500">{admin.createdAt}</td>
                <td className="px-4 py-3 text-[10px] font-bold uppercase text-right space-x-2">
                  <button className="p-1.5 text-blue-400 hover:bg-blue-500/10 rounded transition-colors" title="Sửa">
                    <Edit className="w-4 h-4" />
                  </button>
                  <button className="p-1.5 text-red-400 hover:bg-red-500/10 rounded transition-colors" title="Xóa">
                    <Trash2 className="w-4 h-4" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
