package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/fish-game-backend/internal/models"
	"github.com/yourname/fish-game-backend/internal/services"
)

// FishHandler holds the FishService dependency.
type FishHandler struct {
	svc services.FishService
}

// NewFishHandler creates a FishHandler wired to the given service.
func NewFishHandler(svc services.FishService) *FishHandler {
	return &FishHandler{svc: svc}
}

// ListFish returns all fish configurations.
func (h *FishHandler) ListFish(c *gin.Context) {
	list, err := h.svc.ListFish(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "total": len(list)})
}

// CreateFish adds a new fish type.
func (h *FishHandler) CreateFish(c *gin.Context) {
	var req models.CreateFishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fish, err := h.svc.CreateFish(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, fish)
}

// UpdateFish modifies an existing fish configuration.
func (h *FishHandler) UpdateFish(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateFishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateFish(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Fish " + id + " updated"})
}

// DeleteFish removes a fish type from the configuration.
func (h *FishHandler) DeleteFish(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteFish(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Fish " + id + " deleted"})
}
