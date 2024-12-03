package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"one-day-server/configs"
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

func ValidateUserAuth(c *gin.Context) {
	authHeader := c.GetHeader(utils.OneDayAuthorization)
	if authHeader == "" {
		response.SendError(c, response.MissingRequiredHeader, utils.OneDayAuthorization)
		c.Abort()
		return
	}

	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return configs.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		response.SendError(c, response.UnauthorizedJWTAccessToken)
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("uid", int64(claims["uid"].(float64)))
		c.Set("username", claims["username"].(string))
	} else {
		response.SendError(c, response.InvalidJWTTokenClaims)
		c.Abort()
		return
	}
	c.Next()
}
