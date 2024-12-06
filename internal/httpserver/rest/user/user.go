package user

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	logger "github.com/sirupsen/logrus"

	"one-day-server/configs"
	"one-day-server/internal/httpserver/rest"
	"one-day-server/internal/management"
	"one-day-server/response"
)

func isBlacklisted(jti string) bool {
	ctx := context.Background()
	val, err := rest.RedisClient.GetResult(ctx, "blacklist:"+jti)
	return err == nil && val == "1"
}

func addToBlacklist(jti string, exp time.Time) error {
	ctx := context.Background()
	ttl := time.Until(exp)
	return rest.RedisClient.WriteResultWithTTL(ctx, "blacklist:"+jti, "1", ttl)
}

func LogOut(c *gin.Context) {
	accessTokenJTI, exists := c.Get("jti")
	if !exists {
		response.SendError(c, response.InvalidJWTTokenClaims)
		return
	}

	// 获取 Refresh Token 的 jti（假设它来自请求头或 cookie 中）
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		response.SendError(c, response.MissingRequiredHeader, "Refresh-Token")
		return
	}

	refreshTokenObj, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return configs.JWTSecret, nil
	})

	if err != nil || !refreshTokenObj.Valid {
		response.SendError(c, response.UnauthorizedJWTRefreshToken)
		return
	}

	// 获取 Refresh Token 的 jti
	refreshClaims, ok := refreshTokenObj.Claims.(jwt.MapClaims)
	if !ok || refreshClaims["jti"] == nil {
		response.SendError(c, response.InvalidJWTTokenClaims)
		return
	}
	refreshTokenJTI := refreshClaims["jti"].(string)

	// 获取 Access Token 和 Refresh Token 的过期时间
	accessExp := time.Unix(int64(refreshClaims["exp"].(float64)), 0)
	refreshExp := time.Unix(int64(refreshClaims["exp"].(float64)), 0)

	// 添加到 Redis 黑名单
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

	// 返回成功响应
	c.JSON(200, gin.H{"message": "Logout successful"})
}

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
