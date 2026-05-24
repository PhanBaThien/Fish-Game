package http

import (
	"net/http"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/transport/http/middleware"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/usecase"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
	"github.com/gin-gonic/gin"
)

const refreshTokenCookie = "refresh_token"
const refreshTokenCookiePath = "/api/v1/auth"

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	tokenMaker  utils.TokenMaker
}

func NewAuthHandler(u usecase.AuthUsecase, m utils.TokenMaker) *AuthHandler {
	return &AuthHandler{authUsecase: u, tokenMaker: m}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
		auth.POST("/refresh", h.Refresh)

		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware(h.tokenMaker))
		{
			protected.GET("/me", h.Me)
			protected.POST("/logout", h.Logout)
		}
	}
}

func setRefreshCookie(c *gin.Context, token string, expiresAt int64) {
	maxAge := int(time.Until(time.Unix(expiresAt, 0)).Seconds())
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshTokenCookie, token, maxAge, refreshTokenCookiePath, "", false, true)
}

func clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshTokenCookie, "", -1, refreshTokenCookiePath, "", false, true)
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

	setRefreshCookie(c, resp.RefreshToken, resp.RefreshTokenExpiresAt)
	Success(c, resp) // RefreshToken có json:"-" nên không xuất hiện trong response
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

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err != nil {
		Fail(c, apperror.ErrInvalidToken)
		return
	}

	resp, err := h.authUsecase.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		clearRefreshCookie(c)
		Fail(c, err)
		return
	}

	setRefreshCookie(c, resp.RefreshToken, resp.RefreshTokenExpiresAt)
	Success(c, resp)
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

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err == nil {
		_ = h.authUsecase.Logout(c.Request.Context(), refreshToken)
	}
	clearRefreshCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "đăng xuất thành công"})
}
