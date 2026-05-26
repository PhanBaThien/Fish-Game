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

type FishRepository interface {
	List(ctx context.Context) ([]models.Fish, error)
	GetByID(ctx context.Context, id int32) (*models.Fish, error)
	Create(ctx context.Context, fish *models.Fish) error
	Update(ctx context.Context, fish *models.Fish) error
	Delete(ctx context.Context, id int32) error
}

type fishPgRepo struct {
	queries *dbgen.Queries
}

func NewFishRepository(pool *pgxpool.Pool) FishRepository {
	return &fishPgRepo{queries: dbgen.New(pool)}
}

func mapToModelFish(f dbgen.Fish) models.Fish {
	return models.Fish{
		ID:               f.ID,
		Name:             f.Name,
		Health:           f.Health,
		RewardMultiplier: f.RewardMultiplier,
		Speed:            f.Speed,
		AssetPath:        f.AssetPath,
		CreatedAt:        f.CreatedAt.Time,
		UpdatedAt:        f.UpdatedAt.Time,
	}
}

func (r *fishPgRepo) List(ctx context.Context) ([]models.Fish, error) {
	rows, err := r.queries.ListFishes(ctx)
	if err != nil {
		return nil, apperror.Wrap("repository", "fishRepo.List", err)
	}
	fishes := make([]models.Fish, len(rows))
	for i, row := range rows {
		fishes[i] = mapToModelFish(row)
	}
	return fishes, nil
}

func (r *fishPgRepo) GetByID(ctx context.Context, id int32) (*models.Fish, error) {
	row, err := r.queries.GetFishByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrFishNotFound
		}
		return nil, apperror.Wrap("repository", "fishRepo.GetByID", err)
	}
	fish := mapToModelFish(row)
	return &fish, nil
}

func (r *fishPgRepo) Create(ctx context.Context, fish *models.Fish) error {
	row, err := r.queries.CreateFish(ctx, dbgen.CreateFishParams{
		Name:             fish.Name,
		Health:           fish.Health,
		RewardMultiplier: fish.RewardMultiplier,
		Speed:            fish.Speed,
		AssetPath:        fish.AssetPath,
	})
	if err != nil {
		return apperror.Wrap("repository", "fishRepo.Create", err)
	}
	*fish = mapToModelFish(row)
	return nil
}

func (r *fishPgRepo) Update(ctx context.Context, fish *models.Fish) error {
	row, err := r.queries.UpdateFish(ctx, dbgen.UpdateFishParams{
		ID:               fish.ID,
		Name:             fish.Name,
		Health:           fish.Health,
		RewardMultiplier: fish.RewardMultiplier,
		Speed:            fish.Speed,
		AssetPath:        fish.AssetPath,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrFishNotFound
		}
		return apperror.Wrap("repository", "fishRepo.Update", err)
	}
	*fish = mapToModelFish(row)
	return nil
}

func (r *fishPgRepo) Delete(ctx context.Context, id int32) error {
	_, err := r.queries.DeleteFish(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrFishNotFound
		}
		return apperror.Wrap("repository", "fishRepo.Delete", err)
	}
	return nil
}
