package config

import "os"

type Config struct {
	Host        string
	Port        string
	LogLevel    string
	AnalyzerURL string
}

func Load() *Config {
	return &Config{
		Host:        getEnv("GATEWAY_HOST", "0.0.0.0"),
		Port:        getEnv("GATEWAY_PORT", "8002"),
		LogLevel:    getEnv("LOG_LEVEL", "INFO"),
		AnalyzerURL: getEnv("ANALYZER_URL", "http://localhost:8001"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
