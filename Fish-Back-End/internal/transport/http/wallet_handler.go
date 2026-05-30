package http

import (
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/transport/http/middleware"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/usecase"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletUsecase usecase.WalletUsecase
	tokenMaker    utils.TokenMaker
}

func NewWalletHandler(u usecase.WalletUsecase, m utils.TokenMaker) *WalletHandler {
	return &WalletHandler{walletUsecase: u, tokenMaker: m}
}

func (h *WalletHandler) RegisterRoutes(router *gin.RouterGroup) {
	wallet := router.Group("/wallet")
	wallet.Use(middleware.AuthMiddleware(h.tokenMaker))
	{
		wallet.GET("", h.GetWallet)
		wallet.GET("/transactions", h.GetTransactions)
		wallet.POST("/deposit", h.Deposit)
		wallet.POST("/withdraw", h.Withdraw)
		wallet.POST("/session/start", h.StartSession)
		wallet.POST("/session/end", h.EndSession)
	}
}

// GetWallet godoc
// GET /api/v1/wallet
func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	wallet, err := h.walletUsecase.GetWallet(c.Request.Context(), userID)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, wallet)
}

// GetTransactions godoc
// GET /api/v1/wallet/transactions?limit=20&offset=0
func (h *WalletHandler) GetTransactions(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	var req domain.TransactionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	txs, total, err := h.walletUsecase.GetTransactions(c.Request.Context(), userID, req.Limit, req.Offset)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, gin.H{
		"transactions": txs,
		"total":        total,
		"limit":        req.Limit,
		"offset":       req.Offset,
	})
}

// Deposit godoc
// POST /api/v1/wallet/deposit — nạp vàng
func (h *WalletHandler) Deposit(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	var req domain.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	wallet, err := h.walletUsecase.Deposit(c.Request.Context(), userID, &req)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, wallet)
}

// Withdraw godoc
// POST /api/v1/wallet/withdraw — rút vàng
func (h *WalletHandler) Withdraw(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	var req domain.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	wallet, err := h.walletUsecase.Withdraw(c.Request.Context(), userID, &req)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, wallet)
}

// StartSession godoc
// POST /api/v1/wallet/session/start
// Tạo game session khi người chơi vào phòng
func (h *WalletHandler) StartSession(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	var req domain.StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	session, err := h.walletUsecase.StartSession(c.Request.Context(), userID, &req)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, session)
}

// EndSession godoc
// POST /api/v1/wallet/session/end
// Kết thúc session, gộp earn/spend thành 2 transactions, cập nhật balance
func (h *WalletHandler) EndSession(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	var req domain.EndSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	session, wallet, err := h.walletUsecase.EndSession(c.Request.Context(), userID, &req)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, gin.H{
		"session": session,
		"wallet":  wallet,
	})
}
