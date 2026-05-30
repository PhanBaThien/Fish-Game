package http

import (
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth   *AuthHandler
	Room   *RoomHandler
	Fish   *FishHandler
	Wallet *WalletHandler
	WS     *WSHandler
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
		if h.Room != nil {
			h.Room.RegisterRoutes(v1)
		}
		if h.Fish != nil {
			h.Fish.RegisterRoutes(v1)
		}
		if h.Wallet != nil {
			h.Wallet.RegisterRoutes(v1)
		}
		if h.WS != nil {
			h.WS.RegisterRoutes(v1)
		}
	}

	return router
}
