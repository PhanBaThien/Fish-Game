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
		wallet.POST("/earn", h.Earn)
		wallet.POST("/spend", h.Spend)
	}
}

// GetWallet godoc
// GET /api/v1/wallet
// Trả về số dư hiện tại của user đang đăng nhập
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
// Lịch sử giao dịch của user, mới nhất trước
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

// Earn godoc
// POST /api/v1/wallet/earn
// Cộng vàng (khi bắn hạ cá)
func (h *WalletHandler) Earn(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	var req domain.EarnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	wallet, err := h.walletUsecase.Earn(c.Request.Context(), userID, &req)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, wallet)
}

// Spend godoc
// POST /api/v1/wallet/spend
// Trừ vàng (khi đặt cược). Trả ErrInsufficientBalance nếu không đủ
func (h *WalletHandler) Spend(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	var req domain.SpendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperror.ErrBadRequest)
		return
	}
	wallet, err := h.walletUsecase.Spend(c.Request.Context(), userID, &req)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, wallet)
}
