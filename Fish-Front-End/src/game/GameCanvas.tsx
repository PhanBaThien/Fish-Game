import { useEffect, useRef, useCallback } from 'react'
import { GameScene } from './scenes/GameScene'
import type { Fish, Room } from '../types'

interface GameCanvasProps {
  room: Room
  fishList: Fish[]
  onHitFish: (fishId: number, instanceId: string) => void
  onShot: (x: number, y: number, angle: number) => boolean
  // Ref mà parent gán để gọi confirmFishDeath từ bên ngoài
  confirmDeathRef: { current: ((instanceId: string) => void) | null }
}

export default function GameCanvas({ room, fishList, onHitFish, onShot, confirmDeathRef }: GameCanvasProps) {
  const canvasRef    = useRef<HTMLCanvasElement>(null)
  const gameSceneRef = useRef<GameScene | null>(null)

  const handleHitFish = useCallback(
    (fishId: number, instanceId: string) => onHitFish(fishId, instanceId),
    [onHitFish],
  )

  const handleShot = useCallback(
    (x: number, y: number, angle: number): boolean => onShot(x, y, angle),
    [onShot],
  )

  useEffect(() => {
    if (!canvasRef.current) return

    const timeout = setTimeout(() => {
      if (!canvasRef.current) return
      const scene = new GameScene({
        canvas: canvasRef.current,
        fishList,
        roomRtp: room.rtp,
        onHitFish: handleHitFish,
        onShot: handleShot,
      })
      gameSceneRef.current = scene
      // Expose confirmFishDeath về parent qua ref
      confirmDeathRef.current = (instanceId) => scene.confirmFishDeath(instanceId)
    }, 50)

    return () => {
      clearTimeout(timeout)
      gameSceneRef.current?.dispose()
      gameSceneRef.current = null
      confirmDeathRef.current = null
    }
  }, [fishList, handleHitFish, handleShot, confirmDeathRef])

  return (
    <canvas
      ref={canvasRef}
      style={{ width: '100%', height: '100%', display: 'block', cursor: 'none' }}
    />
  )
}
