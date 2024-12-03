package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"one-day-server/internal/httpserver/rest"
	"one-day-server/internal/management"
	"one-day-server/response"
)

type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Phone    string `json:"phone" binding:"required"`
}

func RegisterUser(c *gin.Context) {
	var req CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, response.InvalidRequestBody)
		return
	}

	// encrypt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("hash password failed, err: %v", err)
		response.SendInternalServerError(c)
		return
	}

	if err = management.AddUser(&management.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Phone:    req.Phone,
	}); err != nil {
		response.SendInternalServerError(c)
		return
	}
	response.SendSuccess(c, map[string]interface{}{
		"message": "success",
	})
}

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
	user.RefreshToken = refreshToken
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
