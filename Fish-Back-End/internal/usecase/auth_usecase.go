package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/models"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
)

var (
	ErrInvalidCredentials = errors.New("tài khoản hoặc mật khẩu không đúng")
	ErrUsernameExisted    = errors.New("tài khoản đã tồn tại")
	ErrInternalServer     = errors.New("lỗi khi tạo token nội bộ")
)

// AuthUsecase định nghĩa các nghiệp vụ liên quan đến xác thực
type AuthUsecase interface {
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error)
}

type authUsecase struct {
	adminRepo  repository.AdminRepository
	hasher     utils.PasswordHasher
	tokenMaker utils.TokenMaker
}

// NewAuthUsecase khởi tạo usecase với các dependencies (Repo, Hasher, TokenMaker)
func NewAuthUsecase(repo repository.AdminRepository, hasher utils.PasswordHasher, token utils.TokenMaker) AuthUsecase {
	return &authUsecase{
		adminRepo:  repo,
		hasher:     hasher,
		tokenMaker: token,
	}
}

// Login xử lý kiểm tra thông tin đăng nhập và cấp phát JWT token
func (u *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	admin, err := u.adminRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = u.hasher.CompareHashAndPassword(admin.PasswordHash, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	tokenString, expiresAt, err := u.tokenMaker.CreateToken(admin.ID, admin.Role, 24*time.Hour)
	if err != nil {
		return nil, ErrInternalServer
	}

	return &domain.LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		Admin:     *admin,
	}, nil
}

// Register xử lý logic tạo mới tài khoản admin vào hệ thống (PostgreSQL)
func (u *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	// 1. Kiểm tra tên đăng nhập đã tồn tại chưa
	exists, err := u.adminRepo.ExistsByUsername(ctx, req.Username)
	if exists {
		return nil, ErrUsernameExisted
	}

	// 2. Băm mật khẩu (Hash)
	passwordHash, err := u.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, ErrInternalServer
	}

	// 3. Chuẩn bị dữ liệu Admin
	// Lưu ý: Không gán ID ở đây, để Postgres tự sinh (Auto Increment)
	admin := &models.Admin{
		Username:     req.Username,
		PasswordHash: passwordHash,
		Email:        req.Email,
		Role:         "admin",
	}

	// 4. Gọi Repository để lưu vào DB
	err = u.adminRepo.Create(ctx, admin)
	if err != nil {
		return nil, err
	}

	// 5. Trả về kết quả (ID lúc này đã có giá trị từ DB trả về)
	return &domain.RegisterResponse{
		ID:       admin.ID, // ID đã được DB cấp phát
		Username: admin.Username,
		Role:     admin.Role,
	}, nil
}


func (u *authUsecase) Me(ctx context.Context, tokenString string) (*models.Admin, error) {
	token, err := u.tokenMaker.ExtractToken(tokenString)
	
	if err != nil {
		return nil, err
	}

	adminID := (*token)["admin_id"].(string)
	log.Print(adminID)
	return u.adminRepo.GetByID(ctx, adminID)
}

