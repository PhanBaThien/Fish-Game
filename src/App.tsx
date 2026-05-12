import { useState } from 'react';
import { 
  BarChart3, 
  Users, 
  Settings, 
  Fish, 
  Activity,
  LogOut,
  Search,
  Filter,
  DollarSign,
  Gamepad2,
  Play
} from 'lucide-react';
import { 
  AreaChart, 
  Area, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  ResponsiveContainer 
} from 'recharts';

const data = [
  { name: 'T2', players: 4000, revenue: 2400 },
  { name: 'T3', players: 3000, revenue: 1398 },
  { name: 'T4', players: 2000, revenue: 9800 },
  { name: 'T5', players: 2780, revenue: 3908 },
  { name: 'T6', players: 1890, revenue: 4800 },
  { name: 'T7', players: 2390, revenue: 3800 },
  { name: 'CN', players: 3490, revenue: 4300 },
];

export default function App() {
  const [activeTab, setActiveTab] = useState('dashboard');

  return (
    <div className="flex h-screen bg-[#0a0f1d] text-[#e2e8f0] font-sans selection:bg-blue-500/30 overflow-hidden">
      {/* Sidebar */}
      <aside className="w-60 bg-[#0d1425] border-r border-slate-800 flex flex-col shrink-0">
        <div className="h-14 flex items-center px-4 border-b border-slate-800 shrink-0">
          <div className="w-8 h-8 bg-blue-600 rounded flex items-center justify-center shadow-lg shadow-blue-900/20 mr-3">
            <Fish className="w-5 h-5 text-white" />
          </div>
          <div className="flex flex-col">
            <span className="font-bold text-sm tracking-tight text-[#e2e8f0]">FISHGAME ADMIN</span>
            <span className="text-[10px] text-slate-500 uppercase font-medium">Management Console</span>
          </div>
        </div>
        
        <nav className="flex-1 p-3 space-y-1">
          <div className="text-[10px] font-bold text-slate-500 uppercase px-3 py-2 tracking-widest">Main Menu</div>
          <NavItem 
            icon={<Activity />} 
            label="Tổng quan (Dashboard)" 
            active={activeTab === 'dashboard'} 
            onClick={() => setActiveTab('dashboard')} 
          />
          <NavItem 
            icon={<Users />} 
            label="Quản lý Người chơi" 
            active={activeTab === 'players'} 
            onClick={() => setActiveTab('players')} 
          />
          <NavItem 
            icon={<Fish />} 
            label="Cấu hình Chỉ số Cá" 
            active={activeTab === 'fish'} 
            onClick={() => setActiveTab('fish')} 
          />
          <NavItem 
            icon={<Gamepad2 />} 
            label="Quản lý Phòng chơi" 
            active={activeTab === 'rooms'} 
            onClick={() => setActiveTab('rooms')} 
          />
          <NavItem 
            icon={<Settings />} 
            label="Cài đặt Hệ thống" 
            active={activeTab === 'settings'} 
            onClick={() => setActiveTab('settings')} 
          />
          
          <div className="pt-4 text-[10px] font-bold text-slate-500 uppercase px-3 py-2 tracking-widest">Game Hub Ecosystem</div>
          <NavItem 
            icon={<Play />} 
            label="Vào màn hình Game" 
            active={activeTab === 'gameplay'} 
            onClick={() => setActiveTab('gameplay')} 
          />
        </nav>
        
        <div className="mt-auto p-4 border-t border-slate-800 bg-[#0a0f1d]">
          <div className="p-2 bg-slate-800/50 rounded border border-slate-700/50 mb-3">
            <div className="text-[10px] text-slate-500 mb-1">SERVER STATUS</div>
            <div className="flex items-center gap-1.5 text-xs font-semibold text-emerald-400">
              <span className="w-1.5 h-1.5 bg-emerald-400 rounded-full animate-pulse"></span>
              Online (99.8%)
            </div>
          </div>
          <button className="flex items-center w-full px-3 py-2 text-sm font-medium text-red-400 hover:bg-red-500/10 rounded-md transition-colors">
            <LogOut className="w-4 h-4 mr-3" />
            Đăng xuất
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col overflow-hidden w-full">
        <div className="h-14 bg-[#111827] border-b border-slate-800 flex items-center justify-between px-4 shrink-0">
          <h2 className="text-sm font-semibold flex items-center gap-2 text-[#e2e8f0]">
            {activeTab === 'dashboard' && 'Dashboard Thống kê'}
            {activeTab === 'players' && 'Quản lý Người chơi'}
            {activeTab === 'fish' && 'Cấu hình Chỉ số Cá 2D'}
            {activeTab === 'rooms' && 'Quản lý Phòng chơi'}
            {activeTab === 'settings' && 'Cài đặt'}
            {activeTab === 'gameplay' && 'Màn hình Game'}
          </h2>
          <div className="flex items-center space-x-4">
            <div className="relative">
              <Search className="w-4 h-4 absolute left-2.5 top-1/2 transform -translate-y-1/2 text-slate-500" />
              <input 
                type="text" 
                placeholder="Tìm kiếm..." 
                className="pl-8 pr-3 py-1.5 bg-slate-800 border-none rounded text-[10px] w-40 text-[#e2e8f0] placeholder-slate-500 outline-none ring-1 ring-slate-700 focus:ring-blue-500"
              />
            </div>
            <button className="p-1.5 bg-slate-800 hover:bg-slate-700 rounded transition-colors text-slate-400">
              <Activity className="w-4 h-4" />
            </button>
            <div className="w-8 h-8 rounded-full bg-blue-500/20 border border-blue-500/40 flex items-center justify-center text-[10px] font-bold text-blue-400">
              AD
            </div>
          </div>
        </div>

        <div className="flex-1 overflow-y-auto p-4 bg-[#0a0f1d]">
          {activeTab === 'dashboard' && <DashboardTab />}
          {activeTab === 'players' && <PlayersTab />}
          {activeTab === 'fish' && <FishConfigTab />}
          {activeTab === 'rooms' && <RoomsTab />}
          {activeTab === 'gameplay' && <GamePlayPlaceholder />}
          {activeTab === 'settings' && (
            <div className="bg-[#111827] p-5 rounded-lg border border-slate-800">
              <p className="text-slate-500 text-xs text-center p-8 border border-dashed border-slate-700 rounded">Cài đặt hệ thống đang được cập nhật...</p>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}

function NavItem({ icon, label, active, onClick }: { icon: React.ReactNode, label: string, active: boolean, onClick: () => void }) {
  return (
    <button
      onClick={onClick}
      className={`flex items-center gap-3 w-full px-3 py-2 text-sm font-medium rounded-md transition-colors ${
        active 
          ? 'bg-blue-600/10 text-blue-400 border border-blue-500/20' 
          : 'text-slate-400 hover:bg-slate-800 border border-transparent'
      }`}
    >
      <div className={`[&>svg]:w-4 [&>svg]:h-4 ${active ? 'text-blue-400' : 'text-slate-400'}`}>
        {icon}
      </div>
      {label}
    </button>
  );
}

function DashboardTab() {
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

function PlayersTab() {
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

function FishConfigTab() {
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

function RoomsTab() {
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

function GamePlayPlaceholder() {
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


