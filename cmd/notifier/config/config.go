package config

import (
	"github.com/bearname/videohost/pkg/common/util"
)

type Config struct {
	SendBlueAddress      string
	SendInBlueApiKey     string
	MessageBrokerAddress string
	AuthServerAddress    string
}

func ParseConfig() *Config {
	sendBlueAddress := util.ParseEnvString("SEND_BLUE_ADDRESS", "https://api.sendinblue.com")
	messageBrokerAddress := util.ParseEnvString("MESSAGE_BROKER_ADDRESS", "amqp://guest:guest@localhost:5672/")
	sendInBlueApiKey := util.ParseEnvString("SEND_BLUE_API_KEY", "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")
	authServerAddress := util.ParseEnvString("AUTH_SERVER_ADDRESS", "http://localhost:8001")

	return &Config{
		sendBlueAddress,
		sendInBlueApiKey,
		messageBrokerAddress,
		authServerAddress,
	}
}
