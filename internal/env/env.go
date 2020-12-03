package env

import (
	"os"
)

func Get(key string, defaultVal string) string  {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}