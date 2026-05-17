package usecase

import (
	"context"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomUsecase interface {
	List(ctx context.Context) ([]models.Room, error)
	GetByID(ctx context.Context, id int64) (*models.Room, error)
	Create(ctx context.Context, req *domain.CreateRoomRequest) (*models.Room, error)
	Update(ctx context.Context, id int64, req *domain.UpdateRoomRequest) (*models.Room, error)
	Delete(ctx context.Context, id int64) error
}

type roomUsecase struct {
	roomRepo repository.RoomRepository
}

func NewRoomUsecase(repo repository.RoomRepository) RoomUsecase {
	return &roomUsecase{roomRepo: repo}
}

func (u *roomUsecase) List(ctx context.Context) ([]models.Room, error) {
	return u.roomRepo.List(ctx)
}

func (u *roomUsecase) GetByID(ctx context.Context, id int64) (*models.Room, error) {
	return u.roomRepo.GetByID(ctx, id)
}

func (u *roomUsecase) Create(ctx context.Context, req *domain.CreateRoomRequest) (*models.Room, error) {
	room := &models.Room{
		Name:       req.Name,
		MinBet:     req.MinBet,
		MaxPlayers: req.MaxPlayers,
	}
	if req.Description != nil {
		room.Description = pgtype.Text{String: *req.Description, Valid: true}
	}
	if err := u.roomRepo.Create(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}

func (u *roomUsecase) Update(ctx context.Context, id int64, req *domain.UpdateRoomRequest) (*models.Room, error) {
	room, err := u.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		room.Name = *req.Name
	}
	if req.MinBet != nil {
		room.MinBet = *req.MinBet
	}
	if req.MaxPlayers != nil {
		room.MaxPlayers = *req.MaxPlayers
	}
	if req.Description != nil {
		room.Description = pgtype.Text{String: *req.Description, Valid: true}
	}
	if err := u.roomRepo.Update(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}

func (u *roomUsecase) Delete(ctx context.Context, id int64) error {
	return u.roomRepo.Delete(ctx, id)
}
