package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
)

var ErrFishNotFound = errors.New("không tìm thấy cấu hình cá")

// FishRepository định nghĩa interface thao tác bảng fish
type FishRepository interface {
	List(ctx context.Context, role string, isActive *bool, limit, offset int) ([]models.Fish, int64, error)
	GetByID(ctx context.Context, id string) (*models.Fish, error)
	Create(ctx context.Context, f *models.Fish) error
	Update(ctx context.Context, f *models.Fish) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, q string, limit int) ([]models.Fish, error)
}

type fishPgRepo struct {
	db *sql.DB
}

func NewFishRepository(db *sql.DB) FishRepository {
	return &fishPgRepo{db: db}
}

func (r *fishPgRepo) List(ctx context.Context, role string, isActive *bool, limit, offset int) ([]models.Fish, int64, error) {
	query := `
		SELECT id, name, multiplier, base_prob, speed, role, is_active, updated_at
		FROM fish
		WHERE ($1 = '' OR role = $1)
		  AND ($2::boolean IS NULL OR is_active = $2)
		ORDER BY id
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, role, isActive, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fishRepo.List: %w", err)
	}
	defer rows.Close()

	var fishes []models.Fish
	for rows.Next() {
		var f models.Fish
		if err := rows.Scan(&f.ID, &f.Name, &f.Multiplier, &f.BaseProb, &f.Speed, &f.Role, &f.IsActive, &f.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("fishRepo.List scan: %w", err)
		}
		fishes = append(fishes, f)
	}

	var total int64
	countQ := `SELECT COUNT(*) FROM fish WHERE ($1 = '' OR role = $1) AND ($2::boolean IS NULL OR is_active = $2)`
	if err := r.db.QueryRowContext(ctx, countQ, role, isActive).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("fishRepo.List count: %w", err)
	}

	return fishes, total, nil
}

func (r *fishPgRepo) GetByID(ctx context.Context, id string) (*models.Fish, error) {
	query := `SELECT id, name, multiplier, base_prob, speed, role, is_active, updated_at FROM fish WHERE id = $1`

	var f models.Fish
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&f.ID, &f.Name, &f.Multiplier, &f.BaseProb, &f.Speed, &f.Role, &f.IsActive, &f.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFishNotFound
		}
		return nil, fmt.Errorf("fishRepo.GetByID: %w", err)
	}
	return &f, nil
}

func (r *fishPgRepo) Create(ctx context.Context, f *models.Fish) error {
	query := `
		INSERT INTO fish (id, name, multiplier, base_prob, speed, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query,
		f.ID, f.Name, f.Multiplier, f.BaseProb, f.Speed, f.Role, f.IsActive,
	).Scan(&f.UpdatedAt)
}

func (r *fishPgRepo) Update(ctx context.Context, f *models.Fish) error {
	query := `
		UPDATE fish
		SET name=$2, multiplier=$3, base_prob=$4, speed=$5, role=$6, is_active=$7, updated_at=NOW()
		WHERE id=$1
		RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query,
		f.ID, f.Name, f.Multiplier, f.BaseProb, f.Speed, f.Role, f.IsActive,
	).Scan(&f.UpdatedAt)
}

func (r *fishPgRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM fish WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrFishNotFound
	}
	return nil
}

func (r *fishPgRepo) Search(ctx context.Context, q string, limit int) ([]models.Fish, error) {
	query := `
		SELECT id, name, multiplier, base_prob, speed, role, is_active, updated_at
		FROM fish
		WHERE name ILIKE '%' || $1 || '%' OR id ILIKE '%' || $1 || '%'
		ORDER BY id
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, q, limit)
	if err != nil {
		return nil, fmt.Errorf("fishRepo.Search: %w", err)
	}
	defer rows.Close()

	var fishes []models.Fish
	for rows.Next() {
		var f models.Fish
		var updatedAt time.Time
		if err := rows.Scan(&f.ID, &f.Name, &f.Multiplier, &f.BaseProb, &f.Speed, &f.Role, &f.IsActive, &updatedAt); err != nil {
			return nil, err
		}
		f.UpdatedAt = updatedAt
		fishes = append(fishes, f)
	}
	return fishes, nil
}
