package http

import (
	"errors"
	"net/http"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/middleware"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/usecase"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
	"github.com/gin-gonic/gin"
)

// PlayerHandler xử lý các HTTP request liên quan đến quản lý người chơi
type PlayerHandler struct {
	playerUsecase usecase.PlayerUsecase
	tokenMaker    utils.TokenMaker
}

func NewPlayerHandler(u usecase.PlayerUsecase, m utils.TokenMaker) *PlayerHandler {
	return &PlayerHandler{playerUsecase: u, tokenMaker: m}
}

// RegisterRoutes gắn tất cả route /players vào router group (đều yêu cầu JWT)
func (h *PlayerHandler) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/players")
	g.Use(middleware.AuthMiddleware(h.tokenMaker))
	{
		g.GET("", h.ListPlayers)
		g.POST("", h.CreatePlayer)
		g.GET("/:id", h.GetPlayer)
		g.PUT("/:id", h.UpdatePlayer)
		g.DELETE("/:id", h.DeletePlayer)
		g.POST("/:id/ban", h.BanPlayer)
		g.POST("/:id/gift", h.GiftPlayer)
		g.PUT("/:id/rtp", h.UpdatePlayerRTP)
	}
}

// ListPlayers godoc
// @Summary  Lấy danh sách người chơi
// @Tags     players
// @Param    status query string false "Lọc theo status: active|banned|suspended"
// @Param    page   query int    false "Số trang (mặc định 1)"
// @Param    limit  query int    false "Số bản ghi mỗi trang (mặc định 20, tối đa 100)"
// @Success  200 {object} domain.ListPlayersResponse
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	var req domain.ListPlayersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.playerUsecase.ListPlayers(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp, "error": nil})
}

// GetPlayer godoc
// @Summary  Lấy thông tin một người chơi
// @Tags     players
// @Param    id path string true "Player UUID"
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	id := c.Param("id")
	player, err := h.playerUsecase.GetPlayer(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrPlayerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": player, "error": nil})
}

// CreatePlayer godoc
// @Summary  Tạo mới người chơi
// @Tags     players
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var req domain.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := h.playerUsecase.CreatePlayer(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrPlayerUsernameTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": player, "error": nil})
}

// UpdatePlayer godoc
// @Summary  Cập nhật thông tin người chơi
// @Tags     players
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := h.playerUsecase.UpdatePlayer(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrPlayerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": player, "error": nil})
}

// DeletePlayer godoc
// @Summary  Xóa người chơi
// @Tags     players
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	id := c.Param("id")
	err := h.playerUsecase.DeletePlayer(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrPlayerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Xóa người chơi thành công"}, "error": nil})
}

// BanPlayer godoc
// @Summary  Ban người chơi
// @Tags     players
func (h *PlayerHandler) BanPlayer(c *gin.Context) {
	id := c.Param("id")
	var req domain.BanPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.playerUsecase.BanPlayer(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrPlayerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Ban người chơi thành công"}, "error": nil})
}

// GiftPlayer godoc
// @Summary  Tặng vàng cho người chơi
// @Tags     players
func (h *PlayerHandler) GiftPlayer(c *gin.Context) {
	id := c.Param("id")
	actorID, _ := c.Get("admin_id")
	actorIDStr, _ := actorID.(string)

	var req domain.GiftPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.playerUsecase.GiftPlayer(c.Request.Context(), id, actorIDStr, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrPlayerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Tặng vàng thành công"}, "error": nil})
}

// UpdatePlayerRTP godoc
// @Summary  Cập nhật RTP cá nhân của người chơi
// @Tags     players
func (h *PlayerHandler) UpdatePlayerRTP(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdatePlayerRTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.playerUsecase.UpdatePlayerRTP(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrPlayerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Cập nhật RTP thành công"}, "error": nil})
}
