package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourname/fish-game-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles admin authentication and JWT token generation.
type AuthService interface {
	Login(ctx context.Context, username, password string) (*TokenPair, error)
	CreateAdmin(ctx context.Context, username, password, role string) (*AdminResponse, error)
}

// TokenPair holds the access and refresh tokens returned after login.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // seconds
}

// AdminResponse is the safe (no password hash) representation returned to the caller.
type AdminResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type authService struct {
	repo repository.AuthRepository
}

// NewAuthService creates an AuthService with the given repository.
func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

// Login verifies credentials against the DB and returns a JWT token pair.
func (s *authService) Login(ctx context.Context, username, password string) (*TokenPair, error) {
	admin, err := s.repo.FindAdminByUsername(ctx, username)
	if err != nil {
		// Don't reveal whether user exists or not
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password against bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "fallback-dev-secret"
	}

	// Access token — 24 hours
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  admin.Username,
		"role": admin.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	})
	accessTokenStr, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("authService: could not sign access token: %w", err)
	}

	// Refresh token — 7 days
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": admin.Username,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	refreshTokenStr, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("authService: could not sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    86400,
	}, nil
}

// CreateAdmin hashes the plain-text password and inserts a new admin record.
func (s *authService) CreateAdmin(ctx context.Context, username, password, role string) (*AdminResponse, error) {
	if role == "" {
		role = "admin"
	}
	if role != "admin" && role != "superadmin" {
		return nil, fmt.Errorf("invalid role %q: must be 'admin' or 'superadmin'", role)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("authService.CreateAdmin: failed to hash password: %w", err)
	}

	admin, err := s.repo.CreateAdmin(ctx, username, string(hash), role)
	if err != nil {
		return nil, fmt.Errorf("authService.CreateAdmin: %w", err)
	}

	return &AdminResponse{
		ID:       admin.ID,
		Username: admin.Username,
		Role:     admin.Role,
	}, nil
}
