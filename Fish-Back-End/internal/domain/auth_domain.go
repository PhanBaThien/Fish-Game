package domain

import (
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
)

// LoginRequest là dữ liệu admin gửi lên
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse là dữ liệu trả về cho Frontend
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expiresAt"`
	Admin     models.Admin `json:"admin"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type RegisterResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
