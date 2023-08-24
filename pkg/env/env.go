package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVariableByKey(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
