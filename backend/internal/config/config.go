package config

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	Auth     AuthConfig
	LogLevel string
	Seeder   SeederConfig
}

type SeederConfig struct {
	PortfolioName      string
	PortfolioTitle     string
	PortfolioEmail     string
	PortfolioBio       string
	PortfolioGithub    string
	PortfolioLinkedin  string
	PortfolioPhone     string
	PortfolioTheme     string
	PortfolioLanguage  string
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

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers []string
}

type AuthConfig struct {
	ServiceURL string
	JWTSecret  string
}

func Load() (*Config, error) {
	// Determine environment
	env := getEnv("ENV", "development")
	
	// Load environment-specific .env file
	// Priority: .env.{env} > .env > system env vars
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

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("REDIS_DB", 0)

	viper.AutomaticEnv()

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "portfolio"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "portfolio"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       viper.GetInt("REDIS_DB"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
		},
		Auth: AuthConfig{
			ServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
			JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Seeder: SeederConfig{
			PortfolioName:     getEnv("PORTFOLIO_NAME", "Muhsin Kılıç"),
			PortfolioTitle:    getEnv("PORTFOLIO_TITLE", "Fullstack Developer"),
			PortfolioEmail:    getEnv("PORTFOLIO_EMAIL", "kmuhsinn@gmail.com"),
			PortfolioBio:      getEnv("PORTFOLIO_BIO", "Driving web platform optimizations for enhanced performance and user experience across complex systems. Team mate, improving engineering efficiency through cutting-edge technology and best practices. Experienced in developing scalable AI solutions and automation, streamlining processes and delivering critical features for teams."),
			PortfolioGithub:   getEnv("PORTFOLIO_GITHUB", "https://github.com/kmuhsinn"),
			PortfolioLinkedin: getEnv("PORTFOLIO_LINKEDIN", "https://linkedin.com/in/muhsin-kilic"),
			PortfolioPhone:    getEnv("PORTFOLIO_PHONE", "+905377812189"),
			PortfolioTheme:    getEnv("PORTFOLIO_THEME", "dark"),
			PortfolioLanguage: getEnv("PORTFOLIO_LANGUAGE", "en"),
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
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

