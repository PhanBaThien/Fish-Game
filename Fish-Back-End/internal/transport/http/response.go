package http

import (
	"errors"
	"log"
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
		c.JSON(appErr.HTTPStatus, gin.H{"error": gin.H{
			"code":    appErr.Code,
			"message": appErr.Error(),
		}})
		return
	}

	var internalErr *apperror.InternalError
	if errors.As(err, &internalErr) {
		log.Printf("ERROR [%s] %s: %v", internalErr.Layer, internalErr.Op, internalErr.Err)
	} else {
		log.Printf("UNEXPECTED ERROR: %v", err)
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
		"code":    apperror.ErrInternalServer.Code,
		"message": apperror.ErrInternalServer.Error(),
	}})
}
