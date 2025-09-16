package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort         string
	AppEnv          string
	DatabaseURL     string
	RedisAddr       string
	RedisDB         int
	RedisTTLSeconds int
	ESAddr          string
	ESIndex         string
	Timeout         time.Duration
}

func Load() Config {
	v := viper.New()
	v.SetConfigFile(".env")
	_ = v.ReadInConfig()
	v.AutomaticEnv()

	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("REDIS_TTL_SECONDS", 300)
	v.SetDefault("ES_ADDR", "http://localhost:9200")
	v.SetDefault("ES_INDEX", "posts")

	return Config{
		AppPort:         v.GetString("APP_PORT"),
		AppEnv:          v.GetString("APP_ENV"),
		DatabaseURL:     v.GetString("DATABASE_URL"),
		RedisAddr:       v.GetString("REDIS_ADDR"),
		RedisDB:         v.GetInt("REDIS_DB"),
		RedisTTLSeconds: v.GetInt("REDIS_TTL_SECONDS"),
		ESAddr:          v.GetString("ES_ADDR"),
		ESIndex:         v.GetString("ES_INDEX"),
		Timeout:         5 * time.Second,
	}
}
