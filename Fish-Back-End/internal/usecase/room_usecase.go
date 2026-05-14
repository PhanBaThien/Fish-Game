package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository"
)

var ErrRoomNotFound = errors.New("không tìm thấy phòng chơi")

// RoomUsecase định nghĩa nghiệp vụ quản lý phòng chơi
type RoomUsecase interface {
	ListRooms(ctx context.Context, req *domain.ListRoomsRequest) (*domain.ListRoomsResponse, error)
	GetRoom(ctx context.Context, id string) (*models.Room, error)
	CreateRoom(ctx context.Context, req *domain.CreateRoomRequest) (*models.Room, error)
	UpdateRoom(ctx context.Context, id string, req *domain.UpdateRoomRequest) (*models.Room, error)
	DeleteRoom(ctx context.Context, id string) error
	UpdateRoomRTP(ctx context.Context, id string, req *domain.UpdateRoomRTPRequest) error
}

type roomUsecase struct {
	roomRepo repository.RoomRepository
}

func NewRoomUsecase(roomRepo repository.RoomRepository) RoomUsecase {
	return &roomUsecase{roomRepo: roomRepo}
}

func (u *roomUsecase) ListRooms(ctx context.Context, req *domain.ListRoomsRequest) (*domain.ListRoomsResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}
	offset := (req.Page - 1) * req.Limit

	rooms, total, err := u.roomRepo.List(ctx, req.Status, req.Type, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("roomUsecase.ListRooms: %w", err)
	}

	return &domain.ListRoomsResponse{
		Items: rooms,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

func (u *roomUsecase) GetRoom(ctx context.Context, id string) (*models.Room, error) {
	r, err := u.roomRepo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrRoomNotFound) {
		return nil, ErrRoomNotFound
	}
	return r, err
}

func (u *roomUsecase) CreateRoom(ctx context.Context, req *domain.CreateRoomRequest) (*models.Room, error) {
	r := &models.Room{
		Name:       req.Name,
		Type:       req.Type,
		BetAmount:  req.BetAmount,
		MaxPlayers: req.MaxPlayers,
		Status:     "waiting",
		BaseRTP:    req.BaseRTP,
	}

	if err := u.roomRepo.Create(ctx, r); err != nil {
		return nil, fmt.Errorf("roomUsecase.CreateRoom: %w", err)
	}
	return r, nil
}

func (u *roomUsecase) UpdateRoom(ctx context.Context, id string, req *domain.UpdateRoomRequest) (*models.Room, error) {
	r, err := u.roomRepo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrRoomNotFound) {
		return nil, ErrRoomNotFound
	}
	if err != nil {
		return nil, err
	}

	r.Name = req.Name
	r.Type = req.Type
	r.BetAmount = req.BetAmount
	r.MaxPlayers = req.MaxPlayers
	r.Status = req.Status
	r.BaseRTP = req.BaseRTP

	if err := u.roomRepo.Update(ctx, r); err != nil {
		return nil, fmt.Errorf("roomUsecase.UpdateRoom: %w", err)
	}
	return r, nil
}

func (u *roomUsecase) DeleteRoom(ctx context.Context, id string) error {
	err := u.roomRepo.Delete(ctx, id)
	if errors.Is(err, repository.ErrRoomNotFound) {
		return ErrRoomNotFound
	}
	return err
}

func (u *roomUsecase) UpdateRoomRTP(ctx context.Context, id string, req *domain.UpdateRoomRTPRequest) error {
	err := u.roomRepo.UpdateRTP(ctx, id, req.BaseRTP)
	if errors.Is(err, repository.ErrRoomNotFound) {
		return ErrRoomNotFound
	}
	return err
}
