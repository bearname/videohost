package config

import (
	"github.com/bearname/videohost/pkg/common/util"
	"strconv"
)

type Config struct {
	Port                 int
	DbName               string
	DbAddress            string
	DbUser               string
	DbPassword           string
	MessageBrokerAddress string
	AuthServerAddress    string
	RedisAddress         string
	RedisPassword        string
}

func ParseConfig() (*Config, error) {
	var err error
	port := util.ParseEnvString("PORT", "8000")
	dbName := util.ParseEnvString("DATABASE_NAME", "video")
	dbAddress := util.ParseEnvString("DATABASE_ADDRESS", "")
	dbUser := util.ParseEnvString("DATABASE_USER", "root")
	dbPassword := util.ParseEnvString("DATABASE_PASSWORD", "123")
	messageBrokerAddress := util.ParseEnvString("MESSAGE_BROKER_ADDRESS", "amqp://guest:guest@localhost:5672/")
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
		messageBrokerAddress,
		authServerAddress,
		redisAddress,
		redisPassword,
	}, nil
}
