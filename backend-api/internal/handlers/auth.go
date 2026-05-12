package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/fish-game-backend/internal/services"
)

// LoginRequest represents the login payload.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateAdminRequest is the payload for creating a new admin account.
type CreateAdminRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`  // "admin" | "superadmin" — defaults to "admin"
}

// TokenResponse is returned after successful login.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// AuthHandler holds the AuthService dependency.
type AuthHandler struct {
	svc services.AuthService
}

// NewAuthHandler creates an AuthHandler wired to the given service.
func NewAuthHandler(svc services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login authenticates an admin and returns a JWT token pair.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pair, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		// "invalid credentials" is always 401 — do not leak whether user exists
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
	})
}

// RefreshToken issues a new access token given a valid refresh token.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// TODO: implement full refresh token rotation
	c.JSON(http.StatusOK, gin.H{"message": "Token refresh endpoint — coming soon"})
}

// CreateAdmin creates a new admin account. Requires a valid JWT (any admin can call this).
func (h *AuthHandler) CreateAdmin(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := h.svc.CreateAdmin(c.Request.Context(), req.Username, req.Password, req.Role)
	if err != nil {
		// Duplicate username → 409 Conflict
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Admin created successfully",
		"data":    admin,
	})
}

// GetDashboardStats returns aggregated statistics for the dashboard.
// TODO: inject StatsService and aggregate from DB + Redis
func GetDashboardStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"active_players":    1245,
		"revenue_today":     23400.00,
		"new_installs":      342,
		"fish_eliminated":   1200000,
		"active_rooms":      5,
		"server_uptime_pct": 99.8,
		"chart_data": []gin.H{
			{"name": "T2", "players": 4000, "revenue": 2400},
			{"name": "T3", "players": 3000, "revenue": 1398},
			{"name": "T4", "players": 2000, "revenue": 9800},
			{"name": "T5", "players": 2780, "revenue": 3908},
			{"name": "T6", "players": 1890, "revenue": 4800},
			{"name": "T7", "players": 2390, "revenue": 3800},
			{"name": "CN", "players": 3490, "revenue": 4300},
		},
	})
}
