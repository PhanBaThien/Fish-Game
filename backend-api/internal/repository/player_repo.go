package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/yourname/fish-game-backend/internal/models"
	"github.com/yourname/fish-game-backend/pkg/database"
)

// PlayerRepository defines the interface for player data access.
type PlayerRepository interface {
	FindAll(ctx context.Context) ([]models.Player, error)
	FindByID(ctx context.Context, id string) (*models.Player, error)
	Update(ctx context.Context, id string, req models.UpdatePlayerRequest) error
	Delete(ctx context.Context, id string) error
}

// playerSQLServerRepo is the SQL Server implementation of PlayerRepository.
type playerSQLServerRepo struct {
	db *database.DB
}

// NewPlayerRepository creates a new SQL Server–backed PlayerRepository.
func NewPlayerRepository(db *database.DB) PlayerRepository {
	return &playerSQLServerRepo{db: db}
}

// FindAll returns all players ordered by creation date descending.
func (r *playerSQLServerRepo) FindAll(ctx context.Context) ([]models.Player, error) {
	const query = `
		SELECT id, username, email, gold_balance, status, win_rate, created_at, last_login_at
		FROM players
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("playerRepo.FindAll: %w", err)
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var p models.Player
		if err := rows.Scan(
			&p.ID, &p.Username, &p.Email,
			&p.GoldBalance, &p.Status, &p.WinRate,
			&p.CreatedAt, &p.LastLoginAt,
		); err != nil {
			return nil, fmt.Errorf("playerRepo.FindAll scan: %w", err)
		}
		players = append(players, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("playerRepo.FindAll rows: %w", err)
	}
	return players, nil
}

// FindByID returns a single player by primary key.
func (r *playerSQLServerRepo) FindByID(ctx context.Context, id string) (*models.Player, error) {
	const query = `
		SELECT id, username, email, gold_balance, status, win_rate, created_at, last_login_at
		FROM players
		WHERE id = @p1`

	var p models.Player
	err := r.db.QueryRowContext(ctx, query, sql.Named("p1", id)).Scan(
		&p.ID, &p.Username, &p.Email,
		&p.GoldBalance, &p.Status, &p.WinRate,
		&p.CreatedAt, &p.LastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("playerRepo.FindByID: player %q not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("playerRepo.FindByID: %w", err)
	}
	return &p, nil
}

// Update applies a partial update to a player record (only non-nil fields are updated).
func (r *playerSQLServerRepo) Update(ctx context.Context, id string, req models.UpdatePlayerRequest) error {
	setClauses := []string{}
	args := []interface{}{}
	paramIdx := 1

	if req.GoldBalance != nil {
		setClauses = append(setClauses, fmt.Sprintf("gold_balance = @p%d", paramIdx))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), *req.GoldBalance))
		paramIdx++
	}
	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = @p%d", paramIdx))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), *req.Status))
		paramIdx++
	}
	if req.WinRate != nil {
		setClauses = append(setClauses, fmt.Sprintf("win_rate = @p%d", paramIdx))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), *req.WinRate))
		paramIdx++
	}

	if len(setClauses) == 0 {
		return nil // nothing to update
	}

	args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), id))
	query := fmt.Sprintf(
		"UPDATE players SET %s WHERE id = @p%d",
		strings.Join(setClauses, ", "), paramIdx,
	)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("playerRepo.Update: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("playerRepo.Update: player %q not found", id)
	}
	return nil
}

// Delete removes a player record by ID (soft-delete: sets status to 'banned').
func (r *playerSQLServerRepo) Delete(ctx context.Context, id string) error {
	const query = `UPDATE players SET status = 'banned' WHERE id = @p1`
	result, err := r.db.ExecContext(ctx, query, sql.Named("p1", id))
	if err != nil {
		return fmt.Errorf("playerRepo.Delete: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("playerRepo.Delete: player %q not found", id)
	}
	return nil
}
