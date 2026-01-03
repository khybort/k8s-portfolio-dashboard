package config

import (
	"os"
	"time"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	LogLevel string
	Seeder   SeederConfig
}

type SeederConfig struct {
	AdminEmail    string
	AdminPassword string
	AdminName     string
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret         string
	AccessExpiry   time.Duration
	RefreshExpiry  time.Duration
}

func Load() (*Config, error) {
	// Determine environment
	env := getEnv("ENV", "development")
	
	// Load environment-specific .env file
	if env == "development" {
		// Try .env.dev first, then .env
		_ = godotenv.Load(".env.dev")
		_ = godotenv.Load(".env")
	} else if env == "production" {
		_ = godotenv.Load(".env.prod")
		_ = godotenv.Load(".env")
	} else {
		// Fallback to .env
		_ = godotenv.Load()
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("AUTH_SERVICE_PORT", "8081"),
			Host: getEnv("AUTH_SERVICE_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "portfolio"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("AUTH_DB_NAME", "auth_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "your-secret-key"),
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 7 * 24 * time.Hour,
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Seeder: SeederConfig{
			AdminEmail:    getEnv("ADMIN_EMAIL", "admin@portfolio.com"),
			AdminPassword: getEnv("ADMIN_PASSWORD", "Admin123!"),
			AdminName:     getEnv("ADMIN_NAME", "Admin User"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *DatabaseConfig) DSN() string {
	return "host=" + c.Host + " port=" + c.Port + " user=" + c.User + " password=" + c.Password + " dbname=" + c.DBName + " sslmode=" + c.SSLMode
}

