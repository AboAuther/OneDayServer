package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
)

func init() {
	//var err error
	err := godotenv.Load()
	if err != nil {
		if os.IsNotExist(err) {
			logger.Infof(".env file is not set, use default environment variables")
		} else {
			logger.Panicf("Error loading .env file, err: %s", err)
		}
	}
	level, err := logger.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logger.SetLevel(logger.InfoLevel)
	} else {
		logger.SetLevel(level)
	}
}

func GetEnvDefault(key string, def string) string {
	env, exist := os.LookupEnv(key)
	if !exist {
		return def
	}
	return env
}

func GetEnvDefaultInt(key string, def int) int {
	env, exist := os.LookupEnv(key)
	if !exist {
		return def
	}
	v, err := strconv.Atoi(env)
	if err != nil {
		log.Error("env %s value is %s, expect int, error: %v", key, env, err)
		return def
	}
	return v
}

func MustGetEnv(key string) string {
	env, exist := os.LookupEnv(key)
	if !exist {
		panic(fmt.Sprintf("Env %s is not set", key))
	}
	return env
}
