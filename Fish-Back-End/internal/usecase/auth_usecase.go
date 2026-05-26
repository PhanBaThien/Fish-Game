package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authUsecase struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	hasher           utils.PasswordHasher
	tokenMaker       utils.TokenMaker
}

func NewAuthUsecase(
	repo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	hasher utils.PasswordHasher,
	token utils.TokenMaker,
) AuthUsecase {
	return &authUsecase{
		userRepo:         repo,
		refreshTokenRepo: refreshTokenRepo,
		hasher:           hasher,
		tokenMaker:       token,
	}
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func (u *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperror.ErrInvalidCredentials
	}
	if err = u.hasher.CompareHashAndPassword(user.Password, req.Password); err != nil {
		return nil, apperror.ErrInvalidCredentials
	}

	accessToken, accessExp, err := u.tokenMaker.CreateAccessToken(user.ID, user.RoleID)
	if err != nil {
		return nil, apperror.Wrap("usecase", "authUsecase.Login.CreateAccessToken", err)
	}

	refreshToken, refreshExp, err := u.tokenMaker.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, apperror.Wrap("usecase", "authUsecase.Login.CreateRefreshToken", err)
	}

	if err = u.refreshTokenRepo.Create(ctx, user.ID, hashToken(refreshToken),
		time.Unix(refreshExp, 0)); err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExp,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExp,
		User:                  *user,
	}, nil
}

func (u *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	exists, err := u.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperror.ErrUsernameExisted
	}

	passwordHash, err := u.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.Wrap("usecase", "authUsecase.Register.HashPassword", err)
	}

	user := &models.User{
		Username: req.Username,
		Password: passwordHash,
		Email:    req.Email,
		RoleID:   1,
	}
	if err = u.userRepo.Create(ctx, user); err != nil {
		return nil, err
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

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenResponse, error) {
	claims, err := u.tokenMaker.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	record, err := u.refreshTokenRepo.GetByHash(ctx, hashToken(refreshToken))
	if err != nil {
		return nil, apperror.ErrInvalidToken
	}
	if time.Now().After(record.ExpiresAt) {
		_ = u.refreshTokenRepo.DeleteByHash(ctx, hashToken(refreshToken))
		return nil, apperror.ErrExpiredToken
	}

	userID := utils.ToInt64((*claims)["user_id"])
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	_ = u.refreshTokenRepo.DeleteByHash(ctx, hashToken(refreshToken))

	accessToken, accessExp, err := u.tokenMaker.CreateAccessToken(user.ID, user.RoleID)
	if err != nil {
		return nil, apperror.Wrap("usecase", "authUsecase.RefreshToken.CreateAccessToken", err)
	}
	newRefreshToken, refreshExp, err := u.tokenMaker.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, apperror.Wrap("usecase", "authUsecase.RefreshToken.CreateRefreshToken", err)
	}

	if err = u.refreshTokenRepo.Create(ctx, user.ID, hashToken(newRefreshToken),
		time.Unix(refreshExp, 0)); err != nil {
		return nil, err
	}

	return &domain.RefreshTokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExp,
		RefreshToken:          newRefreshToken,
		RefreshTokenExpiresAt: refreshExp,
	}, nil
}

func (u *authUsecase) Logout(ctx context.Context, refreshToken string) error {
	return u.refreshTokenRepo.DeleteByHash(ctx, hashToken(refreshToken))
}
