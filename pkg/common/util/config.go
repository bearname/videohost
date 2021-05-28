package util

import (
	"os"
)

func ParseEnvString(key string, defaultValue string) string {
	str, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return str
}
