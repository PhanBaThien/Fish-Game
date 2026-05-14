import React, { useState } from 'react';
import { Save } from 'lucide-react';

export default function Settings() {
  const [settings, setSettings] = useState({
    gameName: 'Fish Game Pro 2026',
    feePercentage: '2.5',
    rtpAlertThreshold: '95',
    bannerUrl: 'https://example.com/banner.jpg',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSettings({
      ...settings,
      [e.target.name]: e.target.value,
    });
  };

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    alert('Đã lưu cấu hình thành công!');
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center p-4 bg-[#111827] border border-slate-800 rounded-lg">
        <p className="text-slate-400 text-xs">Cấu hình các thông số chung của toàn hệ thống game.</p>
        <button 
          onClick={handleSave}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded text-[10px] font-bold uppercase transition-colors flex items-center"
        >
          <Save className="w-4 h-4 mr-2" /> Lưu Cài Đặt
        </button>
      </div>

      <div className="bg-[#111827] p-6 rounded-lg border border-slate-800 max-w-2xl">
        <form onSubmit={handleSave} className="space-y-6">
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">Tên Game</label>
              <input 
                type="text" 
                name="gameName"
                value={settings.gameName}
                onChange={handleChange}
                className="w-full bg-[#0d1425] border border-slate-700 text-[#e2e8f0] rounded-lg px-4 py-2.5 focus:outline-none focus:border-blue-500 transition-colors"
              />
            </div>
            
            <div>
              <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">% Phí thu (Rake)</label>
              <div className="relative">
                <input 
                  type="number" 
                  name="feePercentage"
                  step="0.1"
                  value={settings.feePercentage}
                  onChange={handleChange}
                  className="w-full bg-[#0d1425] border border-slate-700 text-[#e2e8f0] rounded-lg px-4 py-2.5 focus:outline-none focus:border-blue-500 transition-colors"
                />
                <span className="absolute right-4 top-1/2 -translate-y-1/2 text-slate-500 font-bold">%</span>
              </div>
            </div>

            <div>
              <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">Ngưỡng cảnh báo RTP</label>
              <div className="relative">
                <input 
                  type="number" 
                  name="rtpAlertThreshold"
                  value={settings.rtpAlertThreshold}
                  onChange={handleChange}
                  className="w-full bg-[#0d1425] border border-slate-700 text-[#e2e8f0] rounded-lg px-4 py-2.5 focus:outline-none focus:border-blue-500 transition-colors"
                />
                <span className="absolute right-4 top-1/2 -translate-y-1/2 text-slate-500 font-bold">%</span>
              </div>
              <p className="text-[10px] text-slate-500 mt-1">Hệ thống sẽ gửi cảnh báo nếu RTP vượt mức này.</p>
            </div>
            
            <div className="md:col-span-2">
              <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">URL Banner Sự kiện</label>
              <input 
                type="text" 
                name="bannerUrl"
                value={settings.bannerUrl}
                onChange={handleChange}
                className="w-full bg-[#0d1425] border border-slate-700 text-[#e2e8f0] rounded-lg px-4 py-2.5 focus:outline-none focus:border-blue-500 transition-colors"
              />
            </div>
          </div>

        </form>
      </div>
    </div>
  );
}
