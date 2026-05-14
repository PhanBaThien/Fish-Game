package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository"
)

// TransactionUsecase định nghĩa nghiệp vụ quản lý giao dịch
type TransactionUsecase interface {
	ListTransactions(ctx context.Context, req *domain.ListTransactionsRequest) (*domain.ListTransactionsResponse, error)
}

// SettingUsecase định nghĩa nghiệp vụ cài đặt hệ thống
type SettingUsecase interface {
	GetSettings(ctx context.Context) ([]models.Setting, error)
	UpsertSetting(ctx context.Context, key, actorID string, req *domain.UpsertSettingRequest) (*models.Setting, error)
}

// StatsUsecase định nghĩa nghiệp vụ thống kê dashboard
type StatsUsecase interface {
	GetOverview(ctx context.Context) (*domain.StatsOverviewResponse, error)
	GetTimeseries(ctx context.Context, req *domain.TimeseriesRequest) (*domain.TimeseriesResponse, error)
}

// SearchUsecase định nghĩa nghiệp vụ tìm kiếm tổng
type SearchUsecase interface {
	Search(ctx context.Context, q string) (*domain.SearchResult, error)
}

// ─── Transaction ──────────────────────────────────────────────────────────────

type transactionUsecase struct {
	txRepo repository.TransactionRepository
}

func NewTransactionUsecase(txRepo repository.TransactionRepository) TransactionUsecase {
	return &transactionUsecase{txRepo: txRepo}
}

func (u *transactionUsecase) ListTransactions(ctx context.Context, req *domain.ListTransactionsRequest) (*domain.ListTransactionsResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}
	offset := (req.Page - 1) * req.Limit

	txs, total, err := u.txRepo.List(ctx, req.PlayerID, req.Type, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("txUsecase.ListTransactions: %w", err)
	}

	return &domain.ListTransactionsResponse{
		Items: txs,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// ─── Setting ──────────────────────────────────────────────────────────────────

type settingUsecase struct {
	settingRepo repository.SettingRepository
}

func NewSettingUsecase(settingRepo repository.SettingRepository) SettingUsecase {
	return &settingUsecase{settingRepo: settingRepo}
}

func (u *settingUsecase) GetSettings(ctx context.Context) ([]models.Setting, error) {
	return u.settingRepo.GetAll(ctx)
}

func (u *settingUsecase) UpsertSetting(ctx context.Context, key, actorID string, req *domain.UpsertSettingRequest) (*models.Setting, error) {
	s := &models.Setting{
		Key:       key,
		Value:     req.Value,
		UpdatedBy: actorID,
	}
	if err := u.settingRepo.Upsert(ctx, s); err != nil {
		return nil, fmt.Errorf("settingUsecase.UpsertSetting: %w", err)
	}
	return s, nil
}

// ─── Stats ────────────────────────────────────────────────────────────────────

type statsUsecase struct {
	statsRepo repository.StatsRepository
}

func NewStatsUsecase(statsRepo repository.StatsRepository) StatsUsecase {
	return &statsUsecase{statsRepo: statsRepo}
}

func (u *statsUsecase) GetOverview(ctx context.Context) (*domain.StatsOverviewResponse, error) {
	return u.statsRepo.Overview(ctx)
}

func (u *statsUsecase) GetTimeseries(ctx context.Context, req *domain.TimeseriesRequest) (*domain.TimeseriesResponse, error) {
	now := time.Now().UTC()
	var from time.Time

	switch req.Range {
	case "30d":
		from = now.AddDate(0, 0, -30)
	case "90d":
		from = now.AddDate(0, 0, -90)
	default: // "7d" và mặc định
		from = now.AddDate(0, 0, -7)
		req.Range = "7d"
	}

	points, err := u.statsRepo.Timeseries(ctx, from, now)
	if err != nil {
		return nil, fmt.Errorf("statsUsecase.GetTimeseries: %w", err)
	}

	return &domain.TimeseriesResponse{
		Range:  req.Range,
		Points: points,
	}, nil
}

// ─── Search ───────────────────────────────────────────────────────────────────

type searchUsecase struct {
	playerRepo repository.PlayerRepository
	roomRepo   repository.RoomRepository
	fishRepo   repository.FishRepository
}

func NewSearchUsecase(
	playerRepo repository.PlayerRepository,
	roomRepo repository.RoomRepository,
	fishRepo repository.FishRepository,
) SearchUsecase {
	return &searchUsecase{
		playerRepo: playerRepo,
		roomRepo:   roomRepo,
		fishRepo:   fishRepo,
	}
}

func (u *searchUsecase) Search(ctx context.Context, q string) (*domain.SearchResult, error) {
	if len(q) < 2 {
		return &domain.SearchResult{}, nil
	}

	const searchLimit = 5

	players, err := u.playerRepo.Search(ctx, q, searchLimit)
	if err != nil {
		return nil, fmt.Errorf("searchUsecase.Search players: %w", err)
	}

	rooms, err := u.roomRepo.Search(ctx, q, searchLimit)
	if err != nil {
		return nil, fmt.Errorf("searchUsecase.Search rooms: %w", err)
	}

	fishes, err := u.fishRepo.Search(ctx, q, searchLimit)
	if err != nil {
		return nil, fmt.Errorf("searchUsecase.Search fishes: %w", err)
	}

	return &domain.SearchResult{
		Players: players,
		Rooms:   rooms,
		Fishes:  fishes,
	}, nil
}
