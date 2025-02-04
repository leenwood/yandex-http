package config

import "os"

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Hostname string
	Port     string
}

type DatabaseConfig struct {
	Hostname string
	Port     string
	Username string
	Password string
	Database string
}

func NewConfig() Config {
	return Config{
		App: AppConfig{
			Hostname: getEnv("HOSTNAME", "localhost"),
			Port:     getEnv("PORT", "9000"),
		},
		Database: DatabaseConfig{
			Hostname: getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnv("DATABASE_PORT", "5432"),
			Username: getEnv("DATABASE_USER", "postgres"),
			Password: getEnv("DATABASE_PASS", "postgres"),
			Database: getEnv("DATABASE_NAME", "app_db"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
