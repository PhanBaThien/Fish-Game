package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check related requests.
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler creates a new HealthHandler instance.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// HealthResponse represents the response body for the health endpoint.
type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
	GoVersion string `json:"go_version"`
}

// Check godoc
// @Summary      Health Check
// @Description  Returns the current health status of the API server.
// @Tags         System
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /api/v1/health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	uptime := time.Since(h.startTime).Round(time.Second).String()

	c.JSON(http.StatusOK, HealthResponse{
		Status:    "ok",
		Service:   "fish-game-backend-api",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    uptime,
		GoVersion: runtime.Version(),
	})
}
