package domain

import (
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	ExpiresAt int64       `json:"expiresAt"`
	User      models.User `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type RegisterResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	RoleID   int32  `json:"role_id"`
}