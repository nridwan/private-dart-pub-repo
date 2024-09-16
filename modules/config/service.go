package config

import "os"

type ConfigService interface {
	Getenv(key string, fallback string) string
}

// impl `ConfigService` start

func (*ConfigModule) Getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// impl `ConfigService` end
