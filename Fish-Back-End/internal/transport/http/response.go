package http

import (
	"errors"
	"net/http"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data, "error": nil})
}

func Fail(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatus, gin.H{"error": gin.H{"code": appErr.Code, "message": appErr.Message}})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
		"code":    apperror.ErrInternalServer.Code,
		"message": apperror.ErrInternalServer.Message,
	}})
}
