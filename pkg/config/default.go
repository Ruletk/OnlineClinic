package config

import (
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config/logging"
	"os"
	"strconv"
)

func GetDefaultConfiguration() (*Config, error) {
	appPort := GetEnvWithDefault("APP_PORT", "8080")
	appHost := GetEnvWithDefault("APP_HOST", "0.0.0.0")

	dbHost := GetEnvWithDefault("DB_HOST", "localhost")
	dbPort := GetEnvWithDefault("DB_PORT", "5432")
	dbUser := GetEnvWithDefault("DB_USER", "postgres")
	dbPassword := GetEnvWithDefault("DB_PASSWORD", "postgres")
	dbName := GetEnvWithDefault("DB_NAME", "postgres")
	dbSSLMode := GetEnvWithDefault("DB_SSL_MODE", "disable")
	dbCharset := GetEnvWithDefault("DB_CHARSET", "utf8")

	loggerName := GetEnvWithDefault("LOGGER_NAME", "default")
	loggerLevel := GetEnvWithDefault("LOGGER_LEVEL", "info")
	loggerFormat := GetEnvWithDefault("LOGGER_FORMAT", "json")

	// TODO: Make logger output configurable. For now, it is hardcoded to stdout.
	// This is not important part, but it's good to have.
	// I temporarily commented out the loggerOutput variable to avoid confusion.
	//loggerOutput := getEnvWithDefault("LOGGER_OUTPUT", "stdout")

	loggerEnableCaller := GetEnvWithDefault("LOGGER_ENABLE_CALLER", "true")

	natsUrl := getEnvWithDefault("NATS_URL", "nats://localhost:4222")

	appPortInt, err := strconv.Atoi(appPort)
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT value: %w", err)
	}

	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT value: %w", err)
	}

	loggerEnableCallerBool, err := strconv.ParseBool(loggerEnableCaller)
	if err != nil {
		return nil, fmt.Errorf("invalid LOGGER_ENABLE_CALLER value: %w", err)
	}

	dbConfig := DatabaseConfig{
		Host:     dbHost,
		Port:     dbPortInt,
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
		SSLMode:  dbSSLMode,
		Charset:  dbCharset,
	}

	backendConfig := BackendConfig{
		ListenAddress: appHost,
		ListenPort:    appPortInt,
	}

	loggerConfig := LoggerConfig{
		Level:        logging.LoggerLevel(loggerLevel),
		Format:       logging.LoggerFormat(loggerFormat),
		EnableCaller: loggerEnableCallerBool,
		Output:       os.Stdout,
		TestMode:     false,
		LoggerName:   loggerName,
	}

	natsConfig := NatsConfig{
		Url: natsUrl,
	}

	if err := dbConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	if err := loggerConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid logger configuration: %w", err)
	}

	if err := backendConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid backend configuration: %w", err)
	}

	if err := natsConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid nats configuration: %w", err)
	}

	return &Config{
		Database: dbConfig,
		Backend:  backendConfig,
		Logger:   loggerConfig,
		Nats:     natsConfig,
	}, nil
}

func GetEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}
