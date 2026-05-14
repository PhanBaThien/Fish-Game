import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Fish, KeyRound, User } from 'lucide-react';

export default function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault();
    // Logic đăng nhập ảo: Chỉ cần username/password giống nhau hoặc admin/admin là pass
    if (username === 'admin' && password === 'admin') {
      localStorage.setItem('fish_token', 'dummy_token');
      navigate('/');
    } else {
      setError('Tài khoản hoặc mật khẩu không đúng (Gợi ý: admin/admin)');
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-[#0a0f1d]">
      <div className="w-full max-w-md p-8 bg-[#111827] rounded-xl border border-slate-800 shadow-2xl">
        <div className="flex flex-col items-center mb-8">
          <div className="w-16 h-16 bg-blue-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-900/40 mb-4">
            <Fish className="w-10 h-10 text-white" />
          </div>
          <h1 className="text-2xl font-bold text-[#e2e8f0]">FISHGAME CMS</h1>
          <p className="text-slate-500 text-sm mt-1">Đăng nhập hệ thống quản trị</p>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-500/10 border border-red-500/30 rounded text-red-400 text-sm text-center">
            {error}
          </div>
        )}

        <form onSubmit={handleLogin} className="space-y-5">
          <div>
            <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">Tên đăng nhập</label>
            <div className="relative">
              <User className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
              <input 
                type="text" 
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full bg-[#0d1425] border border-slate-700 text-[#e2e8f0] rounded-lg pl-10 pr-4 py-2.5 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-colors"
                placeholder="Nhập tên đăng nhập"
                required
              />
            </div>
          </div>

          <div>
            <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">Mật khẩu</label>
            <div className="relative">
              <KeyRound className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
              <input 
                type="password" 
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full bg-[#0d1425] border border-slate-700 text-[#e2e8f0] rounded-lg pl-10 pr-4 py-2.5 focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-colors"
                placeholder="••••••••"
                required
              />
            </div>
          </div>

          <button 
            type="submit"
            className="w-full bg-blue-600 hover:bg-blue-500 text-white font-bold py-3 rounded-lg transition-colors mt-6"
          >
            Đăng nhập
          </button>
        </form>
      </div>
    </div>
  );
}
