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

// RoomHandler xử lý các HTTP request liên quan đến quản lý phòng chơi
type RoomHandler struct {
	roomUsecase usecase.RoomUsecase
	tokenMaker  utils.TokenMaker
}

func NewRoomHandler(u usecase.RoomUsecase, m utils.TokenMaker) *RoomHandler {
	return &RoomHandler{roomUsecase: u, tokenMaker: m}
}

// RegisterRoutes gắn tất cả route /rooms vào router group (đều yêu cầu JWT)
func (h *RoomHandler) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/rooms")
	g.Use(middleware.AuthMiddleware(h.tokenMaker))
	{
		g.GET("", h.ListRooms)
		g.POST("", h.CreateRoom)
		g.GET("/:id", h.GetRoom)
		g.PUT("/:id", h.UpdateRoom)
		g.DELETE("/:id", h.DeleteRoom)
		g.PUT("/:id/rtp", h.UpdateRoomRTP)
	}
}

// ListRooms godoc
// @Summary  Lấy danh sách phòng chơi
// @Tags     rooms
// @Param    status query string false "waiting|playing|closed"
// @Param    type   query string false "beginner|advanced|expert|vip|boss"
// @Param    page   query int    false "Số trang"
// @Param    limit  query int    false "Số bản ghi mỗi trang"
func (h *RoomHandler) ListRooms(c *gin.Context) {
	var req domain.ListRoomsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.roomUsecase.ListRooms(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp, "error": nil})
}

// GetRoom godoc
// @Summary  Lấy thông tin một phòng chơi
// @Tags     rooms
// @Param    id path string true "Room UUID"
func (h *RoomHandler) GetRoom(c *gin.Context) {
	id := c.Param("id")
	room, err := h.roomUsecase.GetRoom(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": room, "error": nil})
}

// CreateRoom godoc
// @Summary  Tạo mới phòng chơi
// @Tags     rooms
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req domain.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.roomUsecase.CreateRoom(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": room, "error": nil})
}

// UpdateRoom godoc
// @Summary  Cập nhật phòng chơi
// @Tags     rooms
func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.roomUsecase.UpdateRoom(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": room, "error": nil})
}

// DeleteRoom godoc
// @Summary  Xóa phòng chơi
// @Tags     rooms
func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	err := h.roomUsecase.DeleteRoom(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Xóa phòng chơi thành công"}, "error": nil})
}

// UpdateRoomRTP godoc
// @Summary  Cập nhật RTP của phòng chơi
// @Tags     rooms
// @Param    id path string true "Room UUID"
func (h *RoomHandler) UpdateRoomRTP(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateRoomRTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.roomUsecase.UpdateRoomRTP(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Cập nhật RTP phòng thành công"}, "error": nil})
}
