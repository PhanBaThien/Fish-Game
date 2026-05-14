package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
)

var (
	ErrPlayerNotFound     = errors.New("không tìm thấy người chơi")
	ErrPlayerUsernameTaken = errors.New("tên đăng nhập đã được sử dụng")
)

// PlayerUsecase định nghĩa nghiệp vụ quản lý người chơi
type PlayerUsecase interface {
	ListPlayers(ctx context.Context, req *domain.ListPlayersRequest) (*domain.ListPlayersResponse, error)
	GetPlayer(ctx context.Context, id string) (*models.Player, error)
	CreatePlayer(ctx context.Context, req *domain.CreatePlayerRequest) (*models.Player, error)
	UpdatePlayer(ctx context.Context, id string, req *domain.UpdatePlayerRequest) (*models.Player, error)
	DeletePlayer(ctx context.Context, id string) error
	BanPlayer(ctx context.Context, id string, req *domain.BanPlayerRequest) error
	GiftPlayer(ctx context.Context, id, actorID string, req *domain.GiftPlayerRequest) error
	UpdatePlayerRTP(ctx context.Context, id string, req *domain.UpdatePlayerRTPRequest) error
}

type playerUsecase struct {
	playerRepo  repository.PlayerRepository
	txRepo      repository.TransactionRepository
	hasher      utils.PasswordHasher
}

func NewPlayerUsecase(
	playerRepo repository.PlayerRepository,
	txRepo repository.TransactionRepository,
	hasher utils.PasswordHasher,
) PlayerUsecase {
	return &playerUsecase{
		playerRepo: playerRepo,
		txRepo:     txRepo,
		hasher:     hasher,
	}
}

func (u *playerUsecase) ListPlayers(ctx context.Context, req *domain.ListPlayersRequest) (*domain.ListPlayersResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}
	offset := (req.Page - 1) * req.Limit

	players, total, err := u.playerRepo.List(ctx, req.Status, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("playerUsecase.ListPlayers: %w", err)
	}

	return &domain.ListPlayersResponse{
		Items: players,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

func (u *playerUsecase) GetPlayer(ctx context.Context, id string) (*models.Player, error) {
	p, err := u.playerRepo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrPlayerNotFound) {
		return nil, ErrPlayerNotFound
	}
	return p, err
}

func (u *playerUsecase) CreatePlayer(ctx context.Context, req *domain.CreatePlayerRequest) (*models.Player, error) {
	exists, err := u.playerRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPlayerUsernameTaken
	}

	hash, err := u.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, ErrInternalServer
	}

	p := &models.Player{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		GoldBalance:  0,
		Status:       "active",
		WinRate:      100.0,
	}

	if err := u.playerRepo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("playerUsecase.CreatePlayer: %w", err)
	}
	return p, nil
}

func (u *playerUsecase) UpdatePlayer(ctx context.Context, id string, req *domain.UpdatePlayerRequest) (*models.Player, error) {
	p, err := u.playerRepo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrPlayerNotFound) {
		return nil, ErrPlayerNotFound
	}
	if err != nil {
		return nil, err
	}

	p.Username = req.Username
	p.Email = req.Email
	p.Status = req.Status

	if err := u.playerRepo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("playerUsecase.UpdatePlayer: %w", err)
	}
	return p, nil
}

func (u *playerUsecase) DeletePlayer(ctx context.Context, id string) error {
	err := u.playerRepo.Delete(ctx, id)
	if errors.Is(err, repository.ErrPlayerNotFound) {
		return ErrPlayerNotFound
	}
	return err
}

func (u *playerUsecase) BanPlayer(ctx context.Context, id string, req *domain.BanPlayerRequest) error {
	_, err := u.playerRepo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrPlayerNotFound) {
		return ErrPlayerNotFound
	}
	if err != nil {
		return err
	}
	return u.playerRepo.UpdateStatus(ctx, id, "banned")
}

func (u *playerUsecase) GiftPlayer(ctx context.Context, id, actorID string, req *domain.GiftPlayerRequest) error {
	p, err := u.playerRepo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrPlayerNotFound) {
		return ErrPlayerNotFound
	}
	if err != nil {
		return err
	}

	// Tạo idempotency key dựa trên actorID + playerID + timestamp (nano)
	idempotencyKey := fmt.Sprintf("gift-%s-%s-%d", actorID, id, time.Now().UnixNano())

	newBalance := p.GoldBalance + req.Amount

	tx := &models.Transaction{
		PlayerID:       id,
		Type:           "gift",
		Amount:         req.Amount,
		BalanceAfter:   newBalance,
		IdempotencyKey: idempotencyKey,
	}

	if err := u.txRepo.Create(ctx, tx); err != nil {
		return fmt.Errorf("playerUsecase.GiftPlayer create tx: %w", err)
	}

	return nil
}

func (u *playerUsecase) UpdatePlayerRTP(ctx context.Context, id string, req *domain.UpdatePlayerRTPRequest) error {
	_, err := u.playerRepo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrPlayerNotFound) {
		return ErrPlayerNotFound
	}
	if err != nil {
		return err
	}
	return u.playerRepo.UpdateWinRate(ctx, id, req.WinRate)
}
