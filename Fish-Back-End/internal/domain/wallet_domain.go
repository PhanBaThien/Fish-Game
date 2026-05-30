package domain

const (
	TxTypePlay     = "play"     // chơi game (gộp earn - spend), luôn có session_id
	TxTypeDeposit  = "deposit"  // nạp vàng, amount > 0
	TxTypeWithdraw = "withdraw" // rút vàng, amount < 0

	SessionStatusActive   = "active"
	SessionStatusFinished = "finished"
)

type DepositRequest struct {
	Amount      int64   `json:"amount"      binding:"required,min=1"`
	Description *string `json:"description"`
}

type WithdrawRequest struct {
	Amount      int64   `json:"amount"      binding:"required,min=1"`
	Description *string `json:"description"`
}

type TransactionListRequest struct {
	Limit  int32 `form:"limit"  binding:"omitempty,min=1,max=100"`
	Offset int32 `form:"offset" binding:"omitempty,min=0"`
}

// StartSessionRequest — vào phòng, tạo session mới
type StartSessionRequest struct {
	RoomID int64 `json:"room_id" binding:"required,min=1"`
}

// EndSessionRequest — thoát phòng, gộp thành 1 transaction type=play
type EndSessionRequest struct {
	SessionID  int64 `json:"session_id"  binding:"required,min=1"`
	ShotsFired int32 `json:"shots_fired" binding:"min=0"`
	FishKilled int32 `json:"fish_killed" binding:"min=0"`
	TotalSpend int64 `json:"total_spend" binding:"min=0"`
	TotalEarn  int64 `json:"total_earn"  binding:"min=0"`
}
