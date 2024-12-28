package config

import (
	"os"
	"strconv"
)

type Config struct {
	JWTExpirationInSeconds int64
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP_SECONDS", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "SECRET"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}

	return fallback
}
