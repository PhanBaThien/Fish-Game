package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/usecase"
	gorillaws "github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMsgSize = 1024
	dbTimeout  = 5 * time.Second // timeout cho mọi DB call
)

// msgHandlerFunc là chữ ký thống nhất cho mọi handler.
type msgHandlerFunc func(payload json.RawMessage)

// Client đại diện cho 1 kết nối WebSocket của 1 người chơi.
type Client struct {
	hub  *Hub
	conn *gorillaws.Conn
	send chan []byte
	once sync.Once // đảm bảo send channel chỉ bị close 1 lần

	userID int64

	// Dispatch table: type → handler
	handlers map[string]msgHandlerFunc

	// Usecase
	walletUsecase usecase.WalletUsecase
	roomUsecase   usecase.RoomUsecase
	fishUsecase   usecase.FishUsecase

	// fishMap: fishID → rewardMultiplier (load từ DB lúc join_room)
	// Server tự tra, không tin giá trị client gửi lên
	fishMap map[int32]int32

	// Trạng thái session hiện tại
	sessionID int64
	roomID    int64

	// Counters theo dõi server-side (authoritative)
	shotsFired int32
	fishKilled int32
	totalSpend int64
	totalEarn  int64
	lastBet    int64 // bet của viên đạn cuối cùng — dùng để tính earn khi cá chết

	// Balance ước tính trong ván (tránh gọi DB mỗi lần)
	estimatedBalance int64
}

func NewClient(
	hub *Hub,
	conn *gorillaws.Conn,
	userID int64,
	walletUC usecase.WalletUsecase,
	roomUC usecase.RoomUsecase,
	fishUC usecase.FishUsecase,
) *Client {
	c := &Client{
		hub:           hub,
		conn:          conn,
		send:          make(chan []byte, 512), // tăng buffer: game bắn nhanh
		userID:        userID,
		walletUsecase: walletUC,
		roomUsecase:   roomUC,
		fishUsecase:   fishUC,
		fishMap:       make(map[int32]int32),
	}
	c.handlers = map[string]msgHandlerFunc{
		MsgJoinRoom:   c.handleJoinRoom,
		MsgShoot:      c.handleShoot,
		MsgFishKilled: c.handleFishKilled,
		MsgLeaveRoom:  c.handleLeaveRoom,
		MsgPing:       c.handlePing,
	}
	return c
}

// dbCtx trả về context với timeout 5s cho mọi DB call.
func dbCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), dbTimeout)
}

// ── Pumps ─────────────────────────────────────────────────────────────────────

// ReadPump đọc message từ WebSocket và xử lý.
// Chạy trong goroutine riêng.
func (c *Client) ReadPump() {
	defer func() {
		c.endSessionIfActive() // kết thúc session khi disconnect
		if c.roomID != 0 {
			c.hub.LeaveRoom(c, c.roomID)
		}
		c.closeSend()
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.handleMessage(raw)
	}
}

// WritePump ghi message từ channel send ra WebSocket.
// Chạy trong goroutine riêng.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(gorillaws.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(gorillaws.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(gorillaws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ── Message dispatcher ────────────────────────────────────────────────────────

func (c *Client) handleMessage(raw []byte) {
	var msg InMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		c.sendError("BAD_REQUEST", "invalid message format")
		return
	}

	handler, ok := c.handlers[msg.Type]
	if !ok {
		c.sendError("UNKNOWN_TYPE", "unknown message type: "+msg.Type)
		return
	}
	handler(msg.Payload)
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func (c *Client) handleJoinRoom(payload json.RawMessage) {
	// Fix: không cho join khi đang trong phòng → tránh orphan session
	if c.sessionID != 0 {
		c.sendError("ALREADY_IN_ROOM", "leave current room first")
		return
	}

	var p JoinRoomPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.sendError("BAD_REQUEST", "invalid join_room payload")
		return
	}

	ctx, cancel := dbCtx()
	defer cancel()

	// Kiểm tra phòng tồn tại
	_, err := c.roomUsecase.GetByID(ctx, p.RoomID)
	if err != nil {
		c.sendError("ROOM_NOT_FOUND", "room not found")
		return
	}

	// Lấy balance hiện tại làm điểm khởi đầu ước tính
	wallet, err := c.walletUsecase.GetWallet(ctx, c.userID)
	if err != nil {
		c.sendError("WALLET_ERROR", "cannot get wallet")
		return
	}

	// Load danh sách cá → fishMap (server tự tra multiplier, không tin client)
	fishList, err := c.fishUsecase.List(ctx)
	if err != nil {
		c.sendError("FISH_ERROR", "cannot load fish data")
		return
	}
	for _, f := range fishList {
		c.fishMap[f.ID] = f.RewardMultiplier
	}

	// Tạo session
	session, err := c.walletUsecase.StartSession(ctx, c.userID, &domain.StartSessionRequest{
		RoomID: p.RoomID,
	})
	if err != nil {
		c.sendError("SESSION_ERROR", "cannot start session")
		return
	}

	// Reset trạng thái ván chơi
	c.sessionID        = session.ID
	c.roomID           = p.RoomID
	c.shotsFired       = 0
	c.fishKilled       = 0
	c.totalSpend       = 0
	c.totalEarn        = 0
	c.estimatedBalance = wallet.Balance

	c.hub.JoinRoom(c, p.RoomID)

	c.sendJSON(MsgSessionStarted, SessionStartedPayload{SessionID: session.ID})
}

func (c *Client) handleShoot(payload json.RawMessage) {
	if c.sessionID == 0 {
		c.sendError("NO_SESSION", "join a room first")
		return
	}

	var p ShootPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.sendError("BAD_REQUEST", "invalid shoot payload")
		return
	}

	// Validate bet_amount: phải trong khoảng 10–100 (cố định)
	const minBet, maxBet = int64(10), int64(100)
	if p.BetAmount < minBet || p.BetAmount > maxBet {
		c.sendError("INVALID_BET", fmt.Sprintf("bet must be between %d and %d", minBet, maxBet))
		return
	}

	// Kiểm tra balance đủ không
	if c.estimatedBalance < p.BetAmount {
		c.sendError("INSUFFICIENT_BALANCE", "not enough balance to shoot")
		return
	}

	c.shotsFired++
	c.totalSpend       += p.BetAmount
	c.estimatedBalance -= p.BetAmount
	c.lastBet           = p.BetAmount // lưu lại để tính earn khi cá chết

	c.sendJSON(MsgShootAck, ShootAckPayload{
		ShotsFired: c.shotsFired,
		TotalSpend: c.totalSpend,
		Balance:    c.estimatedBalance,
	})
}

func (c *Client) handleFishKilled(payload json.RawMessage) {
	if c.sessionID == 0 {
		c.sendError("NO_SESSION", "join a room first")
		return
	}

	var p FishKilledPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.sendError("BAD_REQUEST", "invalid fish_killed payload")
		return
	}

	// Fix: server tự tra multiplier từ fishMap, không dùng giá trị client gửi
	multiplier, ok := c.fishMap[p.FishID]
	if !ok {
		c.sendError("INVALID_FISH", "unknown fish id")
		return
	}

	// dùng lastBet — viên đạn bắn trúng con cá này
	bet := c.lastBet
	if bet == 0 {
		bet = 10 // fallback nếu chưa có lastBet
	}
	earned := int64(multiplier) * bet
	c.fishKilled++
	c.totalEarn        += earned
	c.estimatedBalance += earned

	c.sendJSON(MsgEarnAck, EarnAckPayload{
		Amount:     earned,
		Balance:    c.estimatedBalance,
		TotalEarn:  c.totalEarn,
		FishKilled: c.fishKilled,
	})
}

func (c *Client) handleLeaveRoom(_ json.RawMessage) {
	c.endSessionIfActive()
	if c.roomID != 0 {
		c.hub.LeaveRoom(c, c.roomID)
		c.roomID = 0
	}
}

func (c *Client) handlePing(_ json.RawMessage) {
	c.sendJSON(MsgPong, nil)
}

// ── Session lifecycle ─────────────────────────────────────────────────────────

// endSessionIfActive gọi EndSession và gửi session_ended về client.
// Dùng sessionID=0 làm guard để tránh double-end.
func (c *Client) endSessionIfActive() {
	if c.sessionID == 0 {
		return
	}
	sid := c.sessionID
	c.sessionID = 0 // guard: ngăn gọi lại

	ctx, cancel := dbCtx()
	defer cancel()

	session, wallet, err := c.walletUsecase.EndSession(ctx, c.userID, &domain.EndSessionRequest{
		SessionID:  sid,
		ShotsFired: c.shotsFired,
		FishKilled: c.fishKilled,
		TotalSpend: c.totalSpend,
		TotalEarn:  c.totalEarn,
	})
	if err != nil {
		log.Printf("[ws] endSession user=%d err=%v", c.userID, err)
		return
	}

	c.sendJSON(MsgSessionEnded, SessionEndedPayload{
		Session: session,
		Wallet:  wallet,
	})
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func (c *Client) sendJSON(msgType string, payload any) {
	data, err := json.Marshal(OutMessage{Type: msgType, Payload: payload})
	if err != nil {
		return
	}
	select {
	case c.send <- data:
	default:
		// buffer đầy → client quá chậm, bỏ qua
	}
}

func (c *Client) sendError(code, message string) {
	c.sendJSON(MsgError, ErrorPayload{Code: code, Message: message})
}

func (c *Client) closeSend() {
	c.once.Do(func() { close(c.send) })
}
