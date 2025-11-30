package config

import (
	"confam-api/internal/database"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   database.Config
	Port string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	dbConfig := database.Config{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DB:   dbConfig,
		Port: port,
	}, nil
}
