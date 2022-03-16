package config

import (
	"os"
	"time"
)

type Setting struct {
	LAToken       string
	LAHTTPHost    string
	LAHTTPTimeout time.Duration
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, err := time.ParseDuration(getEnv(key, "")); err == nil {
		return value
	}
	return fallback
}

func NewSetting() Setting {
	return Setting{
		LAToken:       getEnv("LA_TOKEN", "test-token"),
		LAHTTPHost:    getEnv("LA_HTTP_HOST", "https://lostark.game.onstove.com"),
		LAHTTPTimeout: getEnvDuration("LA_HTTP_TIMEOUT", 10*time.Second),
	}
}
