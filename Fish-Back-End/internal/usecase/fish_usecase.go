package usecase

import (
	"context"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository"
)

type FishUsecase interface {
	List(ctx context.Context) ([]models.Fish, error)
	GetByID(ctx context.Context, id int32) (*models.Fish, error)
	Create(ctx context.Context, req *domain.CreateFishRequest) (*models.Fish, error)
	Update(ctx context.Context, id int32, req *domain.UpdateFishRequest) (*models.Fish, error)
	Delete(ctx context.Context, id int32) error
}

type fishUsecase struct {
	fishRepo repository.FishRepository
}

func NewFishUsecase(repo repository.FishRepository) FishUsecase {
	return &fishUsecase{fishRepo: repo}
}

func (u *fishUsecase) List(ctx context.Context) ([]models.Fish, error) {
	return u.fishRepo.List(ctx)
}

func (u *fishUsecase) GetByID(ctx context.Context, id int32) (*models.Fish, error) {
	return u.fishRepo.GetByID(ctx, id)
}

func (u *fishUsecase) Create(ctx context.Context, req *domain.CreateFishRequest) (*models.Fish, error) {
	fish := &models.Fish{
		Name:             req.Name,
		Health:           req.Health,
		RewardMultiplier: req.RewardMultiplier,
		Speed:            req.Speed,
		AssetPath:        req.AssetPath,
	}
	if err := u.fishRepo.Create(ctx, fish); err != nil {
		return nil, err
	}
	return fish, nil
}

func (u *fishUsecase) Update(ctx context.Context, id int32, req *domain.UpdateFishRequest) (*models.Fish, error) {
	fish, err := u.fishRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		fish.Name = *req.Name
	}
	if req.Health != nil {
		fish.Health = *req.Health
	}
	if req.RewardMultiplier != nil {
		fish.RewardMultiplier = *req.RewardMultiplier
	}
	if req.Speed != nil {
		fish.Speed = *req.Speed
	}
	if req.AssetPath != nil {
		fish.AssetPath = *req.AssetPath
	}
	if err := u.fishRepo.Update(ctx, fish); err != nil {
		return nil, err
	}
	return fish, nil
}

func (u *fishUsecase) Delete(ctx context.Context, id int32) error {
	return u.fishRepo.Delete(ctx, id)
}
