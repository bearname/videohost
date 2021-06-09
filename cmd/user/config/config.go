package config

import (
	"github.com/bearname/videohost/internal/common/util"
	"strconv"
)

type Config struct {
	Port            int
	DbName          string
	DbAddress       string
	DbUser          string
	DbPassword      string
	MaxOpenConns    int
	ConnMaxIdleTime int
}

func ParseConfig() (*Config, error) {
	var err error
	port := util.ParseEnvString("PORT", "8001")
	dbName := util.ParseEnvString("DATABASE_NAME", "video")
	dbAddress := util.ParseEnvString("DATABASE_ADDRESS", "")
	dbUser := util.ParseEnvString("DATABASE_USER", "root")
	dbPassword := util.ParseEnvString("DATABASE_PASSWORD", "123")
	MaxOpenConns := util.ParseEnvString("MAX_OPEN_CONNS", "")
	ConnMaxIdleTime := util.ParseEnvString("CONN_MAX_IDLE_TIME", "")

	atoi, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	maxOpenConns, err := strconv.Atoi(MaxOpenConns)
	if err != nil {
		return nil, err
	}
	connMaxIdleTime, err := strconv.Atoi(ConnMaxIdleTime)
	if err != nil {
		return nil, err
	}

	return &Config{
		atoi,
		dbName,
		dbAddress,
		dbUser,
		dbPassword,
		maxOpenConns,
		connMaxIdleTime,
	}, nil
}
