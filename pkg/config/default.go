package config

import (
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config/logging"
	"os"
	"strconv"
)

func GetDefaultConfiguration() (*Config, error) {
	appPort := getEnvWithDefault("APP_PORT", "8080")
	appHost := getEnvWithDefault("APP_HOST", "0.0.0.0")

	dbHost := getEnvWithDefault("DB_HOST", "localhost")
	dbPort := getEnvWithDefault("DB_PORT", "5432")
	dbUser := getEnvWithDefault("DB_USER", "postgres")
	dbPassword := getEnvWithDefault("DB_PASSWORD", "postgres")
	dbName := getEnvWithDefault("DB_NAME", "postgres")
	dbSSLMode := getEnvWithDefault("DB_SSL_MODE", "disable")
	dbCharset := getEnvWithDefault("DB_CHARSET", "utf8")

	loggerName := getEnvWithDefault("LOGGER_NAME", "default")
	loggerLevel := getEnvWithDefault("LOGGER_LEVEL", "info")
	loggerFormat := getEnvWithDefault("LOGGER_FORMAT", "json")
	// TODO: Make logger output configurable. For now, it is hardcoded to stdout.
	//loggerOutput := getEnvWithDefault("LOGGER_OUTPUT", "stdout")
	loggerEnableCaller := getEnvWithDefault("LOGGER_ENABLE_CALLER", "true")

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

	if err := dbConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	if err := loggerConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid logger configuration: %w", err)
	}

	if err := backendConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid backend configuration: %w", err)
	}

	return &Config{
		Database: dbConfig,
		Backend:  backendConfig,
		Logger:   loggerConfig,
	}, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
