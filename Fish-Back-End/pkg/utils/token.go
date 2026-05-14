package utils

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("token không hợp lệ")
	ErrExpiredToken = errors.New("token đã hết hạn")
)

type TokenMaker interface {
	CreateToken(adminID string, role string, duration time.Duration) (string, int64, error)
	ExtractToken(tokenString string) (*jwt.MapClaims, error)
}

type jwtMaker struct {
	secretKey     string
	signingMethod jwt.SigningMethod
}

func NewTokenMaker(secretKey string, signingMethod jwt.SigningMethod) TokenMaker {
	return &jwtMaker{
		secretKey:     secretKey,
		signingMethod: signingMethod,
	}
}

func (m *jwtMaker) CreateToken(adminID string, role string, duration time.Duration) (string, int64, error) {
	expirationTime := time.Now().Add(duration)

	claims := jwt.MapClaims{
		"admin_id":  adminID,
		"role":      role,
		"exp":       expirationTime.Unix(),
		"issued_at": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(m.signingMethod, claims)

	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", 0, fmt.Errorf("không thể ký token: %w", err)
	}

	return tokenString, expirationTime.Unix(), nil
}

func (m *jwtMaker) ExtractToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != m.signingMethod.Alg() {
			return nil, fmt.Errorf("thuật toán ký không khớp: kỳ vọng %s", m.signingMethod.Alg())
		}

		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	log.Print(claims)
	return &claims, nil
}
