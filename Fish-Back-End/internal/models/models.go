package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	RoleID    int32     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Role struct {
	ID       int32  `json:"id"`
	RoleName string `json:"role_name"`
}

type Room struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	MaxPlayers  int32     `json:"max_players"`
	Description *string   `json:"description"`
	RTP         float64   `json:"rtp"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Fish struct {
	ID               int32     `json:"id"`
	Name             string    `json:"name"`
	Health           int32     `json:"health"`
	RewardMultiplier int32     `json:"reward_multiplier"`
	Speed            float64   `json:"speed"`
	AssetPath        string    `json:"asset_path"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Wallet struct {
	UserID    int64     `json:"user_id"`
	Balance   int64     `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	SessionID   *int64    `json:"session_id"` // nil = giao dịch ngoài session
	Amount      int64     `json:"amount"`     // dương = nhận, âm = tiêu
	Type        string    `json:"type"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type GameSession struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	RoomID     int64      `json:"room_id"`
	ShotsFired int32      `json:"shots_fired"`
	FishKilled int32      `json:"fish_killed"`
	TotalSpend int64      `json:"total_spend"`
	TotalEarn  int64      `json:"total_earn"`
	Status     string     `json:"status"` // active | finished
	StartedAt  time.Time  `json:"started_at"`
	EndedAt    *time.Time `json:"ended_at"` // nil khi đang active
}

type RefreshToken struct {
	ID        int64     `json:"-"`
	UserID    int64     `json:"-"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}
