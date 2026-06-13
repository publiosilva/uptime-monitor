package auth

import (
	"fmt"
	"net/mail"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(email, password string) (User, error) {
	email = strings.TrimSpace(email)
	if err := validateRegisterInput(email, password); err != nil {
		return User{}, err
	}

	hash, err := HashPassword(password)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	return s.repo.CreateUser(email, hash)
}

func validateRegisterInput(email, password string) error {
	if email == "" || password == "" {
		return ErrRequiredFieldMissing{Field: "email or password"}
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}

	if len(password) < 8 {
		return ErrInvalidPassword
	}

	return nil
}
