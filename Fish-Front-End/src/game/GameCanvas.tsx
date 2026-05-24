import { useEffect, useRef, useCallback } from 'react'
import { GameScene } from './scenes/GameScene'
import { useGameStore } from '../stores/gameStore'
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

  const handleFishKilled = useCallback(
    (rewardMultiplier: number) => {
      addCoins(Math.round(rewardMultiplier * room.min_bet))
      addScore(1)
    },
    [room.min_bet, addCoins, addScore],
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
