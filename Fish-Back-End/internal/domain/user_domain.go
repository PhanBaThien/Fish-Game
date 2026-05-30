package domain

import "github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken           string      `json:"access_token"`
	AccessTokenExpiresAt  int64       `json:"-"` // ẩn khỏi response
	RefreshToken          string      `json:"-"` // trả qua HttpOnly cookie, không expose trong JSON
	RefreshTokenExpiresAt int64       `json:"-"` // ẩn khỏi response
	User                  models.User `json:"-"` // ẩn khỏi response
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
}

type RegisterResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	RoleID   int32  `json:"role_id"`
}

// RefreshTokenRequest không cần nữa — token đọc từ cookie

type RefreshTokenResponse struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  int64  `json:"access_token_expires_at"`
	RefreshToken          string `json:"-"` // trả qua HttpOnly cookie
	RefreshTokenExpiresAt int64  `json:"-"`
}
