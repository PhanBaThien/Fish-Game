package models

import (
	"encoding/json"
	"time"
)

// ─── QUẢN LÝ TÀI KHOẢN ───────────────────────────────────────────────────────

// Admin represents an admin account stored in the DB.
type Admin struct {
	ID           string    `db:"id" json:"id"`                       // UUID (gen_random_uuid())
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`                   // super_admin | admin | moderator
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// Player represents a game player.
type Player struct {
	ID           string     `db:"id" json:"id"`                        // UUID (gen_random_uuid())
	Username     string     `db:"username" json:"username"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	GoldBalance  int64      `db:"gold_balance" json:"gold_balance"`    // BIGINT
	Status       string     `db:"status" json:"status"`                // active | banned | suspended
	WinRate      float64    `db:"win_rate" json:"win_rate"`            // NUMERIC(6,2)
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	LastLoginAt  *time.Time `db:"last_login_at" json:"last_login_at"`
}

// ─── QUẢN LÝ TÀI CHÍNH (VÍ & SỔ CÁI) ─────────────────────────────────────────

// Wallet represents a player's balance with optimistic locking.
// PK: player_id (FK → players.id)
type Wallet struct {
	PlayerID  string    `db:"player_id" json:"player_id"`          // UUID, FK → players.id
	Balance   int64     `db:"balance" json:"balance"`              // Đơn vị: xu (1 vàng = 100 xu)
	Version   int       `db:"version" json:"version"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Transaction represents a ledger entry for balance changes.
type Transaction struct {
	ID             string          `db:"id" json:"id"`                             // UUID (gen_random_uuid())
	PlayerID       string          `db:"player_id" json:"player_id"`               // UUID, FK → players.id
	Type           string          `db:"type" json:"type"`                         // shot | win | gift | deposit | withdraw | adjust
	Amount         int64           `db:"amount" json:"amount"`
	BalanceAfter   int64           `db:"balance_after" json:"balance_after"`
	RefShotID      *string         `db:"ref_shot_id" json:"ref_shot_id,omitempty"` // UUID, FK → shots.id
	RefKillID      *string         `db:"ref_kill_id" json:"ref_kill_id,omitempty"` // UUID, FK → fish_kills.id
	IdempotencyKey string          `db:"idempotency_key" json:"idempotency_key"`
	Metadata       json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	CreatedAt      time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at" json:"updated_at"`
}

// ─── CẤU HÌNH GAME (PHÒNG CHƠI & CÁ) ─────────────────────────────────────────

// Room represents a game room configuration.
type Room struct {
	ID         string    `db:"id" json:"id"`                        // UUID (gen_random_uuid())
	Name       string    `db:"name" json:"name"`
	Type       string    `db:"type" json:"type"`                    // beginner | advanced | expert | vip | boss
	BetAmount  int64     `db:"bet_amount" json:"bet_amount"`        // BIGINT
	MaxPlayers int       `db:"max_players" json:"max_players"`
	Status     string    `db:"status" json:"status"`                // waiting | playing | closed
	BaseRTP    float64   `db:"base_rtp" json:"base_rtp"`            // NUMERIC(5,2)
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// RoomSeat represents an active seat in a room.
// PK: (room_id, seat_index)
type RoomSeat struct {
	RoomID     string     `db:"room_id" json:"room_id"`              // UUID, FK → rooms.id
	SeatIndex  int        `db:"seat_index" json:"seat_index"`        // 0..3
	OccupiedBy *string    `db:"occupied_by" json:"occupied_by,omitempty"` // UUID, FK → players.id
	JoinedAt   *time.Time `db:"joined_at" json:"joined_at,omitempty"`
}

// Fish represents the configuration and stats of a 2D fish.
type Fish struct {
	ID         string    `db:"id" json:"id"`                        // VARCHAR(20) ví dụ: 'F01'
	Name       string    `db:"name" json:"name"`
	Multiplier int       `db:"multiplier" json:"multiplier"`
	BaseProb   float64   `db:"base_prob" json:"base_prob"`          // NUMERIC(6,3)
	Speed      string    `db:"speed" json:"speed"`                  // fast | medium | slow | very_slow
	Role       string    `db:"role" json:"role"`                    // common | mid | boss
	IsActive   bool      `db:"is_active" json:"is_active"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// ─── THEO DÕI GAMEPLAY & AUDIT ───────────────────────────────────────────────

// GameSession tracks the stats of a player in a room.
type GameSession struct {
	ID         string          `db:"id" json:"id"`                        // UUID (gen_random_uuid())
	RoomID     string          `db:"room_id" json:"room_id"`              // UUID, FK → rooms.id
	PlayerID   string          `db:"player_id" json:"player_id"`          // UUID, FK → players.id
	GoldStart  int64           `db:"gold_start" json:"gold_start"`
	GoldEnd    *int64          `db:"gold_end" json:"gold_end,omitempty"`
	FishCaught json.RawMessage `db:"fish_caught" json:"fish_caught"`      // JSONB
	StartedAt  time.Time       `db:"started_at" json:"started_at"`
	EndedAt    *time.Time      `db:"ended_at" json:"ended_at,omitempty"`
}

// Shot represents a bullet fired by a player.
type Shot struct {
	ID             string    `db:"id" json:"id"`                        // UUID (gen_random_uuid())
	SessionID      string    `db:"session_id" json:"session_id"`        // UUID, FK → game_sessions.id
	PlayerID       string    `db:"player_id" json:"player_id"`          // UUID, FK → players.id
	RoomID         string    `db:"room_id" json:"room_id"`              // UUID, FK → rooms.id
	BetAmount      int64     `db:"bet_amount" json:"bet_amount"`
	Angle          float64   `db:"angle" json:"angle"`
	ShotAt         time.Time `db:"shot_at" json:"shot_at"`
	IdempotencyKey string    `db:"idempotency_key" json:"idempotency_key"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// FishKill represents a successful hit that killed a fish.
type FishKill struct {
	ID        string    `db:"id" json:"id"`                        // UUID (gen_random_uuid())
	ShotID    string    `db:"shot_id" json:"shot_id"`              // UUID, FK → shots.id
	FishID    string    `db:"fish_id" json:"fish_id"`              // VARCHAR(20), FK → fish.id
	Payout    int64     `db:"payout" json:"payout"`
	RngSeed   int64     `db:"rng_seed" json:"rng_seed"`            // Dùng để audit thuật toán random
	ProbUsed  float64   `db:"prob_used" json:"prob_used"`
	KilledAt  time.Time `db:"killed_at" json:"killed_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// AuditLog tracks sensitive actions in the system.
type AuditLog struct {
	ID        string          `db:"id" json:"id"`                        // UUID (gen_random_uuid())
	ActorType string          `db:"actor_type" json:"actor_type"`        // admin | player | system
	ActorID   string          `db:"actor_id" json:"actor_id"`            // UUID, FK → admins.id hoặc players.id
	Action    string          `db:"action" json:"action"`
	Target    string          `db:"target" json:"target"`
	Payload   json.RawMessage `db:"payload" json:"payload"`
	IP        string          `db:"ip" json:"ip"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

// ─── CÀI ĐẶT HỆ THỐNG ────────────────────────────────────────────────────────

// Setting represents a global system configuration.
// PK: key (không có id riêng)
type Setting struct {
	Key       string          `db:"key" json:"key"`
	Value     json.RawMessage `db:"value" json:"value"`
	UpdatedBy string          `db:"updated_by" json:"updated_by"`        // UUID, FK → admins.id
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}
