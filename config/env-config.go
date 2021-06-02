package config

import (
	"github.com/joho/godotenv"
)

func SetupEnv() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}

	// os.Getenv()获取配置
}
