package config

import (
	"os"
)

type Config struct {
	Database          DatabaseConfig
	JWTSecret         string
	Port              string
	Environment       string
	GraphQLEndpoint   string
	GraphQLPlayground bool
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "admin"),
			DBName:   getEnv("DB_NAME", "user_team_management"),
		},
		JWTSecret:         getEnv("JWT_SECRET", "seta-go"),
		Port:              getEnv("PORT", "8080"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		GraphQLEndpoint:   getEnv("GRAPHQL_ENDPOINT", "/graphql"),
		GraphQLPlayground: getEnv("GRAPHQL_PLAYGROUND", "true") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
