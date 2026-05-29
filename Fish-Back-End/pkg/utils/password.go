package utils

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	CompareHashAndPassword(hash string, password string) error
	HashPassword(password string) (string, error)
}

type bcryptHasher struct{}

func NewPasswordHasher() PasswordHasher {
	return &bcryptHasher{}
}

func (h *bcryptHasher) CompareHashAndPassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (h *bcryptHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
