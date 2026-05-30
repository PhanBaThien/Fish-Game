package usecase

import (
	"context"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository"
)

type WalletUsecase interface {
	GetWallet(ctx context.Context, userID int64) (*models.Wallet, error)
	Deposit(ctx context.Context, userID int64, req *domain.DepositRequest) (*models.Wallet, error)
	Withdraw(ctx context.Context, userID int64, req *domain.WithdrawRequest) (*models.Wallet, error)
	GetTransactions(ctx context.Context, userID int64, limit, offset int32) ([]models.Transaction, int64, error)
	StartSession(ctx context.Context, userID int64, req *domain.StartSessionRequest) (*models.GameSession, error)
	EndSession(ctx context.Context, userID int64, req *domain.EndSessionRequest) (*models.GameSession, *models.Wallet, error)
}

type walletUsecase struct {
	walletRepo repository.WalletRepository
}

func NewWalletUsecase(repo repository.WalletRepository) WalletUsecase {
	return &walletUsecase{walletRepo: repo}
}

func (u *walletUsecase) GetWallet(ctx context.Context, userID int64) (*models.Wallet, error) {
	return u.walletRepo.GetOrCreate(ctx, userID, 5000)
}

func (u *walletUsecase) Deposit(ctx context.Context, userID int64, req *domain.DepositRequest) (*models.Wallet, error) {
	return u.walletRepo.Deposit(ctx, userID, req.Amount, req.Description)
}

func (u *walletUsecase) Withdraw(ctx context.Context, userID int64, req *domain.WithdrawRequest) (*models.Wallet, error) {
	return u.walletRepo.Withdraw(ctx, userID, req.Amount, req.Description)
}

func (u *walletUsecase) GetTransactions(ctx context.Context, userID int64, limit, offset int32) ([]models.Transaction, int64, error) {
	txs, err := u.walletRepo.GetTransactions(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := u.walletRepo.CountTransactions(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	return txs, total, nil
}

func (u *walletUsecase) StartSession(ctx context.Context, userID int64, req *domain.StartSessionRequest) (*models.GameSession, error) {
	return u.walletRepo.StartSession(ctx, userID, req.RoomID)
}

func (u *walletUsecase) EndSession(ctx context.Context, userID int64, req *domain.EndSessionRequest) (*models.GameSession, *models.Wallet, error) {
	return u.walletRepo.EndSession(ctx, userID, repository.EndSessionParams{
		SessionID:  req.SessionID,
		ShotsFired: req.ShotsFired,
		FishKilled: req.FishKilled,
		TotalSpend: req.TotalSpend,
		TotalEarn:  req.TotalEarn,
	})
}
