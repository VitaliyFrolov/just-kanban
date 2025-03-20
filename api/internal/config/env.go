package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"path/filepath"
	"strings"

	"os"
)

// Env defines dictionary of env variables app uses
type Env struct {
	// JWTSecret is secret string for working with jwt encryption
	JWTSecret string
	// ServerPort is port of app hosting
	ServerPort string
	// ServerHost is hostname of app hosting
	ServerHost string
	// DBHost is hostname of database hosting
	DBHost string
	// DBPort is port of database hosting
	DBPort string
	// DBUser is username for working with database
	DBUser string
	// DBPassword is password for working with database
	DBPassword string
	// DBName is name of database app works with
	DBName string
}

func loadEnvFile() {
	for dirOutStepsCount := 0; dirOutStepsCount < 3; dirOutStepsCount++ {
		envFilePath := filepath.Join(strings.Repeat("../", dirOutStepsCount), ".env")
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		fmt.Println(exPath)
		log.Println("Searching for env file at path", exPath, envFilePath)

		loadErr := godotenv.Load(envFilePath)
		if loadErr == nil {
			return
		}
	}
	panic(".env file not found, place it to project's root")
}

// NewEnv loads env variables from .env file and returns structure with those fields
func NewEnv() *Env {
	return &Env{
		JWTSecret:  os.Getenv("JWT_SECRET"),
		ServerPort: os.Getenv("SERVER_PORT"),
		ServerHost: os.Getenv("SERVER_HOST"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}
}
