package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/apperror"
	"github.com/golang-jwt/jwt/v5"
)

type TokenMaker interface {
	CreateAccessToken(userID int64, roleID int32) (string, int64, error)
	CreateRefreshToken(userID int64) (string, int64, error)
	VerifyAccessToken(token string) (*jwt.MapClaims, error)
	VerifyRefreshToken(token string) (*jwt.MapClaims, error)
}

type jwtMaker struct {
	accessKey     string
	refreshKey    string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	signingMethod jwt.SigningMethod
}

func NewTokenMaker(
	accessKey, refreshKey string,
	accessExpiry, refreshExpiry time.Duration,
	signingMethod jwt.SigningMethod,
) TokenMaker {
	return &jwtMaker{
		accessKey:     accessKey,
		refreshKey:    refreshKey,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		signingMethod: signingMethod,
	}
}

func (m *jwtMaker) createToken(claims jwt.MapClaims, key string) (string, error) {
	token := jwt.NewWithClaims(m.signingMethod, claims)
	signed, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("ký token thất bại: %w", err)
	}
	return signed, nil
}

func (m *jwtMaker) CreateAccessToken(userID int64, roleID int32) (string, int64, error) {
	exp := time.Now().Add(m.accessExpiry).Unix()
	tokenStr, err := m.createToken(jwt.MapClaims{
		"user_id": userID,
		"role_id": roleID,
		"type":    "access",
		"exp":     exp,
	}, m.accessKey)
	return tokenStr, exp, err
}

func (m *jwtMaker) CreateRefreshToken(userID int64) (string, int64, error) {
	exp := time.Now().Add(m.refreshExpiry).Unix()
	tokenStr, err := m.createToken(jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     exp,
	}, m.refreshKey)
	return tokenStr, exp, err
}

func (m *jwtMaker) verifyToken(tokenStr, key, expectedType string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != m.signingMethod.Alg() {
			return nil, fmt.Errorf("thuật toán không khớp")
		}
		return []byte(key), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apperror.ErrExpiredToken
		}
		return nil, apperror.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, apperror.ErrInvalidToken
	}
	if claims["type"] != expectedType {
		return nil, apperror.ErrInvalidToken
	}
	return &claims, nil
}

func (m *jwtMaker) VerifyAccessToken(token string) (*jwt.MapClaims, error) {
	return m.verifyToken(token, m.accessKey, "access")
}

func (m *jwtMaker) VerifyRefreshToken(token string) (*jwt.MapClaims, error) {
	return m.verifyToken(token, m.refreshKey, "refresh")
}
