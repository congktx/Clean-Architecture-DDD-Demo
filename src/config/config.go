package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string
}

func LoadConfig() *Config {
	err := godotenv.Load("../.env")
	if err != nil {
		if err = godotenv.Load(".env"); err != nil {
			log.Printf("Warning: Error loading .env file, using default environment variables. Error: %v", err)
		}
	}

	return &Config{
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDBName:   os.Getenv("POSTGRES_DB_NAME"),
	}
}
