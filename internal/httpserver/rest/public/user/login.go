package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"one-day-server/internal/httpserver/rest"
	"one-day-server/internal/management"
	"one-day-server/response"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func LoginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, response.InvalidRequestBody)
		return
	}

	var user *management.User
	var err error
	if req.Username != "" {
		user, err = management.GetUserByUsername(req.Username)
		if err != nil {
			user, err = management.GetUserByUsername(req.Username)
			if err != nil {
				response.SendError(c, response.UserNotFound)
				return
			}
		}
	}
	if user == nil {
		response.SendError(c, response.UserNotFound)
		return
	}

	// validate user password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.SendError(c, response.UnauthorizedUserPassword)
		return
	}

	accessToken, refreshToken, err := rest.GenerateTokens(user)
	if err != nil {
		logger.Errorf("generate jwt token failed, err: %s", err)
		response.SendInternalServerError(c)
		return
	}
	if err := management.UpdateUserRefreshToken(user, refreshToken); err != nil {
		logger.Errorf("save jwt refresh token failed, err: %s", err)
		response.SendInternalServerError(c)
		return
	}
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
