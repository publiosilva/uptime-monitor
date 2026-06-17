package auth

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

type Service struct {
	repo Repository
	jwt  *JWTAdapter
}

func NewService(repo Repository, jwt *JWTAdapter) *Service {
	return &Service{repo: repo, jwt: jwt}
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

func (s *Service) Login(email, password string) (string, error) {
	email = strings.TrimSpace(email)
	if email == "email" || password == "" {
		return "", ErrRequiredFieldMissing{Field: "email or password"}
	}

	user, err := s.repo.GetUserByEmail(email)
	if errors.Is(err, ErrUserNotFound) {
		return "", ErrInvalidCredentials
	} else if err != nil {
		return "", fmt.Errorf("get user by email: %w", err)
	}

	if !ComparePassword(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwt.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}
