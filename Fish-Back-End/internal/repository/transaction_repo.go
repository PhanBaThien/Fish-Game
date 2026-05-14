package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
)

// TransactionRepository định nghĩa interface thao tác bảng transactions
type TransactionRepository interface {
	List(ctx context.Context, playerID, txType string, limit, offset int) ([]models.Transaction, int64, error)
	Create(ctx context.Context, tx *models.Transaction) error
}

type transactionPgRepo struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionPgRepo{db: db}
}

func (r *transactionPgRepo) List(ctx context.Context, playerID, txType string, limit, offset int) ([]models.Transaction, int64, error) {
	query := `
		SELECT id, player_id, type, amount, balance_after,
		       ref_shot_id, ref_kill_id, idempotency_key, metadata, created_at, updated_at
		FROM transactions
		WHERE ($1 = '' OR player_id::text = $1)
		  AND ($2 = '' OR type = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, playerID, txType, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("txRepo.List: %w", err)
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var metaRaw []byte
		if err := rows.Scan(
			&t.ID, &t.PlayerID, &t.Type, &t.Amount, &t.BalanceAfter,
			&t.RefShotID, &t.RefKillID, &t.IdempotencyKey, &metaRaw,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("txRepo.List scan: %w", err)
		}
		if metaRaw != nil {
			t.Metadata = json.RawMessage(metaRaw)
		}
		txs = append(txs, t)
	}

	var total int64
	countQ := `SELECT COUNT(*) FROM transactions WHERE ($1 = '' OR player_id::text = $1) AND ($2 = '' OR type = $2)`
	if err := r.db.QueryRowContext(ctx, countQ, playerID, txType).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("txRepo.List count: %w", err)
	}

	return txs, total, nil
}

func (r *transactionPgRepo) Create(ctx context.Context, t *models.Transaction) error {
	query := `
		INSERT INTO transactions
		  (player_id, type, amount, balance_after, ref_shot_id, ref_kill_id, idempotency_key, metadata)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		t.PlayerID, t.Type, t.Amount, t.BalanceAfter,
		t.RefShotID, t.RefKillID, t.IdempotencyKey, t.Metadata,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}
