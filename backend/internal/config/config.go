package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration values loaded from environment variables.
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Server
	ServerPort string

	// JWT
	JWTSecret           string
	JWTAccessExpiryMin  int
	JWTRefreshExpiryDay int

	// Security
	VerificationSecret string

	// SMTP (Email)
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	// Application
	AppURL string
}

// Load reads the .env file (if present) and populates a Config struct.
// Missing values fall back to sensible development defaults.
func Load() (*Config, error) {
	// Load .env file — ignore error if the file doesn't exist (e.g. in Docker)
	_ = godotenv.Load()

	accessExpiry, err := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRY_MINUTES", "60"))
	if err != nil {
		accessExpiry = 60
	}

	refreshExpiry, err := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_DAYS", "7"))
	if err != nil {
		refreshExpiry = 7
	}

	verificationSecret := os.Getenv("VERIFICATION_SECRET")
	if verificationSecret == "" {
		return nil, fmt.Errorf("VERIFICATION_SECRET environment variable is strictly required")
	}

	return &Config{
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "postgres"),
		DBName:              getEnv("DB_NAME", "jewellery_billing"),
		DBSSLMode:           getEnv("DB_SSLMODE", "disable"),
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		JWTSecret:           getEnv("JWT_SECRET", "change-me-in-production"),
		JWTAccessExpiryMin:  accessExpiry,
		JWTRefreshExpiryDay: refreshExpiry,
		VerificationSecret:  verificationSecret,
		SMTPHost:            getEnv("SMTP_HOST", ""),
		SMTPPort:            getEnv("SMTP_PORT", "587"),
		SMTPUser:            getEnv("SMTP_USER", ""),
		SMTPPassword:        getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:            getEnv("SMTP_FROM", "noreply@jewellery-billing.com"),
		AppURL:              getEnv("APP_URL", "http://localhost:5173"),
	}, nil
}

// DatabaseURL returns a PostgreSQL connection string.
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}

// IsSMTPConfigured returns true if SMTP settings are provided.
func (c *Config) IsSMTPConfigured() bool {
	return c.SMTPHost != "" && c.SMTPUser != ""
}

// getEnv reads an environment variable or returns a fallback default.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
