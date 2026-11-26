package config

import (
	"fmt"
	"log"
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

	if err := godotenv.Load(); err != nil {
		log.Printf("Note: .env file not found, reading from system environment variable")
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
		return "", fmt.Errorf("missing env variables: %s (check ur .env file)", key)
	}
	return val, nil
}
