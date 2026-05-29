package http

import (
	"strconv"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/transport/http/middleware"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/usecase"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
	"github.com/gin-gonic/gin"
)

type FishHandler struct {
	fishUsecase FishUsecase
	tokenMaker  utils.TokenMaker
}

type FishUsecase = usecase.FishUsecase

func NewFishHandler(u usecase.FishUsecase, m utils.TokenMaker) *FishHandler {
	return &FishHandler{fishUsecase: u, tokenMaker: m}
}

func (h *FishHandler) RegisterRoutes(router *gin.RouterGroup) {
	fish := router.Group("/fish")
	fish.Use(middleware.AuthMiddleware(h.tokenMaker))
	{
		fish.GET("", h.List)
		fish.GET("/:id", h.GetByID)

		admin := fish.Group("")
		admin.Use(middleware.RequireRoles(domain.RoleAdmin, domain.RoleSuperAdmin))
		{
			admin.POST("", h.Create)
			admin.PUT("/:id", h.Update)
			admin.DELETE("/:id", h.Delete)
		}
	}
}

func parseFishID(c *gin.Context) (int32, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	return int32(id), err
}

func (h *FishHandler) List(c *gin.Context) {
	fishes, err := h.fishUsecase.List(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, fishes)
}

func (h *FishHandler) GetByID(c *gin.Context) {
	id, err := parseFishID(c)
	if err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	fish, err := h.fishUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, fish)
}

func (h *FishHandler) Create(c *gin.Context) {
	var req domain.CreateFishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	fish, err := h.fishUsecase.Create(c.Request.Context(), &req)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, fish)
}

func (h *FishHandler) Update(c *gin.Context) {
	id, err := parseFishID(c)
	if err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	var req domain.UpdateFishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	fish, err := h.fishUsecase.Update(c.Request.Context(), id, &req)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, fish)
}

func (h *FishHandler) Delete(c *gin.Context) {
	id, err := parseFishID(c)
	if err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	if err := h.fishUsecase.Delete(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	Success(c, gin.H{"message": "xóa cá thành công"})
}
