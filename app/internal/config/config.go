package config

import (
	"log"
	"os"
)

type Config struct {
	Env            string
	Port           string
	DBPath         string
	TelegramToken  string
	WebhookURL     string
	TLSCert        string
	TLSKey         string
}

func MustLoad() *Config {
	return &Config{
		Env:           getEnv("ENV", "local"),
		Port:          getEnv("PORT", "8080"),
		DBPath:        getEnv("DB_PATH", "./db/db.db"),
		TelegramToken: mustEnv("TELEGRAM_BOT_TOKEN"),
		WebhookURL:    getEnv("WEBHOOK_URL", ""),
		TLSCert:       getEnv("TLS_CERT", ""),
		TLSKey:        getEnv("TLS_KEY", ""),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("env %s is required", key)
	}
	return val
}
