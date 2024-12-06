package main

import (
	log "github.com/sirupsen/logrus"

	"one-day-server/internal/db/mysql"
	"one-day-server/internal/httpserver"
	"one-day-server/internal/management"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server.
// @termsOfService  https://lixueduan.com

// @contact.name   lixd
// @contact.url    https://lixueduan.com
// @contact.email  xueduan.li@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8049
// @BasePath  /api/v1/oneDay

// SwaggerUI: http://localhost:8049/swagger/index.html
func main() {
	log.Infof("start gateway server")
	mysql.Init()
	management.Init()

	server, err := httpserver.NewOneDayServer()
	if err != nil {
		panic(err)
	}
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
