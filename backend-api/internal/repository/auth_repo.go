package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yourname/fish-game-backend/internal/models"
	"github.com/yourname/fish-game-backend/pkg/database"
)

// AuthRepository defines the interface for admin authentication data access.
type AuthRepository interface {
	FindAdminByUsername(ctx context.Context, username string) (*models.Admin, error)
	CreateAdmin(ctx context.Context, username, passwordHash, role string) (*models.Admin, error)
}

type authSQLServerRepo struct {
	db *database.DB
}

// NewAuthRepository creates a new SQL Server–backed AuthRepository.
func NewAuthRepository(db *database.DB) AuthRepository {
	return &authSQLServerRepo{db: db}
}

// FindAdminByUsername looks up an admin record by username.
// Returns nil, ErrNoRows-wrapped error if not found.
func (r *authSQLServerRepo) FindAdminByUsername(ctx context.Context, username string) (*models.Admin, error) {
	const query = `
		SELECT id, username, password_hash, role
		FROM admins
		WHERE username = @p1`

	var a models.Admin
	err := r.db.QueryRowContext(ctx, query, sql.Named("p1", username)).
		Scan(&a.ID, &a.Username, &a.PasswordHash, &a.Role)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("authRepo: admin %q not found", username)
	}
	if err != nil {
		return nil, fmt.Errorf("authRepo.FindAdminByUsername: %w", err)
	}
	return &a, nil
}

// CreateAdmin inserts a new admin and returns the created record.
func (r *authSQLServerRepo) CreateAdmin(ctx context.Context, username, passwordHash, role string) (*models.Admin, error) {
	const query = `
		INSERT INTO admins (username, password_hash, role)
		OUTPUT INSERTED.id, INSERTED.username, INSERTED.password_hash, INSERTED.role
		VALUES (@p1, @p2, @p3)`

	var a models.Admin
	err := r.db.QueryRowContext(ctx, query,
		sql.Named("p1", username),
		sql.Named("p2", passwordHash),
		sql.Named("p3", role),
	).Scan(&a.ID, &a.Username, &a.PasswordHash, &a.Role)
	if err != nil {
		return nil, fmt.Errorf("authRepo.CreateAdmin: %w", err)
	}
	return &a, nil
}
