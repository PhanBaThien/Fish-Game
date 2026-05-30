import { useEffect, useRef, useCallback } from 'react'
import { GameScene } from './scenes/GameScene'
import type { Fish, Room } from '../types'

interface GameCanvasProps {
  room: Room
  fishList: Fish[]
  onFishKilled: (fishId: number, rewardMultiplier: number) => void
  onShot: (x: number, y: number, angle: number) => boolean
}

export default function GameCanvas({ room: _room, fishList, onFishKilled, onShot }: GameCanvasProps) {
  const canvasRef    = useRef<HTMLCanvasElement>(null)
  const gameSceneRef = useRef<GameScene | null>(null)

  const handleFishKilled = useCallback(
    (fishId: number, rewardMultiplier: number) => onFishKilled(fishId, rewardMultiplier),
    [onFishKilled],
  )

  const handleShot = useCallback(
    (x: number, y: number, angle: number): boolean => onShot(x, y, angle),
    [onShot],
  )

  useEffect(() => {
    if (!canvasRef.current) return

    const timeout = setTimeout(() => {
      if (!canvasRef.current) return
      gameSceneRef.current = new GameScene({
        canvas: canvasRef.current,
        fishList,
        onFishKilled: handleFishKilled,
        onShot: handleShot,
      })
    }, 50)

    return () => {
      clearTimeout(timeout)
      gameSceneRef.current?.dispose()
      gameSceneRef.current = null
    }
  }, [fishList, handleFishKilled, handleShot])

  return (
    <canvas
      ref={canvasRef}
      style={{ width: '100%', height: '100%', display: 'block', cursor: 'none' }}
    />
  )
}
