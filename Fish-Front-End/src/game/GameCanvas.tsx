import { useEffect, useRef, useCallback } from 'react'
import { GameScene } from './scenes/GameScene'
import { useGameStore } from '../stores/gameStore'
import { useWalletStore } from '../stores/walletStore'
import { walletApi } from '../api/wallet'
import type { Fish, Room } from '../types'

interface GameCanvasProps {
  roomId: number
  room: Room
  fishList: Fish[]
}

export default function GameCanvas({ room, fishList }: GameCanvasProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const gameSceneRef = useRef<GameScene | null>(null)
  const { addCoins, addScore } = useGameStore()
  const { optimisticEarn } = useWalletStore()

  const handleFishKilled = useCallback(
    (rewardMultiplier: number) => {
      const earned = Math.round(rewardMultiplier * room.min_bet)
      // Cập nhật local state ngay
      addCoins(earned)
      addScore(1)
      // Cập nhật wallet bar ngay (optimistic)
      optimisticEarn(earned)
      // Sync lên server (fire-and-forget, lỗi im lặng)
      walletApi
        .earn(earned, `Bắn hạ cá ×${rewardMultiplier} tại phòng ${room.name}`)
        .then((w) => {
          // Đồng bộ balance chính xác từ server
          useWalletStore.getState().setBalance(w.balance)
        })
        .catch(() => {
          // Nếu lỗi, rollback optimistic
          useWalletStore.getState().optimisticSpend(earned)
        })
    },
    [room.min_bet, room.name, addCoins, addScore, optimisticEarn],
  )

  const handleScore = useCallback(
    (points: number) => addScore(points),
    [addScore],
  )

  useEffect(() => {
    if (!canvasRef.current) return

    const timeout = setTimeout(() => {
      if (!canvasRef.current) return
      gameSceneRef.current = new GameScene({
        canvas: canvasRef.current,
        fishList,
        onFishKilled: handleFishKilled,
        onScore: handleScore,
      })
    }, 50)

    return () => {
      clearTimeout(timeout)
      gameSceneRef.current?.dispose()
      gameSceneRef.current = null
    }
  }, [fishList, handleFishKilled, handleScore])

  return (
    <canvas
      ref={canvasRef}
      style={{
        width: '100%',
        height: '100%',
        display: 'block',
        cursor: 'none',
      }}
    />
  )
}
