package env

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Error loading .env file")
	}
}

func GetEnvString(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return ""
}

func GetTimeDuration(key string) time.Duration {
	if val, ok := os.LookupEnv(key); ok {
		intVal, err := time.ParseDuration(val)
		if err != nil {
			log.Printf("Error converting env var %s to int", key)
			return 0
		}

		return intVal
	}

	return 0
}
