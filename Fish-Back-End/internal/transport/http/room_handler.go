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

type RoomHandler struct {
	roomUsecase RoomUsecase
	tokenMaker  utils.TokenMaker
}

type RoomUsecase = usecase.RoomUsecase

func NewRoomHandler(u usecase.RoomUsecase, m utils.TokenMaker) *RoomHandler {
	return &RoomHandler{roomUsecase: u, tokenMaker: m}
}

func (h *RoomHandler) RegisterRoutes(router *gin.RouterGroup) {
	rooms := router.Group("/rooms")
	rooms.Use(middleware.AuthMiddleware(h.tokenMaker))
	{
		rooms.GET("", h.List)
		rooms.GET("/:id", h.GetByID)

		admin := rooms.Group("")
		admin.Use(middleware.RequireRoles(domain.RoleAdmin, domain.RoleSuperAdmin))
		{
			admin.POST("", h.Create)
			admin.PUT("/:id", h.Update)
			admin.DELETE("/:id", h.Delete)
		}
	}
}

func (h *RoomHandler) List(c *gin.Context) {
	rooms, err := h.roomUsecase.List(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, rooms)
}

func (h *RoomHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	room, err := h.roomUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, room)
}

func (h *RoomHandler) Create(c *gin.Context) {
	var req domain.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	room, err := h.roomUsecase.Create(c.Request.Context(), &req)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, room)
}

func (h *RoomHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	var req domain.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	room, err := h.roomUsecase.Update(c.Request.Context(), id, &req)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, room)
}

func (h *RoomHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	if err := h.roomUsecase.Delete(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	Success(c, gin.H{"message": "xóa phòng thành công"})
}
