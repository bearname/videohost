package config

import (
	"github.com/bearname/videohost/pkg/common/util"
	"strconv"
)

type Config struct {
	Port       int
	DbName     string
	DbAddress  string
	DbUser     string
	DbPassword string
}

func ParseConfig() (*Config, error) {
	var err error
	port := util.ParseEnvString("PORT", "8001")
	dbName := util.ParseEnvString("DATABASE_NAME", "video")
	dbAddress := util.ParseEnvString("DATABASE_ADDRESS", "")
	dbUser := util.ParseEnvString("DATABASE_USER", "root")
	dbPassword := util.ParseEnvString("DATABASE_PASSWORD", "123")

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
	}, nil
}
