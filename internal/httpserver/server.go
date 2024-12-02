package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"one-day-server/configs"
	"one-day-server/internal/httpserver/middleware"
	"one-day-server/internal/httpserver/middleware/cache"
	"one-day-server/internal/httpserver/rest/public"
	"one-day-server/response"
)

type OneDayServer struct {
	ginEngine *gin.Engine
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func NewOneDayServer() (*OneDayServer, error) {
	ginEngine := gin.Default()
	ginEngine.Use(Cors())

	ginEngine.NoRoute(func(c *gin.Context) {
		response.SendError(c, response.ApiNotFound)
	})
	ginEngine.HandleMethodNotAllowed = true
	ginEngine.NoMethod(func(c *gin.Context) {
		response.SendError(c, response.RequestMethodNotAllowed)
	})

	ginEngine.Use(cache.Instance.UseCache)
	userGroup := ginEngine.Group("/api/v1/oneDay/user")
	userGroup.Use(middleware.ValidateUserAuth)
	{

	}

	publicGroup := ginEngine.Group("/api/v1/oneDay/public")
	{
		publicGroup.GET("/timestamp", public.Time)
		publicGroup.POST("/user", public.CreateUser)
	}

	ginEngine.GET("/api/v1/oneDay/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ginEngine.GET("/metrics", func(context *gin.Context) {
		promhttp.Handler().ServeHTTP(context.Writer, context.Request)
	})
	return &OneDayServer{
		ginEngine: ginEngine,
	}, nil
}

func (s *OneDayServer) Start() error {
	port := configs.GetEnvDefault("SERVER_PORT", "8049")
	log.Infof("server is starting at: http://localhost:%v", port)
	if err := s.ginEngine.Run(":" + port); err != nil {
		return err
	}
	return nil
}
