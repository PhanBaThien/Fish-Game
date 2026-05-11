package services

import (
	"context"
	"fmt"

	"github.com/yourname/fish-game-backend/internal/models"
	"github.com/yourname/fish-game-backend/internal/repository"
)

// FishService defines business logic for fish configuration management.
type FishService interface {
	ListFish(ctx context.Context) ([]models.Fish, error)
	GetFishByID(ctx context.Context, id string) (*models.Fish, error)
	CreateFish(ctx context.Context, req models.CreateFishRequest) (*models.Fish, error)
	UpdateFish(ctx context.Context, id string, req models.UpdateFishRequest) error
	DeleteFish(ctx context.Context, id string) error
}

type fishService struct {
	repo repository.FishRepository
}

// NewFishService creates a FishService with the given repository.
func NewFishService(repo repository.FishRepository) FishService {
	return &fishService{repo: repo}
}

func (s *fishService) ListFish(ctx context.Context) ([]models.Fish, error) {
	return s.repo.FindAll(ctx)
}

func (s *fishService) GetFishByID(ctx context.Context, id string) (*models.Fish, error) {
	f, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fish not found: %w", err)
	}
	return f, nil
}

func (s *fishService) CreateFish(ctx context.Context, req models.CreateFishRequest) (*models.Fish, error) {
	return s.repo.Create(ctx, req)
}

func (s *fishService) UpdateFish(ctx context.Context, id string, req models.UpdateFishRequest) error {
	return s.repo.Update(ctx, id, req)
}

func (s *fishService) DeleteFish(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
