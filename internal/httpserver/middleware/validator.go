package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"one-day-server/internal/management"
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
	user := ValidateAndGetUser(c)
	if user == nil {
		return
	}

	c.Set(utils.UserInContext, user)
	c.Next()
}
