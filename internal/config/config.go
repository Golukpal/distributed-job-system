package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName     string
	AppEnv      string
	AppPort     string
	DatabaseURL string

	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no env file found")
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		log.Fatal("invalid redis DB value")
	}

	return &Config{
		AppName:       getEnv("APP_NAME", "job-distributed"),
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       getEnv("APP_PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}
	return value
}
