package services

import (
	"context"
	"fmt"

	"github.com/yourname/fish-game-backend/internal/models"
	"github.com/yourname/fish-game-backend/internal/repository"
)

// RoomService defines business logic for game room management.
type RoomService interface {
	ListRooms(ctx context.Context) ([]models.Room, error)
	GetRoomByID(ctx context.Context, id string) (*models.Room, error)
	CreateRoom(ctx context.Context, req models.CreateRoomRequest) (*models.Room, error)
	UpdateRoom(ctx context.Context, id string, req models.UpdateRoomRequest) error
	CloseRoom(ctx context.Context, id string) error
}

type roomService struct {
	repo repository.RoomRepository
}

// NewRoomService creates a RoomService with the given repository.
func NewRoomService(repo repository.RoomRepository) RoomService {
	return &roomService{repo: repo}
}

func (s *roomService) ListRooms(ctx context.Context) ([]models.Room, error) {
	return s.repo.FindAll(ctx)
}

func (s *roomService) GetRoomByID(ctx context.Context, id string) (*models.Room, error) {
	r, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}
	return r, nil
}

func (s *roomService) CreateRoom(ctx context.Context, req models.CreateRoomRequest) (*models.Room, error) {
	return s.repo.Create(ctx, req)
}

func (s *roomService) UpdateRoom(ctx context.Context, id string, req models.UpdateRoomRequest) error {
	return s.repo.Update(ctx, id, req)
}

func (s *roomService) CloseRoom(ctx context.Context, id string) error {
	return s.repo.Close(ctx, id)
}
