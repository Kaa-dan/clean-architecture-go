package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment    string
	Port           string
	DatabaseURL    string
	DatabaseName   string
	JWTSecret      string
	JWTExpiryHours int
	LogLevel       string
	RateLimitRPM   int
	BCryptCost     int
}

func Load() *Config {

	// Load .env file if exists
	godotenv.Load()

	jwtExpiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	rateLimitRPM, _ := strconv.Atoi(getEnv("RATE_LIMIT_RPM", "60"))
	bcryptCost, _ := strconv.Atoi(getEnv("BCRYPT_COST", "12"))

	return &Config{
		Environment:    getEnv("ENVIRONMENT", "development"),
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABSE_URL", "mongodb://localhost:27017"),
		DatabaseName:   getEnv("DATABASE_NAME", "userapi"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-this"),
		JWTExpiryHours: jwtExpiryHours,
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		RateLimitRPM:   rateLimitRPM,
		BCryptCost:     bcryptCost,
	}

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
