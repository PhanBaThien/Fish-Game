package repository

import (
	"context"
	"errors"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository/dbgen"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository interface {
	GetOrCreate(ctx context.Context, userID int64, initialBalance int64) (*models.Wallet, error)
	Deposit(ctx context.Context, userID int64, amount int64, description *string) (*models.Wallet, error)
	Withdraw(ctx context.Context, userID int64, amount int64, description *string) (*models.Wallet, error)
	GetTransactions(ctx context.Context, userID int64, limit, offset int32) ([]models.Transaction, error)
	CountTransactions(ctx context.Context, userID int64) (int64, error)
	StartSession(ctx context.Context, userID, roomID int64) (*models.GameSession, error)
	EndSession(ctx context.Context, userID int64, req EndSessionParams) (*models.GameSession, *models.Wallet, error)
}

// EndSessionParams chứa dữ liệu tổng kết từ client gửi lên
type EndSessionParams struct {
	SessionID  int64
	ShotsFired int32
	FishKilled int32
	TotalSpend int64
	TotalEarn  int64
}

type walletPgRepo struct {
	pool    *pgxpool.Pool
	queries *dbgen.Queries
}

func NewWalletRepository(pool *pgxpool.Pool) WalletRepository {
	return &walletPgRepo{pool: pool, queries: dbgen.New(pool)}
}

// ── Helpers map dbgen → models ────────────────────────────────────────────────

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
		SessionID:   pgint8ToPtr(t.SessionID),
		Amount:      t.Amount,
		Type:        t.Type,
		Description: pgtextToPtr(t.Description),
		CreatedAt:   t.CreatedAt.Time,
	}
}

func mapToModelGameSession(s dbgen.GameSession) models.GameSession {
	gs := models.GameSession{
		ID:         s.ID,
		UserID:     s.UserID,
		RoomID:     s.RoomID,
		ShotsFired: s.ShotsFired,
		FishKilled: s.FishKilled,
		TotalSpend: s.TotalSpend,
		TotalEarn:  s.TotalEarn,
		Status:     s.Status,
		StartedAt:  s.StartedAt.Time,
	}
	if s.EndedAt.Valid {
		t := s.EndedAt.Time
		gs.EndedAt = &t
	}
	return gs
}

// pgint8ToPtr / ptrToPgint8 — nullable int64 (dùng cho session_id)
func pgint8ToPtr(i pgtype.Int8) *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

func ptrToPgint8(i *int64) pgtype.Int8 {
	if i == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *i, Valid: true}
}

// ── GetOrCreate ───────────────────────────────────────────────────────────────

func (r *walletPgRepo) GetOrCreate(ctx context.Context, userID int64, initialBalance int64) (*models.Wallet, error) {
	row, err := r.queries.GetWallet(ctx, userID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.Wrap("repository", "walletRepo.GetOrCreate.Get", err)
		}
		row, err = r.queries.CreateWallet(ctx, dbgen.CreateWalletParams{
			UserID:  userID,
			Balance: initialBalance,
		})
		if err != nil {
			return nil, apperror.Wrap("repository", "walletRepo.GetOrCreate.Create", err)
		}
	}
	w := mapToModelWallet(row)
	return &w, nil
}

// ── transfer: cập nhật balance + ghi transaction trong 1 DB tx ───────────────

func (r *walletPgRepo) transfer(
	ctx context.Context,
	userID int64,
	amount int64,
	txType string,
	description *string,
	sessionID *int64,
) (*models.Wallet, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, apperror.Wrap("repository", "walletRepo.transfer.Begin", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	q := dbgen.New(tx)

	row, err := q.UpdateWalletBalance(ctx, dbgen.UpdateWalletBalanceParams{
		Balance: amount,
		UserID:  userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrInsufficientBalance
		}
		return nil, apperror.Wrap("repository", "walletRepo.transfer.UpdateBalance", err)
	}

	_, err = q.CreateTransaction(ctx, dbgen.CreateTransactionParams{
		UserID:      userID,
		SessionID:   ptrToPgint8(sessionID),
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

func (r *walletPgRepo) Deposit(ctx context.Context, userID int64, amount int64, description *string) (*models.Wallet, error) {
	return r.transfer(ctx, userID, amount, "deposit", description, nil)
}

func (r *walletPgRepo) Withdraw(ctx context.Context, userID int64, amount int64, description *string) (*models.Wallet, error) {
	return r.transfer(ctx, userID, -amount, "withdraw", description, nil)
}

// ── GetTransactions / CountTransactions ──────────────────────────────────────

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

// ── StartSession ──────────────────────────────────────────────────────────────

func (r *walletPgRepo) StartSession(ctx context.Context, userID, roomID int64) (*models.GameSession, error) {
	row, err := r.queries.CreateGameSession(ctx, dbgen.CreateGameSessionParams{
		UserID: userID,
		RoomID: roomID,
	})
	if err != nil {
		return nil, apperror.Wrap("repository", "walletRepo.StartSession", err)
	}
	gs := mapToModelGameSession(row)
	return &gs, nil
}

// ── EndSession: đóng session + gộp earn/spend thành transactions trong 1 DB tx

func (r *walletPgRepo) EndSession(ctx context.Context, userID int64, p EndSessionParams) (*models.GameSession, *models.Wallet, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, nil, apperror.Wrap("repository", "walletRepo.EndSession.Begin", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	q := dbgen.New(tx)

	// 1. Đóng session — WHERE status='active' ngăn double-end
	sessionRow, err := q.EndGameSession(ctx, dbgen.EndGameSessionParams{
		ID:         p.SessionID,
		ShotsFired: p.ShotsFired,
		FishKilled: p.FishKilled,
		TotalSpend: p.TotalSpend,
		TotalEarn:  p.TotalEarn,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, apperror.ErrSessionNotFound
		}
		return nil, nil, apperror.Wrap("repository", "walletRepo.EndSession.EndGameSession", err)
	}

	// Xác nhận session thuộc về user đang gọi
	if sessionRow.UserID != userID {
		return nil, nil, apperror.ErrForbidden
	}

	sidPtr := &p.SessionID

	// 2. Tính net và ghi 1 transaction type=play (amount có thể âm nếu thua)
	net := p.TotalEarn - p.TotalSpend
	_, err = q.CreateTransaction(ctx, dbgen.CreateTransactionParams{
		UserID:    userID,
		SessionID: ptrToPgint8(sidPtr),
		Amount:    net,
		Type:      "play",
	})
	if err != nil {
		return nil, nil, apperror.Wrap("repository", "walletRepo.EndSession.CreatePlayTx", err)
	}

	// 3. Cập nhật balance theo net = earn - spend
	var walletRow dbgen.Wallet
	if net != 0 {
		walletRow, err = q.UpdateWalletBalance(ctx, dbgen.UpdateWalletBalanceParams{
			Balance: net,
			UserID:  userID,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil, apperror.ErrInsufficientBalance
			}
			return nil, nil, apperror.Wrap("repository", "walletRepo.EndSession.UpdateBalance", err)
		}
	} else {
		// net = 0: chỉ đọc balance hiện tại, không cần update
		walletRow, err = q.GetWallet(ctx, userID)
		if err != nil {
			return nil, nil, apperror.Wrap("repository", "walletRepo.EndSession.GetWallet", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, nil, apperror.Wrap("repository", "walletRepo.EndSession.Commit", err)
	}

	gs := mapToModelGameSession(sessionRow)
	w := mapToModelWallet(walletRow)
	return &gs, &w, nil
}
