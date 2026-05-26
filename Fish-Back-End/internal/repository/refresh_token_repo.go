package repository

import (
	"context"
	"errors"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository/dbgen"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	GetByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	DeleteByHash(ctx context.Context, tokenHash string) error
	DeleteByUserID(ctx context.Context, userID int64) error
}

type refreshTokenPgRepo struct {
	queries *dbgen.Queries
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) RefreshTokenRepository {
	return &refreshTokenPgRepo{queries: dbgen.New(pool)}
}

func mapToModelRefreshToken(r dbgen.RefreshToken) models.RefreshToken {
	return models.RefreshToken{
		ID:        r.ID,
		UserID:    r.UserID,
		TokenHash: r.TokenHash,
		ExpiresAt: r.ExpiresAt.Time,
		CreatedAt: r.CreatedAt.Time,
	}
}

func (r *refreshTokenPgRepo) Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	_, err := r.queries.CreateRefreshToken(ctx, dbgen.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return apperror.Wrap("repository", "refreshTokenRepo.Create", err)
	}
	return nil
}

func (r *refreshTokenPgRepo) GetByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	row, err := r.queries.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrInvalidToken
		}
		return nil, apperror.Wrap("repository", "refreshTokenRepo.GetByHash", err)
	}
	token := mapToModelRefreshToken(row)
	return &token, nil
}

func (r *refreshTokenPgRepo) DeleteByHash(ctx context.Context, tokenHash string) error {
	if err := r.queries.DeleteRefreshTokenByHash(ctx, tokenHash); err != nil {
		return apperror.Wrap("repository", "refreshTokenRepo.DeleteByHash", err)
	}
	return nil
}

func (r *refreshTokenPgRepo) DeleteByUserID(ctx context.Context, userID int64) error {
	if err := r.queries.DeleteRefreshTokensByUserID(ctx, userID); err != nil {
		return apperror.Wrap("repository", "refreshTokenRepo.DeleteByUserID", err)
	}
	return nil
}
