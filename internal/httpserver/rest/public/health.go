package public

import (
	"time"

	"github.com/gin-gonic/gin"

	"one-day-server/response"
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
