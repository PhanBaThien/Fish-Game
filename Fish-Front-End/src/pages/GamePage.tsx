import { useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { roomsApi } from '../api/rooms'
import { fishApi } from '../api/fish'
import { useGameStore } from '../stores/gameStore'
import GameCanvas from '../game/GameCanvas'

export default function GamePage() {
  const { roomId: roomIdStr } = useParams<{ roomId: string }>()
  const navigate = useNavigate()
  const roomId = Number(roomIdStr)

  const { coins, score, setCurrentRoom, resetGame } = useGameStore()

  const {
    data: room,
    isLoading: roomLoading,
    isError: roomError,
  } = useQuery({
    queryKey: ['room', roomId],
    queryFn: () => roomsApi.get(roomId),
    enabled: !!roomId && !isNaN(roomId),
  })

  const {
    data: fishList,
    isLoading: fishLoading,
    isError: fishError,
  } = useQuery({
    queryKey: ['fish'],
    queryFn: fishApi.list,
  })

  // Set current room in store
  useEffect(() => {
    if (room) {
      setCurrentRoom(room)
    }
  }, [room, setCurrentRoom])

  // Reset game state on mount
  useEffect(() => {
    resetGame()
    return () => {
      setCurrentRoom(null)
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const isLoading = roomLoading || fishLoading
  const isError = roomError || fishError

  if (isNaN(roomId)) {
    return (
      <div className="min-h-screen bg-slate-900 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-400 text-lg font-medium">Invalid room ID</p>
          <button onClick={() => navigate('/lobby')} className="mt-4 text-cyan-400 underline">
            Back to Lobby
          </button>
        </div>
      </div>
    )
  }

  if (isError) {
    return (
      <div className="min-h-screen bg-slate-900 flex items-center justify-center">
        <div className="text-center">
          <span className="text-5xl mb-4 block">🌊</span>
          <p className="text-red-400 text-lg font-medium mb-2">Failed to load game data</p>
          <button
            onClick={() => navigate('/lobby')}
            className="mt-4 px-6 py-2.5 rounded-xl bg-gradient-to-r from-cyan-600 to-teal-600 text-white font-medium"
          >
            Back to Lobby
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="fixed inset-0 bg-slate-900 overflow-hidden">
      {/* Game canvas fills entire screen */}
      {!isLoading && room && fishList && (
        <GameCanvas roomId={roomId} room={room} fishList={fishList} />
      )}

      {/* Loading overlay */}
      {isLoading && (
        <div className="absolute inset-0 flex items-center justify-center bg-slate-900">
          <div className="text-center">
            <div className="w-16 h-16 border-4 border-cyan-500/30 border-t-cyan-400 rounded-full animate-spin mx-auto mb-4" />
            <p className="text-white/60 text-lg">Loading game...</p>
          </div>
        </div>
      )}

      {/* HUD Overlay */}
      <div className="absolute inset-0 pointer-events-none">
        {/* Top bar */}
        <div className="absolute top-0 left-0 right-0 flex items-center justify-between px-4 py-3">
          {/* Room name */}
          <div className="flex items-center gap-2 px-4 py-2 rounded-xl bg-black/40 backdrop-blur border border-white/10">
            <span className="text-lg">🏠</span>
            <span className="text-white font-semibold text-sm">
              {room?.name ?? 'Loading...'}
            </span>
          </div>

          {/* Back button (pointer-events enabled) */}
          <button
            className="pointer-events-auto px-4 py-2 rounded-xl bg-black/40 backdrop-blur border border-white/10 text-white/70 hover:text-white hover:border-white/30 text-sm font-medium transition-all flex items-center gap-1.5"
            onClick={() => navigate('/lobby')}
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Lobby
          </button>
        </div>

        {/* Bottom HUD */}
        <div className="absolute bottom-0 left-0 right-0 flex items-end justify-between px-4 pb-4">
          {/* Coins */}
          <div className="flex flex-col items-center gap-1 px-5 py-3 rounded-2xl bg-black/50 backdrop-blur border border-yellow-500/20">
            <span className="text-2xl">🪙</span>
            <span className="text-yellow-400 font-extrabold text-2xl leading-none">
              {coins.toLocaleString()}
            </span>
            <span className="text-white/40 text-xs uppercase tracking-wider">Coins</span>
          </div>

          {/* Score */}
          <div className="flex flex-col items-center gap-1 px-5 py-3 rounded-2xl bg-black/50 backdrop-blur border border-cyan-500/20">
            <span className="text-2xl">🎯</span>
            <span className="text-cyan-400 font-extrabold text-2xl leading-none">
              {score}
            </span>
            <span className="text-white/40 text-xs uppercase tracking-wider">Fish</span>
          </div>

          {/* Min bet info */}
          <div className="flex flex-col items-center gap-1 px-5 py-3 rounded-2xl bg-black/50 backdrop-blur border border-teal-500/20">
            <span className="text-2xl">💰</span>
            <span className="text-teal-400 font-extrabold text-xl leading-none">
              {room?.min_bet?.toLocaleString() ?? '...'}
            </span>
            <span className="text-white/40 text-xs uppercase tracking-wider">Min Bet</span>
          </div>
        </div>

        {/* Controls hint */}
        <div className="absolute top-1/2 right-4 -translate-y-1/2 text-right">
          <div className="px-3 py-2 rounded-xl bg-black/30 backdrop-blur border border-white/5 text-white/25 text-xs space-y-1">
            <p>Click to shoot</p>
            <p>Scroll to zoom</p>
            <p>Drag to rotate</p>
          </div>
        </div>
      </div>
    </div>
  )
}
