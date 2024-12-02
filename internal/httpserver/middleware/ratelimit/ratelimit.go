package ratelimit

import (
	"net"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	log "github.com/sirupsen/logrus"

	"one-day-server/configs"
	internalRedis "one-day-server/internal/db/redis"
	"one-day-server/response"
	"one-day-server/utils"
)

const (
	RedisPrefix = "GATEWAY:RateLimiter:"
)

func parseKeyFromGin(ctx *gin.Context) (key string, doLimit bool) {
	if strings.HasSuffix(ctx.Request.URL.String(), "/health") {
		return "", false
	}
	ip := net.ParseIP(ctx.ClientIP())
	return RedisPrefix + ip.String(), !ip.IsPrivate()
}

func ValidateRateLimit(ctx *gin.Context, weight int) {
	client := internalRedis.GetClient().GetRateLimiter()
	key, doLimit := parseKeyFromGin(ctx)
	if doLimit {
		res, err := client.AllowN(ctx, key, redis_rate.PerMinute(configs.GetEnvDefaultInt("RATE_LIMIT_PER_MINUTE", 600)), weight)
		if err != nil {
			log.Errorf("rate limit error: %v", err)
			//ctx.Next()
			return
		}
		if res.Allowed <= 0 {
			response.SendError(ctx, response.IpAddressRateLimitExceeded)
			ctx.Abort()
			return
		}
		// set remaining in response header
		ctx.Header(utils.OneDayAPIIPRateLimitRemaining, strconv.Itoa(res.Remaining))
	}
	//ctx.Next()
}
