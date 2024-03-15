package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	env := os.Getenv("ENV")
	if err != nil && env == "" {
		log.Fatal("error loading .env file")
	}
}
