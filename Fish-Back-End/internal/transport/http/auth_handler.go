package http

import (
	"net/http"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/transport/http/middleware"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/usecase"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	tokenMaker  utils.TokenMaker
}

func NewAuthHandler(u usecase.AuthUsecase, m utils.TokenMaker) *AuthHandler {
	return &AuthHandler{
		authUsecase: u,
		tokenMaker:  m,
	}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", h.Login)
		authRoutes.POST("/logout", h.Logout)
		authRoutes.POST("/register", h.Register)

		protected := authRoutes.Group("/")
		protected.Use(middleware.AuthMiddleware(h.tokenMaker))
		{
			protected.GET("/me", h.Me)
		}
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}

	resp, err := h.authUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Đăng xuất thành công"})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}

	result, err := h.authUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c, result)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	data, err := h.authUsecase.Me(c.Request.Context(), userID.(int64))
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c, data)
}
