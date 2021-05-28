package config

import (
	"github.com/bearname/videohost/pkg/common/util"
)

type Config struct {
	MessageBrokerAddress string
	VideoServerAddress   string
	AuthServerAddress    string
	RedisAddress         string
	RedisPassword        string
}

func ParseConfig() *Config {
	messageBrokerAddress := util.ParseEnvString("MESSAGE_BROKER_ADDRESS", "amqp://guest:guest@localhost:5672/")
	videoServerAddress := util.ParseEnvString("VIDEO_SERVER_ADDRESS", "http://localhost:8000")
	authServerAddress := util.ParseEnvString("AUTH_SERVER_ADDRESS", "http://localhost:8001")
	redisAddress := util.ParseEnvString("REDIS_ADDRESS", "localhost:6379")
	redisPassword := util.ParseEnvString("REDIS_PASSWORD", "")

	return &Config{
		messageBrokerAddress,
		videoServerAddress,
		authServerAddress,
		redisAddress,
		redisPassword,
	}
}
