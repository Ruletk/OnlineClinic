package config_test

import (
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/config/logging"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConfigValidationTestSuite struct {
	suite.Suite
	DatabaseConfig *config.DatabaseConfig
}

func TestConfigValidation(t *testing.T) {
	suite.Run(t, new(ConfigValidationTestSuite))
}

func (suite *ConfigValidationTestSuite) SetupTest() {
	suite.DatabaseConfig = &config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "test_db",
		SSLMode:  "disable",
		Charset:  "utf8",
	}
}

func (suite *ConfigValidationTestSuite) TestDatabaseValidation() {
	err := suite.DatabaseConfig.Validate()
	suite.NoError(err, "Expected no error for valid config")
}

func (suite *ConfigValidationTestSuite) TestDatabaseValidation_EmptyHost() {
	suite.DatabaseConfig.Host = ""
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for empty host")
	suite.Contains(err.Error(), "database host cannot be empty", "Expected error to contain 'database host cannot be empty' message")
}

func (suite *ConfigValidationTestSuite) TestDatabaseValidation_InvalidPort() {
	suite.DatabaseConfig.Port = 0
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for invalid port")
	suite.Contains(err.Error(), "database port must be greater than 0", "Expected error to contain 'database port must be greater than 0' message")
}

func (suite *ConfigValidationTestSuite) TestDatabaseValidation_PortTooHigh() {
	suite.DatabaseConfig.Port = 70000
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for port greater than 65535")
	suite.Contains(err.Error(), "database port must be less than 65536", "Expected error to contain 'database port must be less than 65536' message")
}

func (suite *ConfigValidationTestSuite) TestDatabaseValidation_EmptyUser() {
	suite.DatabaseConfig.User = ""
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for empty user")
	suite.Contains(err.Error(), "database user cannot be empty", "Expected error to contain 'database user cannot be empty' message")
}

func (suite *ConfigValidationTestSuite) TestDatabaseValidation_EmptyPassword() {
	suite.DatabaseConfig.Password = ""
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for empty password")
	suite.Contains(err.Error(), "database password cannot be empty", "Expected error to contain 'database password cannot be empty' message")
}

func (suite *ConfigValidationTestSuite) TestDatabaseValidation_EmptyDBName() {
	suite.DatabaseConfig.DBName = ""
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for empty dbname")
	suite.Contains(err.Error(), "database name cannot be empty", "Expected error to contain 'database name cannot be empty' message")
}

func (suite *ConfigValidationTestSuite) TestDatabaseValidation_EmptySSLMode() {
	suite.DatabaseConfig.SSLMode = ""
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for empty sslmode")
	suite.Contains(err.Error(), "database sslmode cannot be empty", "Expected error to contain 'database sslmode cannot be empty' message")
}

func (suite *ConfigValidationTestSuite) TestDatabaseValidation_EmptyCharset() {
	suite.DatabaseConfig.Charset = ""
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for empty charset")
	suite.Contains(err.Error(), "database charset cannot be empty", "Expected error to contain 'database charset cannot be empty' message")
}

func (suite *ConfigValidationTestSuite) TestDatabaseValidation_AllErrors() {
	suite.DatabaseConfig = &config.DatabaseConfig{}
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for all invalid fields")
	suite.Contains(err.Error(), "database host cannot be empty", "Expected error to contain 'database host cannot be empty' message")
	suite.Contains(err.Error(), "database port must be greater than 0", "Expected error to contain 'database port must be greater than 0' message")
	suite.Contains(err.Error(), "database user cannot be empty", "Expected error to contain 'database user cannot be empty' message")
	suite.Contains(err.Error(), "database password cannot be empty", "Expected error to contain 'database password cannot be empty' message")
	suite.Contains(err.Error(), "database name cannot be empty", "Expected error to contain 'database name cannot be empty' message")
	suite.Contains(err.Error(), "database sslmode cannot be empty", "Expected error to contain 'database sslmode cannot be empty' message")
	suite.Contains(err.Error(), "database charset cannot be empty", "Expected error to contain 'database charset cannot be empty' message")
}

type LoggerConfigValidationTestSuite struct {
	suite.Suite
	ValidConfig *config.LoggerConfig
}

func TestLoggerConfigValidation(t *testing.T) {
	suite.Run(t, new(LoggerConfigValidationTestSuite))
}

func (suite *LoggerConfigValidationTestSuite) SetupTest() {
	suite.ValidConfig = &config.LoggerConfig{
		Level:        "info",
		Format:       "json",
		EnableCaller: true,
		LoggerName:   "test-logger",
	}
}

func (suite *LoggerConfigValidationTestSuite) TestValidLevels() {
	validLevels := []string{"debug", "INFO", "WARN", "error", "FATAL"}
	for _, level := range validLevels {
		cfg := *suite.ValidConfig
		cfg.Level = logging.LoggerLevel(level)
		suite.NoError(cfg.Validate())
	}
}

func (suite *LoggerConfigValidationTestSuite) TestInvalidLevel() {
	cfg := *suite.ValidConfig
	cfg.Level = "invalid_level"
	suite.Error(cfg.Validate())
}

func (suite *LoggerConfigValidationTestSuite) TestEmptyLevel() {
	cfg := *suite.ValidConfig
	cfg.Level = ""
	suite.Error(cfg.Validate())
}

func (suite *LoggerConfigValidationTestSuite) TestValidFormats() {
	validFormats := []string{"json", "TEXT"}
	for _, format := range validFormats {
		cfg := *suite.ValidConfig
		cfg.Format = logging.LoggerFormat(format)
		suite.NoError(cfg.Validate())
	}
}

func (suite *LoggerConfigValidationTestSuite) TestInvalidFormat() {
	cfg := *suite.ValidConfig
	cfg.Format = "yaml"
	suite.Error(cfg.Validate())
}

func (suite *LoggerConfigValidationTestSuite) TestEmptyFormat() {
	cfg := *suite.ValidConfig
	cfg.Format = ""
	suite.Error(cfg.Validate())
}

func (suite *LoggerConfigValidationTestSuite) TestEmptyLoggerName() {
	cfg := *suite.ValidConfig
	cfg.LoggerName = ""
	suite.Error(cfg.Validate())
}
