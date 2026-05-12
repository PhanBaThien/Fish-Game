package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/fish-game-backend/internal/models"
	"github.com/yourname/fish-game-backend/internal/services"
)

// RoomHandler holds the RoomService dependency.
type RoomHandler struct {
	svc services.RoomService
}

// NewRoomHandler creates a RoomHandler wired to the given service.
func NewRoomHandler(svc services.RoomService) *RoomHandler {
	return &RoomHandler{svc: svc}
}

// ListRooms returns all game rooms with their current state.
func (h *RoomHandler) ListRooms(c *gin.Context) {
	rooms, err := h.svc.ListRooms(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rooms, "total": len(rooms)})
}

// CreateRoom creates a new game room.
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req models.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	room, err := h.svc.CreateRoom(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, room)
}

// GetRoom returns details for a single room.
func (h *RoomHandler) GetRoom(c *gin.Context) {
	id := c.Param("id")
	room, err := h.svc.GetRoomByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, room)
}

// UpdateRoom modifies a room's configuration (e.g., RTP).
func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateRoom(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Room " + id + " updated"})
}

// CloseRoom closes an active game room.
func (h *RoomHandler) CloseRoom(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.CloseRoom(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Room " + id + " closed"})
}
