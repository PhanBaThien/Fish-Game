package http

import (
	"net/http"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/middleware"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/usecase"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthHandler xử lý các HTTP request liên quan đến xác thực
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	tokenMaker  utils.TokenMaker
}

// NewAuthHandler khởi tạo handler và tiêm (inject) AuthUsecase vào
func NewAuthHandler(u usecase.AuthUsecase, m utils.TokenMaker) *AuthHandler {
	return &AuthHandler{
		authUsecase: u,
		tokenMaker:  m,
	}
}

// RegisterRoutes nhóm và gắn các endpoint của Auth vào Gin Router
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

// Login nhận request đăng nhập, gọi Usecase xử lý và trả về Token
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu yêu cầu không hợp lệ"})
		return
	}

	resp, err := h.authUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  resp,
		"error": nil,
	})
}

// Logout xử lý API đăng xuất hệ thống
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Đăng xuất thành công"})
}

// Register xử lý API đăng ký tài khoản admin mới
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu yêu cầu không hợp lệ"})
		return
	}

	result, err := h.authUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"error": nil,
	})
}

// Me trả về thông tin cơ bản của admin đang đăng nhập dựa trên context
func (h *AuthHandler) Me(c *gin.Context) {
	adminID, _ := c.Get("admin_id")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"admin_id": adminID,
			"role":     role,
		},
		"error": nil,
	})
}
