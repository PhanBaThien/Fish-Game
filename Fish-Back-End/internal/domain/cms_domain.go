package domain

import (
	"encoding/json"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
)

// ─── REQUEST / RESPONSE cho Lịch sử giao dịch ────────────────────────────────

type ListTransactionsRequest struct {
	PlayerID string `form:"playerId"`
	Type     string `form:"type"` // shot | win | gift | deposit | withdraw | adjust
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}

type ListTransactionsResponse struct {
	Items []models.Transaction `json:"items"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}

// ─── REQUEST / RESPONSE cho Cài đặt hệ thống ─────────────────────────────────

type UpsertSettingRequest struct {
	Value json.RawMessage `json:"value" binding:"required"`
}

// ─── REQUEST / RESPONSE cho Dashboard Stats ───────────────────────────────────

type StatsOverviewResponse struct {
	TotalPlayers      int64   `json:"total_players"`
	ActivePlayers     int64   `json:"active_players"`
	BannedPlayers     int64   `json:"banned_players"`
	TotalRooms        int64   `json:"total_rooms"`
	ActiveRooms       int64   `json:"active_rooms"`
	TotalTransactions int64   `json:"total_transactions"`
	TotalGoldIn       int64   `json:"total_gold_in"`   // tổng vàng nạp/gift
	TotalGoldOut      int64   `json:"total_gold_out"`  // tổng vàng rút/thắng
}

type TimeseriesRequest struct {
	Range string `form:"range"` // 7d | 30d | 90d
}

type TimeseriesPoint struct {
	Date       string `json:"date"`        // "2026-05-14"
	GoldIn     int64  `json:"gold_in"`
	GoldOut    int64  `json:"gold_out"`
	NewPlayers int64  `json:"new_players"`
	Shots      int64  `json:"shots"`
}

type TimeseriesResponse struct {
	Range  string            `json:"range"`
	Points []TimeseriesPoint `json:"points"`
}

// ─── REQUEST / RESPONSE cho Tìm kiếm tổng ────────────────────────────────────

type SearchResult struct {
	Players []models.Player `json:"players"`
	Rooms   []models.Room   `json:"rooms"`
	Fishes  []models.Fish   `json:"fishes"`
}
