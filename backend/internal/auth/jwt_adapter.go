package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAdapter struct {
	Secret string
}

func NewJWTAdapter(secret string) *JWTAdapter {
	return &JWTAdapter{Secret: secret}
}

func (j *JWTAdapter) GenerateToken(userID string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "uptime-monitor-backend",
			"sub": userID,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
	s, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return s, nil
}

func (j *JWTAdapter) VerifyToken(token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return "", fmt.Errorf("verify token: %w", err)
	}

	return t.Claims.(jwt.MapClaims)["sub"].(string), nil
}
