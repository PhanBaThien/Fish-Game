package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/fish-game-backend/internal/models"
	"github.com/yourname/fish-game-backend/internal/services"
)

// PlayerHandler holds the PlayerService dependency.
type PlayerHandler struct {
	svc services.PlayerService
}

// NewPlayerHandler creates a PlayerHandler wired to the given service.
func NewPlayerHandler(svc services.PlayerService) *PlayerHandler {
	return &PlayerHandler{svc: svc}
}

// ListPlayers returns a paginated list of all players.
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	players, err := h.svc.ListPlayers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  players,
		"total": len(players),
	})
}

// GetPlayer returns a single player by ID.
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	id := c.Param("id")
	player, err := h.svc.GetPlayerByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, player)
}

// UpdatePlayer updates a player's profile or gold balance.
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdatePlayer(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Player " + id + " updated successfully"})
}

// BanPlayer bans a player by ID.
func (h *PlayerHandler) BanPlayer(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.BanPlayer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Player " + id + " banned successfully"})
}
