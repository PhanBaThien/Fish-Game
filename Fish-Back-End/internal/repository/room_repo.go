package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
)

var ErrRoomNotFound = errors.New("không tìm thấy phòng chơi")

// RoomRepository định nghĩa interface thao tác bảng rooms
type RoomRepository interface {
	List(ctx context.Context, status, roomType string, limit, offset int) ([]models.Room, int64, error)
	GetByID(ctx context.Context, id string) (*models.Room, error)
	Create(ctx context.Context, r *models.Room) error
	Update(ctx context.Context, r *models.Room) error
	UpdateRTP(ctx context.Context, id string, baseRTP float64) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, q string, limit int) ([]models.Room, error)
}

type roomPgRepo struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomPgRepo{db: db}
}

func (r *roomPgRepo) List(ctx context.Context, status, roomType string, limit, offset int) ([]models.Room, int64, error) {
	query := `
		SELECT id, name, type, bet_amount, max_players, status, base_rtp, created_at, updated_at
		FROM rooms
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR type = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, status, roomType, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("roomRepo.List: %w", err)
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var rm models.Room
		if err := rows.Scan(&rm.ID, &rm.Name, &rm.Type, &rm.BetAmount, &rm.MaxPlayers, &rm.Status, &rm.BaseRTP, &rm.CreatedAt, &rm.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("roomRepo.List scan: %w", err)
		}
		rooms = append(rooms, rm)
	}

	var total int64
	countQ := `SELECT COUNT(*) FROM rooms WHERE ($1 = '' OR status = $1) AND ($2 = '' OR type = $2)`
	if err := r.db.QueryRowContext(ctx, countQ, status, roomType).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("roomRepo.List count: %w", err)
	}

	return rooms, total, nil
}

func (r *roomPgRepo) GetByID(ctx context.Context, id string) (*models.Room, error) {
	query := `
		SELECT id, name, type, bet_amount, max_players, status, base_rtp, created_at, updated_at
		FROM rooms WHERE id = $1`

	var rm models.Room
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&rm.ID, &rm.Name, &rm.Type, &rm.BetAmount, &rm.MaxPlayers, &rm.Status, &rm.BaseRTP, &rm.CreatedAt, &rm.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoomNotFound
		}
		return nil, fmt.Errorf("roomRepo.GetByID: %w", err)
	}
	return &rm, nil
}

func (r *roomPgRepo) Create(ctx context.Context, rm *models.Room) error {
	query := `
		INSERT INTO rooms (name, type, bet_amount, max_players, status, base_rtp)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		rm.Name, rm.Type, rm.BetAmount, rm.MaxPlayers, rm.Status, rm.BaseRTP,
	).Scan(&rm.ID, &rm.CreatedAt, &rm.UpdatedAt)
}

func (r *roomPgRepo) Update(ctx context.Context, rm *models.Room) error {
	query := `
		UPDATE rooms
		SET name=$2, type=$3, bet_amount=$4, max_players=$5, status=$6, base_rtp=$7, updated_at=NOW()
		WHERE id=$1
		RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query,
		rm.ID, rm.Name, rm.Type, rm.BetAmount, rm.MaxPlayers, rm.Status, rm.BaseRTP,
	).Scan(&rm.UpdatedAt)
}

func (r *roomPgRepo) UpdateRTP(ctx context.Context, id string, baseRTP float64) error {
	res, err := r.db.ExecContext(ctx, `UPDATE rooms SET base_rtp=$2, updated_at=NOW() WHERE id=$1`, id, baseRTP)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrRoomNotFound
	}
	return nil
}

func (r *roomPgRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM rooms WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrRoomNotFound
	}
	return nil
}

func (r *roomPgRepo) Search(ctx context.Context, q string, limit int) ([]models.Room, error) {
	query := `
		SELECT id, name, type, bet_amount, max_players, status, base_rtp, created_at, updated_at
		FROM rooms
		WHERE name ILIKE '%' || $1 || '%'
		ORDER BY name
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, q, limit)
	if err != nil {
		return nil, fmt.Errorf("roomRepo.Search: %w", err)
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var rm models.Room
		if err := rows.Scan(&rm.ID, &rm.Name, &rm.Type, &rm.BetAmount, &rm.MaxPlayers, &rm.Status, &rm.BaseRTP, &rm.CreatedAt, &rm.UpdatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, rm)
	}
	return rooms, nil
}
