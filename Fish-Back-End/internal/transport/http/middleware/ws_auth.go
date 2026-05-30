package middleware

import (
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
	"github.com/gin-gonic/gin"
)

// WSAuthMiddleware xác thực JWT cho WebSocket.
// Browser không thể set custom header khi connect WS,
// nên token được truyền qua query param ?token=<access_token>.
func WSAuthMiddleware(tokenMaker utils.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			abortWithAppError(c, apperror.ErrInvalidToken)
			return
		}

		claims, err := tokenMaker.VerifyAccessToken(tokenStr)
		if err != nil {
			abortWithAppError(c, apperror.ErrInvalidToken)
			return
		}

		c.Set(ctxUserIDKey, utils.ToInt64((*claims)["user_id"]))
		c.Next()
	}
}
