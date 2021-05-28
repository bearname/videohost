package config

import "github.com/bearname/videohost/pkg/common/util"

type Config struct {
	DbName            string
	DbAddress         string
	DbUser            string
	DbPassword        string
	AuthServerAddress string
	RedisAddress      string
	RedisPassword     string
}

func ParseConfig() *Config {
	dbName := util.ParseEnvString("DATABASE_NAME", "video")
	dbAddress := util.ParseEnvString("DATABASE_ADDRESS", "")
	dbUser := util.ParseEnvString("DATABASE_USER", "root")
	dbPassword := util.ParseEnvString("DATABASE_PASSWORD", "123")
	messageBrokerAddress := util.ParseEnvString("MESSAGE_BROKER_ADDRESS", "amqp://guest:guest@localhost:5672/")
	redisAddress := util.ParseEnvString("REDIS_ADDRESS", "localhost:6379")
	redisPassword := util.ParseEnvString("REDIS_PASSWORD", "")

	return &Config{
		dbName,
		dbAddress,
		dbUser,
		dbPassword,
		messageBrokerAddress,
		redisAddress,
		redisPassword,
	}
}
