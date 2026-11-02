package auth

import (
	"fmt"
	"strings"
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
	token = strings.TrimSpace(token)
	if len(token) >= 7 && strings.EqualFold(token[:7], "bearer ") {
		token = strings.TrimSpace(token[7:])
	}

	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})
	if err != nil {
		return "", fmt.Errorf("verify token: %w", err)
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return "", fmt.Errorf("invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", fmt.Errorf("invalid subject claim")
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return "", fmt.Errorf("invalid expiration claim: %w", err)
	}
	if exp.Before(time.Now()) {
		return "", fmt.Errorf("token expired")
	}

	return sub, nil
}
