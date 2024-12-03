package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	logger "github.com/sirupsen/logrus"

	"one-day-server/configs"
	"one-day-server/internal/httpserver/rest"
	"one-day-server/internal/management"
	"one-day-server/response"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, response.InvalidRequestBody)
		return
	}

	// validate Refresh Token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return configs.JWTSecret, nil
	})
	if err != nil || !token.Valid {
		response.SendError(c, response.UnauthorizedJWTAccessToken)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["uid"] == nil {
		response.SendError(c, response.InvalidJWTTokenClaims)
		return
	}

	user, err := management.GetUserByUid(claims["uid"].(int64))
	if err != nil {
		response.SendError(c, response.UserNotFound)
		return
	}
	// generate new Access Token
	accessToken, _, err := rest.GenerateTokens(user)
	if err != nil {
		logger.Errorf("generate jwt token failed, err: %s", err)
		response.SendInternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
