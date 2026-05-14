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

// FishHandler xử lý các HTTP request liên quan đến cấu hình cá
type FishHandler struct {
	fishUsecase usecase.FishUsecase
	tokenMaker  utils.TokenMaker
}

func NewFishHandler(u usecase.FishUsecase, m utils.TokenMaker) *FishHandler {
	return &FishHandler{fishUsecase: u, tokenMaker: m}
}

// RegisterRoutes gắn tất cả route /fishes vào router group (đều yêu cầu JWT)
func (h *FishHandler) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/fishes")
	g.Use(middleware.AuthMiddleware(h.tokenMaker))
	{
		g.GET("", h.ListFish)
		g.POST("", h.CreateFish)
		g.GET("/:id", h.GetFish)
		g.PUT("/:id", h.UpdateFish)
		g.DELETE("/:id", h.DeleteFish)
	}
}

// ListFish godoc
// @Summary  Lấy danh sách cấu hình cá
// @Tags     fishes
// @Param    role   query string false "Lọc theo role: common|mid|boss"
// @Param    active query bool   false "Lọc theo trạng thái active"
// @Param    page   query int    false "Số trang"
// @Param    limit  query int    false "Số bản ghi mỗi trang"
func (h *FishHandler) ListFish(c *gin.Context) {
	var req domain.ListFishRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.fishUsecase.ListFish(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp, "error": nil})
}

// GetFish godoc
// @Summary  Lấy thông tin một loại cá
// @Tags     fishes
// @Param    id path string true "Fish ID (e.g. F01)"
func (h *FishHandler) GetFish(c *gin.Context) {
	id := c.Param("id")
	fish, err := h.fishUsecase.GetFish(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrFishNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": fish, "error": nil})
}

// CreateFish godoc
// @Summary  Tạo mới cấu hình cá
// @Tags     fishes
func (h *FishHandler) CreateFish(c *gin.Context) {
	var req domain.CreateFishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fish, err := h.fishUsecase.CreateFish(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrFishIDExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": fish, "error": nil})
}

// UpdateFish godoc
// @Summary  Cập nhật cấu hình cá
// @Tags     fishes
func (h *FishHandler) UpdateFish(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateFishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fish, err := h.fishUsecase.UpdateFish(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrFishNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": fish, "error": nil})
}

// DeleteFish godoc
// @Summary  Xóa cấu hình cá
// @Tags     fishes
func (h *FishHandler) DeleteFish(c *gin.Context) {
	id := c.Param("id")
	err := h.fishUsecase.DeleteFish(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrFishNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Xóa cấu hình cá thành công"}, "error": nil})
}
