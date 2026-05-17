package usecase

import (
	"context"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
)

type AuthUsecase interface {
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error)
	Me(ctx context.Context, userID int64) (*models.User, error)
}

type authUsecase struct {
	userRepo   repository.UserRepository
	hasher     utils.PasswordHasher
	tokenMaker utils.TokenMaker
}

func NewAuthUsecase(repo repository.UserRepository, hasher utils.PasswordHasher, token utils.TokenMaker) AuthUsecase {
	return &authUsecase{
		userRepo:   repo,
		hasher:     hasher,
		tokenMaker: token,
	}
}

func (u *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperror.ErrInvalidCredentials
	}

	if err = u.hasher.CompareHashAndPassword(user.Password, req.Password); err != nil {
		return nil, apperror.ErrInvalidCredentials
	}

	tokenString, expiresAt, err := u.tokenMaker.CreateToken(user.ID, user.RoleID, 24*time.Hour)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	return &domain.LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		User:      *user,
	}, nil
}

func (u *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	exists, err := u.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}
	if exists {
		return nil, apperror.ErrUsernameExisted
	}

	passwordHash, err := u.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	user := &models.User{
		Username: req.Username,
		Password: passwordHash,
		Email:    req.Email,
		RoleID:   1,
	}

	if err = u.userRepo.Create(ctx, user); err != nil {
		return nil, apperror.ErrInternalServer
	}

	return &domain.RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		RoleID:   user.RoleID,
	}, nil
}

func (u *authUsecase) Me(ctx context.Context, userID int64) (*models.User, error) {
	return u.userRepo.GetByID(ctx, userID)
}
