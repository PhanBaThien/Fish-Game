package usecase

import (
	"context"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository"
)

type WalletUsecase interface {
	GetWallet(ctx context.Context, userID int64) (*models.Wallet, error)
	Earn(ctx context.Context, userID int64, req *domain.EarnRequest) (*models.Wallet, error)
	Spend(ctx context.Context, userID int64, req *domain.SpendRequest) (*models.Wallet, error)
	GetTransactions(ctx context.Context, userID int64, limit, offset int32) ([]models.Transaction, int64, error)
}

type walletUsecase struct {
	walletRepo repository.WalletRepository
}

func NewWalletUsecase(repo repository.WalletRepository) WalletUsecase {
	return &walletUsecase{walletRepo: repo}
}

func (u *walletUsecase) GetWallet(ctx context.Context, userID int64) (*models.Wallet, error) {
	return u.walletRepo.GetOrCreate(ctx, userID)
}

func (u *walletUsecase) Earn(ctx context.Context, userID int64, req *domain.EarnRequest) (*models.Wallet, error) {
	return u.walletRepo.Credit(ctx, userID, req.Amount, domain.TxTypeEarn, req.Description)
}

func (u *walletUsecase) Spend(ctx context.Context, userID int64, req *domain.SpendRequest) (*models.Wallet, error) {
	return u.walletRepo.Debit(ctx, userID, req.Amount, domain.TxTypeSpend, req.Description)
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
