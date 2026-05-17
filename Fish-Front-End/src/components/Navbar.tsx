import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'

export default function Navbar() {
  const navigate = useNavigate()
  const { user, logout } = useAuthStore()

  const handleLogout = () => {
    logout()
    navigate('/login', { replace: true })
  }

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 flex items-center justify-between px-6 py-3 bg-slate-900/80 backdrop-blur border-b border-white/10">
      <div className="flex items-center gap-3">
        <span className="text-2xl">🐟</span>
        <span className="text-xl font-bold bg-gradient-to-r from-cyan-400 to-teal-400 bg-clip-text text-transparent">
          Fish Game
        </span>
      </div>

      <div className="flex items-center gap-4">
        {user && (
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-full bg-gradient-to-br from-cyan-500 to-teal-600 flex items-center justify-center text-white font-bold text-sm">
              {user.username.charAt(0).toUpperCase()}
            </div>
            <span className="text-white/80 text-sm font-medium">{user.username}</span>
          </div>
        )}

        <button
          onClick={handleLogout}
          className="px-4 py-1.5 rounded-lg text-sm font-medium text-white/70 hover:text-white border border-white/20 hover:border-white/40 hover:bg-white/5 transition-all duration-200"
        >
          Logout
        </button>
      </div>
    </nav>
  )
}
