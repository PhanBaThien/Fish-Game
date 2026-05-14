import React from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Users, DollarSign, BarChart3, Fish } from 'lucide-react';

const data = [
  { name: 'T2', players: 4000, revenue: 2400 },
  { name: 'T3', players: 3000, revenue: 1398 },
  { name: 'T4', players: 2000, revenue: 9800 },
  { name: 'T5', players: 2780, revenue: 3908 },
  { name: 'T6', players: 1890, revenue: 4800 },
  { name: 'T7', players: 2390, revenue: 3800 },
  { name: 'CN', players: 3490, revenue: 4300 },
];

export default function Dashboard() {
  return (
    <div className="space-y-6">
      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <StatCard title="Người chơi hiện tại" value="1,245" trend="+12%" icon={<Users />} color="text-blue-400 border-blue-500/30" />
        <StatCard title="Doanh thu hôm nay" value="23,400.00 đ" trend="+5.4%" icon={<DollarSign />} color="text-emerald-400 border-emerald-500/30" />
        <StatCard title="Lượt tải mới" value="342" trend="-2%" icon={<BarChart3 />} color="text-purple-400 border-purple-500/30" />
        <StatCard title="Số cá bị tiêu diệt" value="1.2M" trend="+24%" icon={<Fish />} color="text-orange-400 border-orange-500/30" />
      </div>

      {/* Chart */}
      <div className="bg-[#111827] p-5 rounded-lg border border-slate-800">
        <h3 className="text-sm font-semibold flex items-center gap-2 mb-4 text-[#e2e8f0]">Biểu đồ người chơi & Doanh thu 7 ngày qua</h3>
        <div className="h-80 w-full">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
              <defs>
                <linearGradient id="colorPlayers" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.8}/>
                  <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
                </linearGradient>
                <linearGradient id="colorRevenue" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#10b981" stopOpacity={0.8}/>
                  <stop offset="95%" stopColor="#10b981" stopOpacity={0}/>
                </linearGradient>
              </defs>
              <XAxis dataKey="name" stroke="#64748b" tick={{fill: '#64748b', fontSize: 10}} />
              <YAxis stroke="#64748b" tick={{fill: '#64748b', fontSize: 10}} />
              <CartesianGrid stroke="#1e293b" strokeDasharray="3 3" vertical={false} />
              <Tooltip contentStyle={{ backgroundColor: '#0d1425', borderColor: '#1e293b', color: '#e2e8f0', fontSize: '10px' }} />
              <Area type="monotone" dataKey="players" stroke="#3b82f6" fillOpacity={1} fill="url(#colorPlayers)" />
              <Area type="monotone" dataKey="revenue" stroke="#10b981" fillOpacity={1} fill="url(#colorRevenue)" />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
}

function StatCard({ title, value, trend, icon, color }: { title: string, value: string, trend: string, icon: React.ReactNode, color: string }) {
  const isPositive = trend.startsWith('+');
  return (
    <div className="bg-slate-800/40 border border-slate-700 rounded-lg p-4 flex items-start justify-between hover:border-slate-500 transition-colors">
      <div>
        <p className="text-[10px] font-bold tracking-widest uppercase text-slate-500 mb-1">{title}</p>
        <h4 className="text-xl font-bold text-[#e2e8f0]">{value}</h4>
        <div className={`text-[10px] mt-2 font-semibold ${isPositive ? 'text-emerald-400' : 'text-red-400'}`}>
          {trend} so với hôm qua
        </div>
      </div>
      <div className={`p-2 rounded bg-slate-800/50 border ${color} [&>svg]:w-5 [&>svg]:h-5`}>
        {icon}
      </div>
    </div>
  );
}
