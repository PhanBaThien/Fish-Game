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

type RoomRepository interface {
	List(ctx context.Context) ([]models.Room, error)
	GetByID(ctx context.Context, id int64) (*models.Room, error)
	Create(ctx context.Context, room *models.Room) error
	Update(ctx context.Context, room *models.Room) error
	Delete(ctx context.Context, id int64) error
}

type roomPgRepo struct {
	queries *dbgen.Queries
}

func NewRoomRepository(pool *pgxpool.Pool) RoomRepository {
	return &roomPgRepo{queries: dbgen.New(pool)}
}

func mapToModelRoom(r dbgen.Room) models.Room {
	return models.Room{
		ID:          r.ID,
		Name:        r.Name,
		MinBet:      r.MinBet,
		MaxPlayers:  r.MaxPlayers,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func (r *roomPgRepo) List(ctx context.Context) ([]models.Room, error) {
	rows, err := r.queries.ListRooms(ctx)
	if err != nil {
		return nil, apperror.Wrap("repository", "roomRepo.List", err)
	}
	rooms := make([]models.Room, len(rows))
	for i, row := range rows {
		rooms[i] = mapToModelRoom(row)
	}
	return rooms, nil
}

func (r *roomPgRepo) GetByID(ctx context.Context, id int64) (*models.Room, error) {
	row, err := r.queries.GetRoomByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrRoomNotFound
		}
		return nil, apperror.Wrap("repository", "roomRepo.GetByID", err)
	}
	room := mapToModelRoom(row)
	return &room, nil
}

func (r *roomPgRepo) Create(ctx context.Context, room *models.Room) error {
	row, err := r.queries.CreateRoom(ctx, dbgen.CreateRoomParams{
		Name:        room.Name,
		MinBet:      room.MinBet,
		MaxPlayers:  room.MaxPlayers,
		Description: room.Description,
	})
	if err != nil {
		return apperror.Wrap("repository", "roomRepo.Create", err)
	}
	*room = mapToModelRoom(row)
	return nil
}

func (r *roomPgRepo) Update(ctx context.Context, room *models.Room) error {
	row, err := r.queries.UpdateRoom(ctx, dbgen.UpdateRoomParams{
		ID:          room.ID,
		Name:        room.Name,
		MinBet:      room.MinBet,
		MaxPlayers:  room.MaxPlayers,
		Description: room.Description,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrRoomNotFound
		}
		return apperror.Wrap("repository", "roomRepo.Update", err)
	}
	*room = mapToModelRoom(row)
	return nil
}

func (r *roomPgRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.queries.DeleteRoom(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrRoomNotFound
		}
		return apperror.Wrap("repository", "roomRepo.Delete", err)
	}
	return nil
}
