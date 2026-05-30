package http

import (
	"net/http"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/transport/http/middleware"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/usecase"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/ws"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Cho phép mọi origin trong dev; production nên check c.Request.Header.Get("Origin")
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
	hub           *ws.Hub
	walletUsecase usecase.WalletUsecase
	roomUsecase   usecase.RoomUsecase
	fishUsecase   usecase.FishUsecase
	tokenMaker    utils.TokenMaker
}

func NewWSHandler(
	hub *ws.Hub,
	walletUC usecase.WalletUsecase,
	roomUC usecase.RoomUsecase,
	fishUC usecase.FishUsecase,
	tm utils.TokenMaker,
) *WSHandler {
	return &WSHandler{
		hub:           hub,
		walletUsecase: walletUC,
		roomUsecase:   roomUC,
		fishUsecase:   fishUC,
		tokenMaker:    tm,
	}
}

func (h *WSHandler) RegisterRoutes(router *gin.RouterGroup) {
	// WSAuthMiddleware đọc token từ ?token= thay vì Authorization header
	// vì browser WebSocket API không hỗ trợ custom header khi connect
	router.GET("/ws", middleware.WSAuthMiddleware(h.tokenMaker), h.ServeWS)
}

// ServeWS godoc
// GET /api/v1/ws?token=<access_token>
func (h *WSHandler) ServeWS(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// upgrader đã tự ghi response lỗi
		return
	}

	client := ws.NewClient(h.hub, conn, userID, h.walletUsecase, h.roomUsecase, h.fishUsecase)
	go client.WritePump()
	go client.ReadPump()
}
