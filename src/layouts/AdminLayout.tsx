import React, { useState, useEffect } from 'react';
import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom';
import { 
  BarChart3, 
  Users, 
  Settings as SettingsIcon, 
  Fish, 
  Activity,
  LogOut,
  Search,
  Gamepad2,
  Play,
  ServerCrash,
  Shield
} from 'lucide-react';

type ApiStatus = 'checking' | 'online' | 'offline';

export default function AdminLayout() {
  const [apiStatus, setApiStatus] = useState<ApiStatus>('checking');
  const [apiUptime, setApiUptime] = useState<string | null>(null);
  const location = useLocation();
  const navigate = useNavigate();

  // Redirect to login if not authenticated
  useEffect(() => {
    const token = localStorage.getItem('fish_token');
    if (!token) {
      navigate('/login');
    }
  }, [navigate]);

  const handleLogout = () => {
    localStorage.removeItem('fish_token');
    navigate('/login');
  };

  // Live backend health check
  useEffect(() => {
    const checkHealth = async () => {
      try {
        const res = await fetch('/api/v1/health', { signal: AbortSignal.timeout(5000) });
        if (res.ok) {
          const data = await res.json();
          setApiStatus('online');
          setApiUptime(data.data?.uptime ?? null);
        } else {
          setApiStatus('offline');
        }
      } catch {
        setApiStatus('offline');
      }
    };

    checkHealth();
    const interval = setInterval(checkHealth, 30_000);
    return () => clearInterval(interval);
  }, []);

  const getPageTitle = () => {
    switch (location.pathname) {
      case '/': return 'Dashboard Thống kê';
      case '/players': return 'Quản lý Người chơi';
      case '/admins': return 'Quản lý Admin';
      case '/fish': return 'Cấu hình Chỉ số Cá 2D';
      case '/rooms': return 'Quản lý Phòng chơi';
      case '/settings': return 'Cài đặt Hệ thống';
      case '/gameplay': return 'Màn hình Game';
      default: return 'CMS';
    }
  };

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
          <NavItem to="/" icon={<Activity />} label="Tổng quan" active={location.pathname === '/'} />
          <NavItem to="/players" icon={<Users />} label="Quản lý Người chơi" active={location.pathname === '/players'} />
          <NavItem to="/admins" icon={<Shield />} label="Quản lý Admin" active={location.pathname === '/admins'} />
          <NavItem to="/fish" icon={<Fish />} label="Cấu hình Chỉ số Cá" active={location.pathname === '/fish'} />
          <NavItem to="/rooms" icon={<Gamepad2 />} label="Quản lý Phòng chơi" active={location.pathname === '/rooms'} />
          <NavItem to="/settings" icon={<SettingsIcon />} label="Cài đặt Hệ thống" active={location.pathname === '/settings'} />
          
          <div className="pt-4 text-[10px] font-bold text-slate-500 uppercase px-3 py-2 tracking-widest">Game Hub Ecosystem</div>
          <NavItem to="/gameplay" icon={<Play />} label="Vào màn hình Game" active={location.pathname === '/gameplay'} />
        </nav>
        
        <div className="mt-auto p-4 border-t border-slate-800 bg-[#0a0f1d]">
          <div className="p-2 bg-slate-800/50 rounded border border-slate-700/50 mb-3">
            <div className="text-[10px] text-slate-500 mb-1">SERVER STATUS</div>
            <div className={`flex items-center gap-1.5 text-xs font-semibold ${apiStatus === 'online' ? 'text-emerald-400' : 'text-red-400'}`}>
              <span className={`w-1.5 h-1.5 rounded-full animate-pulse ${apiStatus === 'online' ? 'bg-emerald-400' : 'bg-red-400'}`}></span>
              {apiStatus === 'online' ? 'Online' : 'Offline'}
            </div>
          </div>
          <button onClick={handleLogout} className="flex items-center w-full px-3 py-2 text-sm font-medium text-red-400 hover:bg-red-500/10 rounded-md transition-colors">
            <LogOut className="w-4 h-4 mr-3" />
            Đăng xuất
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col overflow-hidden w-full">
        {/* Topbar */}
        <div className="h-14 bg-[#111827] border-b border-slate-800 flex items-center justify-between px-4 shrink-0">
          <h2 className="text-sm font-semibold flex items-center gap-2 text-[#e2e8f0]">
            {getPageTitle()}
          </h2>
          <div className="flex items-center space-x-3">
            {/* Backend API Status Badge */}
            <div className={`flex items-center gap-1.5 px-2.5 py-1 rounded border text-[10px] font-bold uppercase tracking-wider ${
              apiStatus === 'online'
                ? 'bg-emerald-500/10 border-emerald-500/20 text-emerald-400'
                : apiStatus === 'offline'
                ? 'bg-red-500/10 border-red-500/20 text-red-400'
                : 'bg-slate-700/50 border-slate-600 text-slate-400'
            }`}>
              {apiStatus === 'online' ? (
                <><span className="w-1.5 h-1.5 bg-emerald-400 rounded-full animate-pulse" />API Online{apiUptime ? ` · ${apiUptime}` : ''}</>
              ) : apiStatus === 'offline' ? (
                <><ServerCrash className="w-3 h-3" />API Offline</>
              ) : (
                <><span className="w-1.5 h-1.5 bg-slate-400 rounded-full animate-pulse" />Connecting...</>
              )}
            </div>
            
            {/* Search Bar */}
            <div className="relative">
              <Search className="w-4 h-4 absolute left-2.5 top-1/2 transform -translate-y-1/2 text-slate-500" />
              <input
                type="text"
                placeholder="Tìm kiếm..."
                className="pl-8 pr-3 py-1.5 bg-slate-800 border-none rounded text-[10px] w-40 text-[#e2e8f0] placeholder-slate-500 outline-none ring-1 ring-slate-700 focus:ring-blue-500"
              />
            </div>
            <div className="w-8 h-8 rounded-full bg-blue-500/20 border border-blue-500/40 flex items-center justify-center text-[10px] font-bold text-blue-400">
              AD
            </div>
          </div>
        </div>

        {/* Content Outlet */}
        <div className="flex-1 overflow-y-auto p-4 bg-[#0a0f1d]">
          <Outlet />
        </div>
      </main>
    </div>
  );
}

function NavItem({ to, icon, label, active }: { to: string, icon: React.ReactNode, label: string, active: boolean }) {
  return (
    <Link
      to={to}
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
    </Link>
  );
}
