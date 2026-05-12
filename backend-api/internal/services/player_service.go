package services

import (
	"context"
	"fmt"

	"github.com/yourname/fish-game-backend/internal/models"
	"github.com/yourname/fish-game-backend/internal/repository"
)

// PlayerService defines business logic for player management.
type PlayerService interface {
	ListPlayers(ctx context.Context) ([]models.Player, error)
	GetPlayerByID(ctx context.Context, id string) (*models.Player, error)
	UpdatePlayer(ctx context.Context, id string, req models.UpdatePlayerRequest) error
	BanPlayer(ctx context.Context, id string) error
}

type playerService struct {
	repo repository.PlayerRepository
}

// NewPlayerService creates a PlayerService with the given repository.
func NewPlayerService(repo repository.PlayerRepository) PlayerService {
	return &playerService{repo: repo}
}

func (s *playerService) ListPlayers(ctx context.Context) ([]models.Player, error) {
	return s.repo.FindAll(ctx)
}

func (s *playerService) GetPlayerByID(ctx context.Context, id string) (*models.Player, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("player not found: %w", err)
	}
	return p, nil
}

func (s *playerService) UpdatePlayer(ctx context.Context, id string, req models.UpdatePlayerRequest) error {
	return s.repo.Update(ctx, id, req)
}

func (s *playerService) BanPlayer(ctx context.Context, id string) error {
	status := "banned"
	return s.repo.Update(ctx, id, models.UpdatePlayerRequest{Status: &status})
}
