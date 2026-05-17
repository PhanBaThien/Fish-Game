package middleware

import (
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
	"github.com/gin-gonic/gin"
)

func RequireRoles(allowedRoles ...int32) gin.HandlerFunc {
	allowed := make(map[int32]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		raw, exists := c.Get(ctxRoleKey)
		if !exists {
			abortWithAppError(c, apperror.ErrInvalidToken)
			return
		}

		roleID := int32(utils.ToInt64(raw))
		if _, ok := allowed[roleID]; !ok {
			abortWithAppError(c, apperror.ErrForbidden)
			return
		}

		c.Next()
	}
}
