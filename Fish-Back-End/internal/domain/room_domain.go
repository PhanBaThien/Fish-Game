package domain

import (
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
)

// ─── REQUEST / RESPONSE cho CRUD Phòng chơi ───────────────────────────────────

type ListRoomsRequest struct {
	Status string `form:"status"` // waiting | playing | closed | "" (all)
	Type   string `form:"type"`   // beginner | advanced | expert | vip | boss
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type ListRoomsResponse struct {
	Items []models.Room `json:"items"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

type CreateRoomRequest struct {
	Name       string  `json:"name"        binding:"required"`
	Type       string  `json:"type"        binding:"required,oneof=beginner advanced expert vip boss"`
	BetAmount  int64   `json:"bet_amount"  binding:"required,min=1"`
	MaxPlayers int     `json:"max_players" binding:"required,min=2,max=4"`
	BaseRTP    float64 `json:"base_rtp"    binding:"required,min=50,max=100"`
}

type UpdateRoomRequest struct {
	Name       string  `json:"name"        binding:"required"`
	Type       string  `json:"type"        binding:"required,oneof=beginner advanced expert vip boss"`
	BetAmount  int64   `json:"bet_amount"  binding:"required,min=1"`
	MaxPlayers int     `json:"max_players" binding:"required,min=2,max=4"`
	Status     string  `json:"status"      binding:"required,oneof=waiting playing closed"`
	BaseRTP    float64 `json:"base_rtp"    binding:"required,min=50,max=100"`
}

type UpdateRoomRTPRequest struct {
	BaseRTP float64 `json:"base_rtp" binding:"required,min=50,max=100"`
}
