package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
)

var ErrSettingNotFound = errors.New("không tìm thấy cài đặt")

// SettingRepository định nghĩa interface thao tác bảng settings
type SettingRepository interface {
	GetAll(ctx context.Context) ([]models.Setting, error)
	GetByKey(ctx context.Context, key string) (*models.Setting, error)
	Upsert(ctx context.Context, s *models.Setting) error
}

type settingPgRepo struct {
	db *sql.DB
}

func NewSettingRepository(db *sql.DB) SettingRepository {
	return &settingPgRepo{db: db}
}

func (r *settingPgRepo) GetAll(ctx context.Context) ([]models.Setting, error) {
	query := `SELECT key, value, updated_by, created_at, updated_at FROM settings ORDER BY key`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("settingRepo.GetAll: %w", err)
	}
	defer rows.Close()

	var settings []models.Setting
	for rows.Next() {
		var s models.Setting
		var valRaw []byte
		if err := rows.Scan(&s.Key, &valRaw, &s.UpdatedBy, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("settingRepo.GetAll scan: %w", err)
		}
		s.Value = json.RawMessage(valRaw)
		settings = append(settings, s)
	}
	return settings, nil
}

func (r *settingPgRepo) GetByKey(ctx context.Context, key string) (*models.Setting, error) {
	query := `SELECT key, value, updated_by, created_at, updated_at FROM settings WHERE key = $1`

	var s models.Setting
	var valRaw []byte
	err := r.db.QueryRowContext(ctx, query, key).
		Scan(&s.Key, &valRaw, &s.UpdatedBy, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSettingNotFound
		}
		return nil, fmt.Errorf("settingRepo.GetByKey: %w", err)
	}
	s.Value = json.RawMessage(valRaw)
	return &s, nil
}

func (r *settingPgRepo) Upsert(ctx context.Context, s *models.Setting) error {
	query := `
		INSERT INTO settings (key, value, updated_by, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (key) DO UPDATE
		SET value=$2, updated_by=$3, updated_at=NOW()
		RETURNING created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, s.Key, s.Value, s.UpdatedBy).
		Scan(&s.CreatedAt, &s.UpdatedAt)
}
