package middleware

import (
	"net/http"
	"strings"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	ctxAdminIDKey           = "admin_id"
	ctxRoleKey              = "role"
)

// AuthMiddleware kiểm tra Bearer Token và lưu thông tin vào context
func AuthMiddleware(tokenMaker utils.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy Header Authorization
		authHeader := c.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Thiếu header Authorization"})
			return
		}

		// 2. Phân tách chuỗi "Bearer <token>"
		fields := strings.Fields(authHeader)
		if len(fields) < 2 || strings.ToLower(fields[0]) != authorizationTypeBearer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Định dạng token không hợp lệ"})
			return
		}

		tokenString := fields[1]

		claims, err := tokenMaker.ExtractToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ hoặc đã hết hạn"})
			return
		}

		adminID, _ := (*claims)["admin_id"].(string)
		role, _ := (*claims)["role"].(string)

		c.Set(ctxAdminIDKey, adminID)
		c.Set(ctxRoleKey, role)

		c.Next()
	}
}
