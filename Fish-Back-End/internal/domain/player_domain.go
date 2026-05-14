package domain

import (
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
)

// ─── REQUEST / RESPONSE cho CRUD Người chơi ───────────────────────────────────

// ListPlayersRequest chứa các tham số lọc & phân trang
type ListPlayersRequest struct {
	Status string `form:"status"` // active | banned | suspended | "" (all)
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

// ListPlayersResponse trả về danh sách player kèm tổng số bản ghi
type ListPlayersResponse struct {
	Items []models.Player `json:"items"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

// CreatePlayerRequest dữ liệu tạo mới người chơi
type CreatePlayerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdatePlayerRequest dữ liệu cập nhật thông tin người chơi
type UpdatePlayerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email"    binding:"required,email"`
	Status   string `json:"status"   binding:"required,oneof=active banned suspended"`
}

// BanPlayerRequest yêu cầu ban người chơi
type BanPlayerRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// GiftPlayerRequest yêu cầu tặng vàng cho người chơi
type GiftPlayerRequest struct {
	Amount int64  `json:"amount" binding:"required,min=1"`
	Note   string `json:"note"`
}

// UpdatePlayerRTPRequest yêu cầu đặt RTP cá nhân cho người chơi
type UpdatePlayerRTPRequest struct {
	WinRate float64 `json:"win_rate" binding:"required,min=0,max=100"`
}
