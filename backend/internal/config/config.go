package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	Environment string
}

func Load() (*Config, error) { //public fn

	if err := godotenv.Load("../.env"); err != nil {
		return nil, fmt.Errorf("Note: .env file not found, reading from system environment variables")
	}
	port, err := getRequiredEnv("PORT")
	if err != nil {
		return nil, err
	}
	dbURL, err := getRequiredEnv("DATABASE_URL")
	if err != nil {
		return nil, err
	}
	jwtSecret, err := getRequiredEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	env, err := getRequiredEnv("ENV")
	if err != nil {
		return nil, err
	}
	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
		Environment: env,
	}, nil
}
func getRequiredEnv(key string) (string, error) {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" {
		return "", fmt.Errorf("missing env varialbes: %s (check ur .env file)", key)
	}
	return val, nil
}
