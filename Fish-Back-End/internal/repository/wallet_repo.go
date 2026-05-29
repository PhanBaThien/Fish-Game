package repository

import (
	"context"
	"errors"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository/dbgen"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository interface {
	GetOrCreate(ctx context.Context, userID int64) (*models.Wallet, error)
	Credit(ctx context.Context, userID int64, amount int64, txType string, description *string) (*models.Wallet, error)
	Debit(ctx context.Context, userID int64, amount int64, txType string, description *string) (*models.Wallet, error)
	GetTransactions(ctx context.Context, userID int64, limit, offset int32) ([]models.Transaction, error)
	CountTransactions(ctx context.Context, userID int64) (int64, error)
}

type walletPgRepo struct {
	pool    *pgxpool.Pool
	queries *dbgen.Queries
}

func NewWalletRepository(pool *pgxpool.Pool) WalletRepository {
	return &walletPgRepo{pool: pool, queries: dbgen.New(pool)}
}

func mapToModelWallet(w dbgen.Wallet) models.Wallet {
	return models.Wallet{
		UserID:    w.UserID,
		Balance:   w.Balance,
		UpdatedAt: w.UpdatedAt.Time,
	}
}

func mapToModelTransaction(t dbgen.Transaction) models.Transaction {
	return models.Transaction{
		ID:          t.ID,
		UserID:      t.UserID,
		Amount:      t.Amount,
		Type:        t.Type,
		Description: pgtextToPtr(t.Description),
		CreatedAt:   t.CreatedAt.Time,
	}
}

func (r *walletPgRepo) GetOrCreate(ctx context.Context, userID int64) (*models.Wallet, error) {
	row, err := r.queries.GetWallet(ctx, userID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.Wrap("repository", "walletRepo.GetOrCreate.Get", err)
		}
		// wallet chưa tồn tại → tạo mới (fallback cho user cũ chưa có trigger)
		row, err = r.queries.CreateWallet(ctx, userID)
		if err != nil {
			return nil, apperror.Wrap("repository", "walletRepo.GetOrCreate.Create", err)
		}
	}
	w := mapToModelWallet(row)
	return &w, nil
}

// transfer thực hiện cập nhật balance + ghi transaction trong 1 DB transaction
func (r *walletPgRepo) transfer(ctx context.Context, userID int64, amount int64, txType string, description *string) (*models.Wallet, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, apperror.Wrap("repository", "walletRepo.transfer.Begin", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	q := dbgen.New(tx) // dùng queries trong context của transaction

	// Cập nhật balance — câu WHERE tự kiểm tra balance >= 0
	row, err := q.UpdateWalletBalance(ctx, dbgen.UpdateWalletBalanceParams{
		Amount: amount,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrInsufficientBalance
		}
		return nil, apperror.Wrap("repository", "walletRepo.transfer.UpdateBalance", err)
	}

	// Ghi lịch sử giao dịch
	_, err = q.CreateTransaction(ctx, dbgen.CreateTransactionParams{
		UserID:      userID,
		Amount:      amount,
		Type:        txType,
		Description: ptrToPgtext(description),
	})
	if err != nil {
		return nil, apperror.Wrap("repository", "walletRepo.transfer.CreateTransaction", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, apperror.Wrap("repository", "walletRepo.transfer.Commit", err)
	}

	w := mapToModelWallet(row)
	return &w, nil
}

func (r *walletPgRepo) Credit(ctx context.Context, userID int64, amount int64, txType string, description *string) (*models.Wallet, error) {
	return r.transfer(ctx, userID, amount, txType, description) // amount dương
}

func (r *walletPgRepo) Debit(ctx context.Context, userID int64, amount int64, txType string, description *string) (*models.Wallet, error) {
	return r.transfer(ctx, userID, -amount, txType, description) // amount âm
}

func (r *walletPgRepo) GetTransactions(ctx context.Context, userID int64, limit, offset int32) ([]models.Transaction, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.queries.ListTransactionsByUserID(ctx, dbgen.ListTransactionsByUserIDParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, apperror.Wrap("repository", "walletRepo.GetTransactions", err)
	}
	txs := make([]models.Transaction, len(rows))
	for i, row := range rows {
		txs[i] = mapToModelTransaction(row)
	}
	return txs, nil
}

func (r *walletPgRepo) CountTransactions(ctx context.Context, userID int64) (int64, error) {
	count, err := r.queries.CountTransactionsByUserID(ctx, userID)
	if err != nil {
		return 0, apperror.Wrap("repository", "walletRepo.CountTransactions", err)
	}
	return count, nil
}
