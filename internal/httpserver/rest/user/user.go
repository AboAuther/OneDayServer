package user

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"

	"one-day-server/configs"
	"one-day-server/internal/httpserver/rest"
	"one-day-server/internal/management"
	"one-day-server/response"
	"one-day-server/utils"
)

func addToBlacklist(jti string, exp time.Time) error {
	ctx := context.Background()
	ttl := time.Until(exp)
	if ttl <= 0 {
		return fmt.Errorf("token already expired")
	}
	return rest.RedisClient.WriteResultWithTTL(ctx, "blacklist:"+jti, "1", ttl)
}

func LogOut(c *gin.Context) {
	accessTokenJTI, exists := c.Get("jti")
	uid := c.GetInt64("uid")
	if !exists {
		response.SendError(c, response.InvalidJWTTokenClaims)
		return
	}

	// get Refresh Token jti
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		response.SendError(c, response.MissingRequiredHeader, "Refresh-Token")
		return
	}

	refreshClaims, err := utils.ParseJWT(refreshToken, configs.JWTSecret)
	if err != nil {
		logger.Errorf("failed to parse refresh token for user %d: %v", uid, err)
		response.SendError(c, response.UnauthorizedJWTRefreshToken)
		return
	}

	refreshTokenJTI, ok := refreshClaims["jti"].(string)
	if !ok || refreshClaims["jti"] == nil {
		logger.Errorf("refresh token missing jti for user %d", uid)
		response.SendError(c, response.InvalidJWTTokenClaims)
		return
	}

	// get Access Token and Refresh Token expired time
	accessExp := time.Unix(c.GetInt64("exp"), 0)
	refreshExp := time.Unix(int64(refreshClaims["exp"].(float64)), 0)

	if err := addToBlacklist(accessTokenJTI.(string), accessExp); err != nil {
		logger.Errorf("failed to blacklist access token jti: %v", err)
		response.SendError(c, response.InternalServerError)
		return
	}

	if err := addToBlacklist(refreshTokenJTI, refreshExp); err != nil {
		logger.Errorf("failed to blacklist refresh token jti: %v", err)
		response.SendError(c, response.InternalServerError)
		return
	}

	response.SendSuccessMessage(c)
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, response.InvalidRequestBody)
		return
	}

	// validate Refresh Token
	tokenClaims, err := utils.ParseJWT(req.RefreshToken, configs.JWTSecret)
	if err != nil {
		response.SendError(c, response.UnauthorizedJWTRefreshToken)
		return
	}

	user, err := management.GetUserByUid(int64(tokenClaims["uid"].(float64)))
	if err != nil {
		logger.Errorf("get uid from refresh token claims failed, claims: %v, err: %v", tokenClaims, err)
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
	response.SendSuccess(c, map[string]interface{}{
		"access_token": accessToken,
	})
}

type UpdateUserProfileRequest struct {
	Username string `json:"username" binding:"required"`
	Phone    string `json:"phone"`
	Email    string `gorm:"column:email"`
	Gender   string `gorm:"column:gender"`
	Age      int    `gorm:"column:age"`
	IsVip    bool   `gorm:"column:is_vip"`
}

func UpdateUserProfile(c *gin.Context) {
	var req UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, response.InvalidRequestBody)
		return
	}

	user, err := management.GetUserByUsername(req.Username)
	if err != nil {
		logger.Errorf("get user by username failed, err: %s", err)
		response.SendError(c, response.UserNotFound)
		return
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Age > 0 {
		user.Age = req.Age
	}

	if err := management.UpdateUser(user); err != nil {
		logger.Errorf("update user failed, err: %s", err)
		response.SendError(c, response.InternalServerError)
		return
	}
	response.SendSuccess(c, user)
}
