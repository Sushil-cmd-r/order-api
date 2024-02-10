package application

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	RedisAddr  string
	ServerPort uint16
}

func LoadConfig() Config {
	_ = godotenv.Load()

	cfg := Config{
		RedisAddr:  "localhost:6379",
		ServerPort: 8080,
	}

	if redisAddr, exists := os.LookupEnv("REDIS_ADDR"); exists {
		cfg.RedisAddr = redisAddr
	}

	if sPort, exists := os.LookupEnv("PORT"); exists {
		if port, err := strconv.ParseUint(sPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	return cfg
}
