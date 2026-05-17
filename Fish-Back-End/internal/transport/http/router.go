package http

import (
	"github.com/gin-gonic/gin"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/transport/http/middleware"
)

type Handlers struct {
	Auth   *AuthHandler
}
func SetupRouter(h Handlers) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	v1 := router.Group("/api/v1")
	{
		if h.Auth != nil {
			h.Auth.RegisterRoutes(v1)
		}
		// if h.Player != nil {
		// 	h.Player.RegisterRoutes(v1)
		// }
	}

	return router
}