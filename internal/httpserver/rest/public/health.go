package public

import (
	"time"

	"github.com/gin-gonic/gin"

	"one-day-server/response"
)

const (
	UPTIME_ROBOT_RESPONSE_KEY   = "GATEWAY:uptimerobot:status"
	UPTIME_ROBOT_RESET          = "GATEWAY:uptimerobot:reset"
	ERROR_MSG                   = "Request Failed"
	RESET_THRESHOLD             = 5
	HTTP_TIMEOUT_SECONDS        = 5
	REDIS_TIMEOUT_SECONDS       = 1
	RESPONSE_CACHE_TIME_SECONDS = 60
)

type TimeResponse struct {
	ServerTime int64 `json:"serverTime"`
}

// Time          查询服务器时间
// @Summary      查询服务器时间
// @Description  查询当前服务器时间
// @Tags         gateway/public
// @Success           200         TimeResponse     TimeResponse
// @Router       /api/v1/gateway/public/time [get]
func Time(c *gin.Context) {
	response.SendSuccess(c, &TimeResponse{ServerTime: time.Now().UnixMilli()})
}

func CreateUser(c *gin.Context) {

}
