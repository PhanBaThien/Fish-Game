import { useEffect, useRef, useCallback, useState } from 'react'
import { useAuthStore } from '../stores/authStore'
import { useWalletStore } from '../stores/walletStore'
import { useGameStore } from '../stores/gameStore'

// ── Message types (mirror backend ws/message.go) ──────────────────────────────

type WsOutMsg =
  | { type: 'join_room';   payload: { room_id: number } }
  | { type: 'shoot';       payload: { x: number; y: number; angle: number; bet_amount: number } }
  | { type: 'fish_killed'; payload: { fish_id: number; reward_multiplier: number } }
  | { type: 'leave_room';  payload: null }
  | { type: 'ping';        payload: null }

interface WsInMsg {
  type: string
  payload: Record<string, unknown>
}

interface ShootAckPayload  { shots_fired: number; total_spend: number; balance: number }
interface EarnAckPayload   { amount: number; balance: number; total_earn: number; fish_killed: number }
interface SessionStartedPayload { session_id: number }
interface SessionEndedPayload   { session: unknown; wallet: { balance: number } }

// ── Hook ──────────────────────────────────────────────────────────────────────

export type WsStatus = 'idle' | 'connecting' | 'connected' | 'disconnected'

interface WsErrorPayload { code: string; message: string }

interface UseGameSocketReturn {
  status:        WsStatus
  sessionId:     number | null
  lastError:     WsErrorPayload | null
  sendShoot:     (x: number, y: number, angle: number, betAmount: number) => void
  sendFishKilled:(fishId: number, rewardMultiplier: number) => void
}

export function useGameSocket(roomId: number | null): UseGameSocketReturn {
  const wsRef    = useRef<WebSocket | null>(null)
  const roomIdRef = useRef(roomId) // stable ref để dùng trong cleanup

  const [status,    setStatus]    = useState<WsStatus>('idle')
  const [sessionId, setSessionId] = useState<number | null>(null)
  const [lastError, setLastError] = useState<WsErrorPayload | null>(null)

  const accessToken = useAuthStore(s => s.accessToken)
  const { setBalance }       = useWalletStore()
  const { addCoins, addScore } = useGameStore()

  // ── send helper ─────────────────────────────────────────────────────────────
  const send = useCallback((msg: WsOutMsg) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg))
    }
  }, [])

  // ── Public send actions ──────────────────────────────────────────────────────
  const sendShoot = useCallback(
    (x: number, y: number, angle: number, betAmount: number) =>
      send({ type: 'shoot', payload: { x, y, angle, bet_amount: betAmount } }),
    [send],
  )

  const sendFishKilled = useCallback(
    (fishId: number, rewardMultiplier: number) =>
      send({ type: 'fish_killed', payload: { fish_id: fishId, reward_multiplier: rewardMultiplier } }),
    [send],
  )

  // ── Message handler ──────────────────────────────────────────────────────────
  const handleMessage = useCallback(
    (msg: WsInMsg) => {
      switch (msg.type) {
        case 'session_started': {
          const p = msg.payload as unknown as SessionStartedPayload
          setSessionId(p.session_id)
          break
        }
        case 'shoot_ack': {
          const p = msg.payload as unknown as ShootAckPayload
          setBalance(p.balance)  // cập nhật balance sau khi trừ tiền đạn
          break
        }

        case 'earn_ack': {
          const p = msg.payload as unknown as EarnAckPayload
          addCoins(p.amount)   // cộng vào coins "ván này"
          addScore(1)
          setBalance(p.balance) // cập nhật balance ước tính ngay lên Navbar
          break
        }
        case 'session_ended': {
          const p = msg.payload as unknown as SessionEndedPayload
          setBalance(p.wallet.balance) // balance chính xác từ DB sau EndSession
          setSessionId(null)
          break
        }
        case 'error': {
          const p = msg.payload as unknown as WsErrorPayload
          setLastError(p)
          // Tự xóa sau 3s
          setTimeout(() => setLastError(null), 3000)
          break
        }
        case 'pong':
          break
        default:
          break
      }
    },
    [addCoins, addScore, setBalance],
  )

  // ── Connect / Disconnect ─────────────────────────────────────────────────────
  useEffect(() => {
    if (!accessToken || !roomId) return
    roomIdRef.current = roomId

    // Xây URL WebSocket: dùng window.location để tự động đúng host
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const url = `${protocol}//${window.location.host}/api/v1/ws?token=${accessToken}`

    setStatus('connecting')
    const ws = new WebSocket(url)
    wsRef.current = ws

    ws.onopen = () => {
      setStatus('connected')
      ws.send(JSON.stringify({ type: 'join_room', payload: { room_id: roomId } }))
    }

    ws.onmessage = (event: MessageEvent<string>) => {
      try {
        const msg: WsInMsg = JSON.parse(event.data)
        handleMessage(msg)
      } catch (e) {
        console.error('[WS] parse error', e)
      }
    }

    ws.onerror = () => setStatus('disconnected')

    ws.onclose = () => setStatus('disconnected')

    // Cleanup: gửi leave_room trước khi đóng (component unmount / roomId thay đổi)
    return () => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: 'leave_room', payload: null }))
        // Cho một nhịp để server nhận leave_room trước khi đóng kết nối
        setTimeout(() => ws.close(1000, 'leaving room'), 150)
      } else {
        ws.close()
      }
      wsRef.current = null
    }
  }, [accessToken, roomId, handleMessage])

  return { status, sessionId, lastError, sendShoot, sendFishKilled }
}
