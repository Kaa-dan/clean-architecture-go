package config

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
