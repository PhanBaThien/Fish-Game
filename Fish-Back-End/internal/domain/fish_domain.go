package domain

import (
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
)

// ─── REQUEST / RESPONSE cho CRUD Cấu hình Cá ──────────────────────────────────

type ListFishRequest struct {
	Role     string `form:"role"`  // common | mid | boss | "" (all)
	IsActive *bool  `form:"active"` // nil = all, true/false = filter
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}

type ListFishResponse struct {
	Items []models.Fish `json:"items"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

type CreateFishRequest struct {
	ID         string  `json:"id"          binding:"required,max=20"` // e.g. "F01"
	Name       string  `json:"name"        binding:"required"`
	Multiplier int     `json:"multiplier"  binding:"required,min=1"`
	BaseProb   float64 `json:"base_prob"   binding:"required,min=0,max=1"`
	Speed      string  `json:"speed"       binding:"required,oneof=fast medium slow very_slow"`
	Role       string  `json:"role"        binding:"required,oneof=common mid boss"`
	IsActive   bool    `json:"is_active"`
}

type UpdateFishRequest struct {
	Name       string  `json:"name"        binding:"required"`
	Multiplier int     `json:"multiplier"  binding:"required,min=1"`
	BaseProb   float64 `json:"base_prob"   binding:"required,min=0,max=1"`
	Speed      string  `json:"speed"       binding:"required,oneof=fast medium slow very_slow"`
	Role       string  `json:"role"        binding:"required,oneof=common mid boss"`
	IsActive   bool    `json:"is_active"`
}
