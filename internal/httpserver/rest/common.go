package rest

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"one-day-server/configs"
	"one-day-server/internal/management"
)

func GenerateTokens(user *management.User) (string, string, error) {
	// Access Token
	accessClaims := jwt.MapClaims{
		"uid":      user.Id,
		"username": user.Username,
		"phone":    user.Phone,
		"exp":      time.Now().Add(15 * time.Minute).Unix(), // 15 分钟有效
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSigned, err := accessToken.SignedString(configs.JWTSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshClaims := jwt.MapClaims{
		"uid": user.Id,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 天有效
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshSigned, err := refreshToken.SignedString(configs.JWTSecret)
	if err != nil {
		return "", "", err
	}

	return accessSigned, refreshSigned, nil
}
