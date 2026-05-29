package domain

const (
	TxTypeEarn  = "earn"  // nhận vàng khi bắn cá
	TxTypeSpend = "spend" // tiêu vàng khi đặt cược
)

type EarnRequest struct {
	Amount      int64   `json:"amount"      binding:"required,min=1"`
	Description *string `json:"description"`
}

type SpendRequest struct {
	Amount      int64   `json:"amount"      binding:"required,min=1"`
	Description *string `json:"description"`
}

type TransactionListRequest struct {
	Limit  int32 `form:"limit"  binding:"omitempty,min=1,max=100"`
	Offset int32 `form:"offset" binding:"omitempty,min=0"`
}
