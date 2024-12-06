package rest

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"one-day-server/configs"
	internalRedis "one-day-server/internal/db/redis"
	"one-day-server/internal/management"
)

var RedisClient internalRedis.IClient

func RegisterRedisClient(client internalRedis.IClient) {
	RedisClient = client
}

func GenerateTokens(user *management.User) (string, string, error) {
	accessJTI := uuid.NewString()
	refreshJTI := uuid.NewString()

	// Access Token Claims
	accessClaims := jwt.MapClaims{
		"uid":      user.Id,
		"username": user.Username,
		"phone":    user.Phone,
		"jti":      accessJTI,
		"exp":      time.Now().Add(15 * time.Minute).Unix(), // 15 分钟有效
		"iat":      time.Now().Unix(),
	}

	// Refresh Token Claims
	refreshClaims := jwt.MapClaims{
		"uid": user.Id,
		"jti": refreshJTI,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 天有效
		"iat": time.Now().Unix(),
	}

	// 签发 Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSigned, err := accessToken.SignedString(configs.JWTSecret)
	if err != nil {
		return "", "", err
	}

	// 签发 Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshSigned, err := refreshToken.SignedString(configs.JWTSecret)
	if err != nil {
		return "", "", err
	}

	return accessSigned, refreshSigned, nil
}
