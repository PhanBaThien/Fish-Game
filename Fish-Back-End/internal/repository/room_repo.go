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

// pgtextToPtr chuyển pgtype.Text (nullable DB) → *string (Go chuẩn)
func pgtextToPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

// ptrToPgtext chuyển *string (Go chuẩn) → pgtype.Text (để ghi xuống DB)
func ptrToPgtext(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func mapToModelRoom(r dbgen.Room) models.Room {
	return models.Room{
		ID:          r.ID,
		Name:        r.Name,
		MaxPlayers:  r.MaxPlayers,
		Description: pgtextToPtr(r.Description),
		RTP:         r.Rtp,
		CreatedAt:   r.CreatedAt.Time,
		UpdatedAt:   r.UpdatedAt.Time,
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
		MaxPlayers:  room.MaxPlayers,
		Description: ptrToPgtext(room.Description),
		Rtp:         room.RTP,
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
		MaxPlayers:  room.MaxPlayers,
		Description: ptrToPgtext(room.Description),
		Rtp:         room.RTP,
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
