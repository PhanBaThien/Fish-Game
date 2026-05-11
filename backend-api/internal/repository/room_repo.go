package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/yourname/fish-game-backend/internal/models"
	"github.com/yourname/fish-game-backend/pkg/database"
)

// RoomRepository defines the interface for game room data access.
type RoomRepository interface {
	FindAll(ctx context.Context) ([]models.Room, error)
	FindByID(ctx context.Context, id string) (*models.Room, error)
	Create(ctx context.Context, req models.CreateRoomRequest) (*models.Room, error)
	Update(ctx context.Context, id string, req models.UpdateRoomRequest) error
	Close(ctx context.Context, id string) error
}

type roomSQLServerRepo struct {
	db *database.DB
}

// NewRoomRepository creates a new SQL Server–backed RoomRepository.
func NewRoomRepository(db *database.DB) RoomRepository {
	return &roomSQLServerRepo{db: db}
}

func (r *roomSQLServerRepo) FindAll(ctx context.Context) ([]models.Room, error) {
	const query = `
		SELECT id, name, type, bet_amount, players, max_players, status, base_rtp
		FROM rooms
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("roomRepo.FindAll: %w", err)
	}
	defer rows.Close()

	var list []models.Room
	for rows.Next() {
		var rm models.Room
		if err := rows.Scan(
			&rm.ID, &rm.Name, &rm.Type, &rm.BetAmount,
			&rm.Players, &rm.MaxPlayers, &rm.Status, &rm.BaseRTP,
		); err != nil {
			return nil, fmt.Errorf("roomRepo.FindAll scan: %w", err)
		}
		list = append(list, rm)
	}
	return list, rows.Err()
}

func (r *roomSQLServerRepo) FindByID(ctx context.Context, id string) (*models.Room, error) {
	const query = `
		SELECT id, name, type, bet_amount, players, max_players, status, base_rtp
		FROM rooms WHERE id = @p1`

	var rm models.Room
	err := r.db.QueryRowContext(ctx, query, sql.Named("p1", id)).Scan(
		&rm.ID, &rm.Name, &rm.Type, &rm.BetAmount,
		&rm.Players, &rm.MaxPlayers, &rm.Status, &rm.BaseRTP,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("roomRepo.FindByID: room %q not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("roomRepo.FindByID: %w", err)
	}
	return &rm, nil
}

func (r *roomSQLServerRepo) Create(ctx context.Context, req models.CreateRoomRequest) (*models.Room, error) {
	const query = `
		INSERT INTO rooms (name, type, bet_amount, players, max_players, status, base_rtp)
		OUTPUT INSERTED.id, INSERTED.name, INSERTED.type, INSERTED.bet_amount,
		       INSERTED.players, INSERTED.max_players, INSERTED.status, INSERTED.base_rtp
		VALUES (@p1, @p2, @p3, 0, @p4, 'waiting', @p5)`

	var rm models.Room
	err := r.db.QueryRowContext(ctx, query,
		sql.Named("p1", req.Name),
		sql.Named("p2", req.Type),
		sql.Named("p3", req.BetAmount),
		sql.Named("p4", req.MaxPlayers),
		sql.Named("p5", req.BaseRTP),
	).Scan(&rm.ID, &rm.Name, &rm.Type, &rm.BetAmount, &rm.Players, &rm.MaxPlayers, &rm.Status, &rm.BaseRTP)
	if err != nil {
		return nil, fmt.Errorf("roomRepo.Create: %w", err)
	}
	return &rm, nil
}

func (r *roomSQLServerRepo) Update(ctx context.Context, id string, req models.UpdateRoomRequest) error {
	setClauses := []string{}
	args := []interface{}{}
	paramIdx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = @p%d", paramIdx))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), *req.Name))
		paramIdx++
	}
	if req.BaseRTP != nil {
		setClauses = append(setClauses, fmt.Sprintf("base_rtp = @p%d", paramIdx))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), *req.BaseRTP))
		paramIdx++
	}
	if req.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = @p%d", paramIdx))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), *req.Status))
		paramIdx++
	}
	if req.BetAmount != nil {
		setClauses = append(setClauses, fmt.Sprintf("bet_amount = @p%d", paramIdx))
		args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), *req.BetAmount))
		paramIdx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, sql.Named(fmt.Sprintf("p%d", paramIdx), id))
	query := fmt.Sprintf(
		"UPDATE rooms SET %s WHERE id = @p%d",
		strings.Join(setClauses, ", "), paramIdx,
	)
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("roomRepo.Update: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("roomRepo.Update: room %q not found", id)
	}
	return nil
}

func (r *roomSQLServerRepo) Close(ctx context.Context, id string) error {
	const query = `UPDATE rooms SET status = 'closed' WHERE id = @p1`
	result, err := r.db.ExecContext(ctx, query, sql.Named("p1", id))
	if err != nil {
		return fmt.Errorf("roomRepo.Close: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("roomRepo.Close: room %q not found", id)
	}
	return nil
}
