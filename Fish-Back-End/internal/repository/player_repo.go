package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
)

var ErrPlayerNotFound = errors.New("không tìm thấy người chơi")

// PlayerRepository định nghĩa interface thao tác bảng players
type PlayerRepository interface {
	List(ctx context.Context, status string, limit, offset int) ([]models.Player, int64, error)
	GetByID(ctx context.Context, id string) (*models.Player, error)
	GetByUsername(ctx context.Context, username string) (*models.Player, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	Create(ctx context.Context, p *models.Player) error
	Update(ctx context.Context, p *models.Player) error
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateWinRate(ctx context.Context, id string, winRate float64) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, q string, limit int) ([]models.Player, error)
}

type playerPgRepo struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) PlayerRepository {
	return &playerPgRepo{db: db}
}

func (r *playerPgRepo) List(ctx context.Context, status string, limit, offset int) ([]models.Player, int64, error) {
	query := `
		SELECT id, username, email, gold_balance, status, win_rate, created_at, last_login_at
		FROM players
		WHERE ($1 = '' OR status = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("playerRepo.List: %w", err)
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var p models.Player
		if err := rows.Scan(&p.ID, &p.Username, &p.Email, &p.GoldBalance, &p.Status, &p.WinRate, &p.CreatedAt, &p.LastLoginAt); err != nil {
			return nil, 0, fmt.Errorf("playerRepo.List scan: %w", err)
		}
		players = append(players, p)
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM players WHERE ($1 = '' OR status = $1)`
	if err := r.db.QueryRowContext(ctx, countQuery, status).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("playerRepo.List count: %w", err)
	}

	return players, total, nil
}

func (r *playerPgRepo) GetByID(ctx context.Context, id string) (*models.Player, error) {
	query := `
		SELECT id, username, email, gold_balance, status, win_rate, created_at, last_login_at
		FROM players WHERE id = $1`

	var p models.Player
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&p.ID, &p.Username, &p.Email, &p.GoldBalance, &p.Status, &p.WinRate, &p.CreatedAt, &p.LastLoginAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlayerNotFound
		}
		return nil, fmt.Errorf("playerRepo.GetByID: %w", err)
	}
	return &p, nil
}

func (r *playerPgRepo) GetByUsername(ctx context.Context, username string) (*models.Player, error) {
	query := `
		SELECT id, username, email, gold_balance, status, win_rate, created_at, last_login_at
		FROM players WHERE username = $1`

	var p models.Player
	err := r.db.QueryRowContext(ctx, query, username).
		Scan(&p.ID, &p.Username, &p.Email, &p.GoldBalance, &p.Status, &p.WinRate, &p.CreatedAt, &p.LastLoginAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlayerNotFound
		}
		return nil, fmt.Errorf("playerRepo.GetByUsername: %w", err)
	}
	return &p, nil
}

func (r *playerPgRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM players WHERE username=$1)`, username).Scan(&exists)
	return exists, err
}

func (r *playerPgRepo) Create(ctx context.Context, p *models.Player) error {
	query := `
		INSERT INTO players (username, email, password_hash, gold_balance, status, win_rate)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query,
		p.Username, p.Email, p.PasswordHash, p.GoldBalance, p.Status, p.WinRate,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *playerPgRepo) Update(ctx context.Context, p *models.Player) error {
	query := `
		UPDATE players SET username=$2, email=$3, status=$4
		WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Username, p.Email, p.Status)
	return err
}

func (r *playerPgRepo) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE players SET status=$2 WHERE id=$1`, id, status)
	return err
}

func (r *playerPgRepo) UpdateWinRate(ctx context.Context, id string, winRate float64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE players SET win_rate=$2 WHERE id=$1`, id, winRate)
	return err
}

func (r *playerPgRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM players WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrPlayerNotFound
	}
	return nil
}

func (r *playerPgRepo) Search(ctx context.Context, q string, limit int) ([]models.Player, error) {
	query := `
		SELECT id, username, email, gold_balance, status, win_rate, created_at, last_login_at
		FROM players
		WHERE username ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%'
		ORDER BY username
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, q, limit)
	if err != nil {
		return nil, fmt.Errorf("playerRepo.Search: %w", err)
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var p models.Player
		var lastLogin *time.Time
		if err := rows.Scan(&p.ID, &p.Username, &p.Email, &p.GoldBalance, &p.Status, &p.WinRate, &p.CreatedAt, &lastLogin); err != nil {
			return nil, err
		}
		p.LastLoginAt = lastLogin
		players = append(players, p)
	}
	return players, nil
}
