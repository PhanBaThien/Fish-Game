package middleware

import (
	"strings"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	ctxUserIDKey            = "user_id"
	ctxRoleKey              = "role_id"
)

func AuthMiddleware(tokenMaker utils.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(apperror.ErrInvalidToken.HTTPStatus, gin.H{"error": gin.H{
				"code":    apperror.ErrInvalidToken.Code,
				"message": "thiếu header Authorization",
			}})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 || strings.ToLower(fields[0]) != authorizationTypeBearer {
			c.AbortWithStatusJSON(apperror.ErrInvalidToken.HTTPStatus, gin.H{"error": gin.H{
				"code":    apperror.ErrInvalidToken.Code,
				"message": "định dạng token không hợp lệ",
			}})
			return
		}

		claims, err := tokenMaker.ExtractToken(fields[1])
		if err != nil {
			appErr, ok := err.(*apperror.AppError)
			if !ok {
				appErr = apperror.ErrInvalidToken
			}
			c.AbortWithStatusJSON(appErr.HTTPStatus, gin.H{"error": gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			}})
			return
		}

		c.Set(ctxUserIDKey, utils.ToInt64((*claims)["user_id"]))
		c.Set(ctxRoleKey, utils.ToInt64((*claims)["role_id"]))

		c.Next()
	}
}
