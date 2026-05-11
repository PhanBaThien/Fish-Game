package models

import "time"

// ─── Admin ───────────────────────────────────────────────────────────────────

// Admin represents an admin account stored in the DB.
type Admin struct {
	ID           string `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	Role         string `db:"role"` // admin | superadmin
}

// ─── Player ──────────────────────────────────────────────────────────────────

// Player represents a game player entity.
type Player struct {
	ID          string    `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	Email       string    `json:"email,omitempty" db:"email"`
	GoldBalance int64     `json:"gold_balance" db:"gold_balance"`
	Status      string    `json:"status" db:"status"` // active | banned | suspended
	WinRate     float64   `json:"win_rate" db:"win_rate"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	LastLoginAt time.Time `json:"last_login_at" db:"last_login_at"`
}

// UpdatePlayerRequest is used to patch a player record.
type UpdatePlayerRequest struct {
	GoldBalance *int64  `json:"gold_balance,omitempty"`
	Status      *string `json:"status,omitempty"`
	WinRate     *float64 `json:"win_rate,omitempty"`
}

// ─── Fish ────────────────────────────────────────────────────────────────────

// Fish represents a fish entity in the game.
type Fish struct {
	ID         string  `json:"id" db:"id"`
	Name       string  `json:"name" db:"name"`
	Multiplier int     `json:"multiplier" db:"multiplier"` // Gold reward multiplier
	BaseProb   float64 `json:"base_prob" db:"base_prob"`   // Base catch probability (%)
	Speed      string  `json:"speed" db:"speed"`           // fast | medium | slow | very_slow
	Role       string  `json:"role" db:"role"`             // common | mid | boss
	IsActive   bool    `json:"is_active" db:"is_active"`
}

// CreateFishRequest is the payload for creating a new fish.
type CreateFishRequest struct {
	Name       string  `json:"name" binding:"required"`
	Multiplier int     `json:"multiplier" binding:"required,min=1"`
	BaseProb   float64 `json:"base_prob" binding:"required,min=0,max=100"`
	Speed      string  `json:"speed" binding:"required,oneof=fast medium slow very_slow"`
	Role       string  `json:"role" binding:"required,oneof=common mid boss"`
}

// UpdateFishRequest allows partial updates to fish config.
type UpdateFishRequest struct {
	Name       *string  `json:"name,omitempty"`
	Multiplier *int     `json:"multiplier,omitempty"`
	BaseProb   *float64 `json:"base_prob,omitempty"`
	Speed      *string  `json:"speed,omitempty"`
	IsActive   *bool    `json:"is_active,omitempty"`
}

// ─── Room ────────────────────────────────────────────────────────────────────

// Room represents a game room entity.
type Room struct {
	ID         string  `json:"id" db:"id"`
	Name       string  `json:"name" db:"name"`
	Type       string  `json:"type" db:"type"` // beginner | advanced | expert | vip | boss
	BetAmount  int64   `json:"bet_amount" db:"bet_amount"`
	Players    int     `json:"players" db:"players"`
	MaxPlayers int     `json:"max_players" db:"max_players"`
	Status     string  `json:"status" db:"status"` // waiting | playing | closed
	BaseRTP    float64 `json:"base_rtp" db:"base_rtp"` // Return to Player %
}

// CreateRoomRequest is the payload for creating a new room.
type CreateRoomRequest struct {
	Name       string  `json:"name" binding:"required"`
	Type       string  `json:"type" binding:"required,oneof=beginner advanced expert vip boss"`
	BetAmount  int64   `json:"bet_amount" binding:"required,min=1"`
	MaxPlayers int     `json:"max_players" binding:"required,min=1,max=8"`
	BaseRTP    float64 `json:"base_rtp" binding:"required,min=50,max=100"`
}

// UpdateRoomRequest allows partial updates (e.g. RTP adjustment).
type UpdateRoomRequest struct {
	Name      *string  `json:"name,omitempty"`
	BaseRTP   *float64 `json:"base_rtp,omitempty"`
	Status    *string  `json:"status,omitempty"`
	BetAmount *int64   `json:"bet_amount,omitempty"`
}
