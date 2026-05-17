package middleware

import (
	"errors"
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

func abortWithAppError(c *gin.Context, err *apperror.AppError) {
	c.AbortWithStatusJSON(err.HTTPStatus, gin.H{"error": gin.H{
		"code":    err.Code,
		"message": err.Error(),
	}})
}

func AuthMiddleware(tokenMaker utils.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			abortWithAppError(c, apperror.New("INVALID_TOKEN", apperror.ErrInvalidToken.HTTPStatus, errors.New("thiếu header Authorization")))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 || strings.ToLower(fields[0]) != authorizationTypeBearer {
			abortWithAppError(c, apperror.New("INVALID_TOKEN", apperror.ErrInvalidToken.HTTPStatus, errors.New("định dạng token không hợp lệ")))
			return
		}

		claims, err := tokenMaker.ExtractToken(fields[1])
		if err != nil {
			var appErr *apperror.AppError
			if !errors.As(err, &appErr) {
				appErr = apperror.ErrInvalidToken
			}
			abortWithAppError(c, appErr)
			return
		}

		c.Set(ctxUserIDKey, utils.ToInt64((*claims)["user_id"]))
		c.Set(ctxRoleKey, utils.ToInt64((*claims)["role_id"]))

		c.Next()
	}
}
