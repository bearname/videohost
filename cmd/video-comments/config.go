package main

import (
	"github.com/bearname/videohost/internal/common/util"
	"strconv"
)

type Config struct {
	Port               int
	DbName             string
	DbAddress          string
	DbUser             string
	DbPassword         string
	VideoServerAddress string
	AuthServerAddress  string
	RedisAddress       string
	RedisPassword      string
}

func ParseConfig() (*Config, error) {
	var err error
	port := util.ParseEnvString("PORT", "8010")
	dbName := util.ParseEnvString("DATABASE_NAME", "video")
	dbAddress := util.ParseEnvString("DATABASE_ADDRESS", "")
	dbUser := util.ParseEnvString("DATABASE_USER", "root")
	dbPassword := util.ParseEnvString("DATABASE_PASSWORD", "123")
	videoServerAddress := util.ParseEnvString("VIDEO_SERVER_ADDRESS", "http://localhost:8000")
	authServerAddress := util.ParseEnvString("AUTH_SERVER_ADDRESS", "http://localhost:8001")
	redisAddress := util.ParseEnvString("REDIS_ADDRESS", "localhost:6379")
	redisPassword := util.ParseEnvString("REDIS_PASSWORD", "")

	atoi, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	return &Config{
		atoi,
		dbName,
		dbAddress,
		dbUser,
		dbPassword,
		videoServerAddress,
		authServerAddress,
		redisAddress,
		redisPassword,
	}, nil
}
