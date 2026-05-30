package ws

import "encoding/json"

// ── Message types Client → Server ────────────────────────────────────────────
const (
	MsgJoinRoom   = "join_room"   // vào phòng
	MsgShoot      = "shoot"       // bắn đạn
	MsgFishKilled = "fish_killed" // cá bị hạ (client báo)
	MsgLeaveRoom  = "leave_room"  // thoát phòng
	MsgPing       = "ping"
)

// ── Message types Server → Client ────────────────────────────────────────────
const (
	MsgSessionStarted = "session_started" // session tạo thành công
	MsgShootAck       = "shoot_ack"       // server nhận shot, trả stats
	MsgEarnAck        = "earn_ack"        // server xác nhận cá chết, trả balance
	MsgSessionEnded   = "session_ended"   // kết thúc ván, trả wallet
	MsgError          = "error"
	MsgPong           = "pong"
)

// ── Envelope ─────────────────────────────────────────────────────────────────

type InMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type OutMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// ── Client → Server payloads ─────────────────────────────────────────────────

type JoinRoomPayload struct {
	RoomID int64 `json:"room_id"`
}

type ShootPayload struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Angle     float64 `json:"angle"`
	BetAmount int64   `json:"bet_amount"` // tiền đặt cược cho viên đạn này
}

type FishKilledPayload struct {
	FishID           int32 `json:"fish_id"`
	RewardMultiplier int32 `json:"reward_multiplier"`
}

// ── Server → Client payloads ─────────────────────────────────────────────────

type SessionStartedPayload struct {
	SessionID int64 `json:"session_id"`
}

type ShootAckPayload struct {
	ShotsFired int32 `json:"shots_fired"`
	TotalSpend int64 `json:"total_spend"`
	Balance    int64 `json:"balance"` // balance ước tính sau khi trừ tiền đạn
}

type EarnAckPayload struct {
	Amount     int64 `json:"amount"`      // vàng kiếm được từ con cá này
	Balance    int64 `json:"balance"`     // balance ước tính (sẽ confirm lúc EndSession)
	TotalEarn  int64 `json:"total_earn"`  // tổng earn ván này
	FishKilled int32 `json:"fish_killed"` // tổng cá bắn hạ ván này
}

type SessionEndedPayload struct {
	Session any `json:"session"`
	Wallet  any `json:"wallet"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
