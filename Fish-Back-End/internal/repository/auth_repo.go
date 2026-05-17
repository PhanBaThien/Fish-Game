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

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type userPgRepo struct {
	pool    *pgxpool.Pool
	queries *dbgen.Queries
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userPgRepo{
		pool:    pool,
		queries: dbgen.New(pool),
	}
}

func mapToModelUser(u dbgen.User) models.User {
	return models.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		RoleID:    u.RoleID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (r *userPgRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	exists, err := r.queries.CheckUserExists(ctx, username)
	if err != nil {
		return false, apperror.Wrap("repository", "userRepo.ExistsByUsername", err)
	}
	return exists, nil
}

func (r *userPgRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	res, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, apperror.Wrap("repository", "userRepo.GetByUsername", err)
	}
	user := mapToModelUser(res)
	return &user, nil
}

func (r *userPgRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	res, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, apperror.Wrap("repository", "userRepo.GetByID", err)
	}
	user := mapToModelUser(res)
	return &user, nil
}

func (r *userPgRepo) Create(ctx context.Context, user *models.User) error {
	res, err := r.queries.CreateUser(ctx, dbgen.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		RoleID:   user.RoleID,
	})
	if err != nil {
		return apperror.Wrap("repository", "userRepo.Create", err)
	}
	user.ID = res.ID
	user.CreatedAt = res.CreatedAt
	user.UpdatedAt = res.UpdatedAt
	return nil
}
