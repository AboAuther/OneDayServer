package public

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"one-day-server/internal/httpserver/middleware/sms"
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
		logger.Errorf("add user failed, err: %v", err)
		response.SendInternalServerError(c)
		return
	}
	response.SendSuccess(c, map[string]interface{}{
		"message": "success",
	})
}

func ForgotPassword(c *gin.Context) {
	var request struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		response.SendError(c, response.InvalidRequestBody)
		return
	}

	// check if user exists
	_, err := management.GetUserByPhone(request.Phone)
	if err != nil {
		logger.Errorf("user not found for phone: %s, err: %v", request.Phone, err)
		response.SendError(c, response.UserNotFound)
		return
	}

	// send sms by aliyun
	err = sms.SendSMSAndStoreCode(request.Phone)
	if err != nil {
		logger.Errorf("failed to send SMS: %v", err)
		response.SendError(c, response.SMSDeliveryFailed)
		return
	}

	response.SendSuccessMessage(c)
}
