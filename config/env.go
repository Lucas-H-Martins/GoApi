package config

import (
	"fmt"
	"goapi/logger"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Dev  Environment = "dev"
	Stag Environment = "stag"
	Prod Environment = "prod"
)

type Config struct {
	Environment Environment
	Server      ServerConfig
	Database    DBConfig
	Logger      LoggerConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type LoggerConfig struct {
	Level     logger.LogLevel
	UseColors bool
}

// LoadConfig loads configuration based on the environment
func LoadConfig() (*Config, error) {
	env := getEnvironment()

	// Initialize logger first
	logConfig, err := initializeLogger(env)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	// Load environment file from main folder
	envFile := fmt.Sprintf("./env.%s", env)

	if err := godotenv.Load(envFile); err != nil {
		logger.Warn("No %s file found in main folder, using default environment variables", envFile)
	}

	config := &Config{
		Environment: env,
		Logger:      logConfig,
		Server: ServerConfig{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
			Host: getEnvOrDefault("SERVER_HOST", "localhost"),
		},
		Database: DBConfig{
			Host:     resolveSecret(getEnvOrDefault("DB_HOST", "localhost")),
			Port:     resolveSecret(getEnvOrDefault("DB_PORT", "5432")),
			User:     resolveSecret(getEnvOrDefault("DB_USER", "postgres")),
			Password: resolveSecret(getEnvOrDefault("DB_PASSWORD", "postgres")),
			DBName:   resolveSecret(getEnvOrDefault("DB_NAME", "goapi_db")),
			SSLMode:  resolveSecret(getEnvOrDefault("DB_SSL_MODE", "disable")),
		},
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	// Log configuration loaded
	logger.Info("Configuration loaded for environment: %s", config.Environment)
	logger.Debug("Server configuration - Host: %s, Port: %s", config.Server.Host, config.Server.Port)
	logger.Debug("Database configuration - Host: %s, Port: %s, Database: %s",
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
	)

	return config, nil
}

func initializeLogger(env Environment) (LoggerConfig, error) {
	// Default log levels per environment
	defaultLevel := map[Environment]string{
		Dev:  "DEBUG",
		Stag: "INFO",
		Prod: "WARN",
	}

	// Get log level from environment or use default
	levelStr := getEnvOrDefault("LOG_LEVEL", defaultLevel[env])
	level, err := logger.ParseLevel(levelStr)
	if err != nil {
		return LoggerConfig{}, err
	}

	// Colors enabled by default for local and dev
	useColors := env != Prod
	if colorStr := os.Getenv("LOG_USE_COLORS"); colorStr != "" {
		useColors = strings.ToLower(colorStr) == "true"
	}

	// Initialize the logger
	logger.InitLogger(level, useColors)

	return LoggerConfig{
		Level:     level,
		UseColors: useColors,
	}, nil
}

func getEnvironment() Environment {
	env := os.Getenv("GO_ENV")
	switch strings.ToLower(env) {
	case "dev":
		return Dev
	case "stag":
		return Stag
	case "prod":
		return Prod
	default:
		return Dev
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// resolveSecret handles secret resolution for values starting with "!"
func resolveSecret(value string) string {
	if !strings.HasPrefix(value, "!") {
		return value
	}

	// Remove the "!" prefix
	secretPath := strings.TrimPrefix(value, "!")

	// In a real application, you would implement secret management here
	// For example, reading from Vault, AWS Secrets Manager, or a local secrets file
	secretContent, err := os.ReadFile(secretPath)
	if err != nil {
		logger.Error("Error reading secret from %s: %v", secretPath, err)
		return ""
	}

	return strings.TrimSpace(string(secretContent))
}

func validateConfig(config *Config) error {
	if config.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if config.Server.Host == "" {
		return fmt.Errorf("server host is required")
	}

	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Port == "" {
		return fmt.Errorf("database port is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Database.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	if config.Database.SSLMode == "" {
		return fmt.Errorf("database ssl mode is required")
	}

	// No need to validate UseColors as it's a bool with default value

	// Add more validation as needed
	return nil
}
