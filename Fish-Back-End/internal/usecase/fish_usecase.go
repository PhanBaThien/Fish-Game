package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository"
)

var ErrFishNotFound = errors.New("không tìm thấy cấu hình cá")
var ErrFishIDExists = errors.New("mã cá đã tồn tại")

// FishUsecase định nghĩa nghiệp vụ quản lý cấu hình cá
type FishUsecase interface {
	ListFish(ctx context.Context, req *domain.ListFishRequest) (*domain.ListFishResponse, error)
	GetFish(ctx context.Context, id string) (*models.Fish, error)
	CreateFish(ctx context.Context, req *domain.CreateFishRequest) (*models.Fish, error)
	UpdateFish(ctx context.Context, id string, req *domain.UpdateFishRequest) (*models.Fish, error)
	DeleteFish(ctx context.Context, id string) error
}

type fishUsecase struct {
	fishRepo repository.FishRepository
}

func NewFishUsecase(fishRepo repository.FishRepository) FishUsecase {
	return &fishUsecase{fishRepo: fishRepo}
}

func (u *fishUsecase) ListFish(ctx context.Context, req *domain.ListFishRequest) (*domain.ListFishResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}
	offset := (req.Page - 1) * req.Limit

	fishes, total, err := u.fishRepo.List(ctx, req.Role, req.IsActive, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("fishUsecase.ListFish: %w", err)
	}

	return &domain.ListFishResponse{
		Items: fishes,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

func (u *fishUsecase) GetFish(ctx context.Context, id string) (*models.Fish, error) {
	f, err := u.fishRepo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrFishNotFound) {
		return nil, ErrFishNotFound
	}
	return f, err
}

func (u *fishUsecase) CreateFish(ctx context.Context, req *domain.CreateFishRequest) (*models.Fish, error) {
	// Kiểm tra ID đã tồn tại chưa
	if _, err := u.fishRepo.GetByID(ctx, req.ID); err == nil {
		return nil, ErrFishIDExists
	}

	f := &models.Fish{
		ID:         req.ID,
		Name:       req.Name,
		Multiplier: req.Multiplier,
		BaseProb:   req.BaseProb,
		Speed:      req.Speed,
		Role:       req.Role,
		IsActive:   req.IsActive,
	}

	if err := u.fishRepo.Create(ctx, f); err != nil {
		return nil, fmt.Errorf("fishUsecase.CreateFish: %w", err)
	}
	return f, nil
}

func (u *fishUsecase) UpdateFish(ctx context.Context, id string, req *domain.UpdateFishRequest) (*models.Fish, error) {
	f, err := u.fishRepo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrFishNotFound) {
		return nil, ErrFishNotFound
	}
	if err != nil {
		return nil, err
	}

	f.Name = req.Name
	f.Multiplier = req.Multiplier
	f.BaseProb = req.BaseProb
	f.Speed = req.Speed
	f.Role = req.Role
	f.IsActive = req.IsActive

	if err := u.fishRepo.Update(ctx, f); err != nil {
		return nil, fmt.Errorf("fishUsecase.UpdateFish: %w", err)
	}
	return f, nil
}

func (u *fishUsecase) DeleteFish(ctx context.Context, id string) error {
	err := u.fishRepo.Delete(ctx, id)
	if errors.Is(err, repository.ErrFishNotFound) {
		return ErrFishNotFound
	}
	return err
}
