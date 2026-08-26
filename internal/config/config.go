package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost    string
	DBUser    string
	DBPass    string
	DBName    string
	DBPort    string
	DBSSLMode string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error Load .env")
	}

	return &Config{
		DBHost:    os.Getenv("DBHost"),
		DBUser:    os.Getenv("DBUser"),
		DBPass:    os.Getenv("DBPass"),
		DBName:    os.Getenv("DBName"),
		DBPort:    os.Getenv("DBPort"),
		DBSSLMode: os.Getenv("DBSSLMode"),
	}

}
