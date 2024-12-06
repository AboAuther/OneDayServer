package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	logger "github.com/sirupsen/logrus"

	"one-day-server/configs"
	internalRedis "one-day-server/internal/db/redis"
	"one-day-server/internal/management"
	"one-day-server/response"
	"one-day-server/utils"
)

func requestExpired(vesselTimestamp int64) bool {
	currentTime := time.Now().UnixMilli()
	return currentTime-vesselTimestamp > utils.RecvWindow || vesselTimestamp-currentTime > utils.RecvWindow
}

func ValidateAndGetUser(c *gin.Context) *management.User {
	return nil
}

func isBlacklisted(jti string) bool {
	ctx := context.Background()
	val, err := internalRedis.GetClient().GetResult(ctx, "blacklist:"+jti)
	return err == nil && val == "1"
}

func ValidateUserAuth(c *gin.Context) {
	authHeader := c.GetHeader(utils.OneDayAuthorization)
	if authHeader == "" {
		response.SendError(c, response.MissingRequiredHeader, utils.OneDayAuthorization)
		c.Abort()
		return
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		response.SendError(c, response.InvalidJWTTokenFormat)
		c.Abort()
		return
	}

	tokenString := authHeader[len(bearerPrefix):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return configs.JWTSecret, nil
	})

	if err != nil {
		logger.Errorf("parse jwt token failed, err: %s", err)
		response.SendError(c, response.UnauthorizedJWTAccessToken)
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		response.SendError(c, response.InvalidJWTTokenClaims)
		c.Abort()
		return
	}

	jti, ok := claims["jti"].(string)
	if !ok || isBlacklisted(jti) {
		response.SendError(c, response.UnauthorizedJWTAccessToken)
		c.Abort()
		return
	}

	// 检查必要字段
	uid, ok := claims["uid"].(float64)
	username, ok2 := claims["username"].(string)
	if !ok || !ok2 {
		response.SendError(c, response.InvalidJWTTokenClaims)
		c.Abort()
		return
	}

	c.Set("uid", int64(uid))
	c.Set("username", username)
	c.Set("jti", jti)
	c.Next()
}
