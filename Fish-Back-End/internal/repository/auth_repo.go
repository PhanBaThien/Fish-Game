package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
)

var (
	ErrAdminNotFound = errors.New("không tìm thấy admin")
)

const (
	queryGetAdminByUsername = `
		SELECT id, username, password_hash, role, created_at
		FROM admin_users
		WHERE username = $1`

	queryGetAdminByID = `
		SELECT id, username, password_hash, role, created_at
		FROM admin_users
		WHERE id = $1`

	queryCreateAdmin = `
		INSERT INTO admin_users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
)

// AdminRepository định nghĩa interface thao tác với bảng admins
type AdminRepository interface {
	Create(ctx context.Context, admin *models.Admin) error
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	GetByUsername(ctx context.Context, username string) (*models.Admin, error)
	GetByID(ctx context.Context, id string) (*models.Admin, error)
}

type adminPgRepo struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) AdminRepository {
	return &adminPgRepo{db: db}
}

func (r *adminPgRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM admins WHERE username = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	return exists, err
}

// GetByUsername tìm kiếm admin theo tên đăng nhập
func (r *adminPgRepo) GetByUsername(ctx context.Context, username string) (*models.Admin, error) {
	var a models.Admin
	err := r.db.QueryRowContext(ctx, queryGetAdminByUsername, username).
		Scan(&a.ID, &a.Username, &a.PasswordHash, &a.Role, &a.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAdminNotFound // Trả về lỗi chuẩn đã định nghĩa
		}
		return nil, fmt.Errorf("adminRepo.GetByUsername: %w", err)
	}

	return &a, nil
}

// GetByID tìm kiếm admin theo ID nguyên bản
func (r *adminPgRepo) GetByID(ctx context.Context, id string) (*models.Admin, error) {
	var a models.Admin
	err := r.db.QueryRowContext(ctx, queryGetAdminByID, id).
		Scan(&a.ID, &a.Username, &a.PasswordHash, &a.Role, &a.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAdminNotFound // Trả về lỗi chuẩn đã định nghĩa
		}
		return nil, fmt.Errorf("adminRepo.GetByID: %w", err)
	}

	return &a, nil
}

// Create thêm một admin mới vào database
func (r *adminPgRepo) Create(ctx context.Context, admin *models.Admin) error {
	// 1. Khởi tạo một Transaction (Giao dịch)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("không thể mở transaction: %w", err)
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, queryCreateAdmin, admin.Username, admin.Email, admin.PasswordHash, admin.Role).
		Scan(&admin.ID, &admin.CreatedAt)

	if err != nil {
		return fmt.Errorf("lỗi insert admin: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("lỗi commit transaction: %w", err)
	}

	return nil
}
