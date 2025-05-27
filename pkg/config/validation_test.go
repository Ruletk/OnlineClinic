package config_test

import (
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/config/logging"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DatabaseConfigValidationTestSuite struct {
	suite.Suite
	DatabaseConfig *config.DatabaseConfig
}

func TestDatabaseConfigValidation(t *testing.T) {
	suite.Run(t, new(DatabaseConfigValidationTestSuite))
}

func (suite *DatabaseConfigValidationTestSuite) SetupTest() {
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

func (suite *DatabaseConfigValidationTestSuite) TestDatabaseValidation() {
	err := suite.DatabaseConfig.Validate()
	suite.NoError(err, "Expected no error for valid config")
}

func (suite *DatabaseConfigValidationTestSuite) TestDatabaseValidation_EmptyHost() {
	suite.DatabaseConfig.Host = ""
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for empty host")
	suite.Contains(err.Error(), "database host cannot be empty", "Expected error to contain 'database host cannot be empty' message")
}

func (suite *DatabaseConfigValidationTestSuite) TestDatabaseValidation_InvalidPort() {
	suite.DatabaseConfig.Port = 0
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for invalid port")
	suite.Contains(err.Error(), "database port must be greater than 0", "Expected error to contain 'database port must be greater than 0' message")
}

func (suite *DatabaseConfigValidationTestSuite) TestDatabaseValidation_PortTooHigh() {
	suite.DatabaseConfig.Port = 70000
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for port greater than 65535")
	suite.Contains(err.Error(), "database port must be less than 65536", "Expected error to contain 'database port must be less than 65536' message")
}

func (suite *DatabaseConfigValidationTestSuite) TestDatabaseValidation_EmptyUser() {
	suite.DatabaseConfig.User = ""
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for empty user")
	suite.Contains(err.Error(), "database user cannot be empty", "Expected error to contain 'database user cannot be empty' message")
}

func (suite *DatabaseConfigValidationTestSuite) TestDatabaseValidation_EmptyPassword() {
	suite.DatabaseConfig.Password = ""
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for empty password")
	suite.Contains(err.Error(), "database password cannot be empty", "Expected error to contain 'database password cannot be empty' message")
}

func (suite *DatabaseConfigValidationTestSuite) TestDatabaseValidation_EmptyDBName() {
	suite.DatabaseConfig.DBName = ""
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for empty dbname")
	suite.Contains(err.Error(), "database name cannot be empty", "Expected error to contain 'database name cannot be empty' message")
}

func (suite *DatabaseConfigValidationTestSuite) TestDatabaseValidation_EmptySSLMode() {
	suite.DatabaseConfig.SSLMode = ""
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for empty sslmode")
	suite.Contains(err.Error(), "database sslmode cannot be empty", "Expected error to contain 'database sslmode cannot be empty' message")
}

func (suite *DatabaseConfigValidationTestSuite) TestDatabaseValidation_EmptyCharset() {
	suite.DatabaseConfig.Charset = ""
	err := suite.DatabaseConfig.Validate()
	suite.Error(err, "Expected error for empty charset")
	suite.Contains(err.Error(), "database charset cannot be empty", "Expected error to contain 'database charset cannot be empty' message")
}

func (suite *DatabaseConfigValidationTestSuite) TestDatabaseValidation_AllErrors() {
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

type BackendConfigValidationTestSuite struct {
	suite.Suite
	ValidConfig *config.BackendConfig
}

func TestBackendConfigValidation(t *testing.T) {
	suite.Run(t, new(BackendConfigValidationTestSuite))
}

func (suite *BackendConfigValidationTestSuite) SetupTest() {
	suite.ValidConfig = &config.BackendConfig{
		ListenAddress: "0.0.0.0",
		ListenPort:    8080,
	}
}

func (suite *BackendConfigValidationTestSuite) TestValidListenAddress_SuccessIP4() {
	suite.ValidConfig.ListenAddress = "127.0.0.1"
	err := suite.ValidConfig.Validate()
	suite.NoError(err, "Expected no error for valid IPv4 address")
}

func (suite *BackendConfigValidationTestSuite) TestValidListenAddress_SuccessIP6() {
	suite.ValidConfig.ListenAddress = "::1"
	err := suite.ValidConfig.Validate()
	suite.NoError(err, "Expected no error for valid IPv6 address")
}

func (suite *BackendConfigValidationTestSuite) TestValidListenAddress_SuccessIP6Full() {
	suite.ValidConfig.ListenAddress = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	err := suite.ValidConfig.Validate()
	suite.NoError(err, "Expected no error for valid full IPv6 address")
}

func (suite *BackendConfigValidationTestSuite) TestValidListenAddress_SuccessDomain() {
	suite.ValidConfig.ListenAddress = "example.com"
	err := suite.ValidConfig.Validate()
	suite.NoError(err, "Expected no error for valid domain name")
}

func (suite *BackendConfigValidationTestSuite) TestValidListenAddress_SuccessDomainLevel1() {
	suite.ValidConfig.ListenAddress = "database"
	err := suite.ValidConfig.Validate()
	suite.NoError(err, "Expected no error for valid domain name")
}

func (suite *BackendConfigValidationTestSuite) TestValidListenAddress_SuccessDomainLevel3() {
	suite.ValidConfig.ListenAddress = "database.example.com"
	err := suite.ValidConfig.Validate()
	suite.NoError(err, "Expected no error for valid domain name")
}

func (suite *BackendConfigValidationTestSuite) TestInvalidListenAddress_Failure() {
	suite.ValidConfig.ListenAddress = ""
	err := suite.ValidConfig.Validate()
	suite.Error(err, "Expected error for empty listen address")
	suite.Contains(err.Error(), "backend listen address cannot be empty", "Expected error to contain 'backend listen address cannot be empty' message")
}

func (suite *BackendConfigValidationTestSuite) TestInvalidListenAddress_FailureSpace() {
	suite.ValidConfig.ListenAddress = "invalid address"
	err := suite.ValidConfig.Validate()
	suite.Error(err, "Expected error for invalid listen address")
	suite.Contains(err.Error(), "backend listen address cannot contain spaces", "Expected error to contain 'backend listen address cannot contain spaces' message")
}

func (suite *BackendConfigValidationTestSuite) TestInvalidPort_Failure() {
	suite.ValidConfig.ListenPort = 0
	err := suite.ValidConfig.Validate()
	suite.Error(err, "Expected error for invalid port")
	suite.Contains(err.Error(), "backend listen port must be greater than 0", "Expected error to contain 'backend listen port must be greater than 0' message")
}

func (suite *BackendConfigValidationTestSuite) TestInvalidPort_FailureTooHigh() {
	suite.ValidConfig.ListenPort = 70000
	err := suite.ValidConfig.Validate()
	suite.Error(err, "Expected error for port greater than 65535")
	suite.Contains(err.Error(), "backend listen port must be less than 65536", "Expected error to contain 'backend listen port must be less than 65536' message")
}

func (suite *BackendConfigValidationTestSuite) TestInvalidPort_FailureNegative() {
	suite.ValidConfig.ListenPort = -1
	err := suite.ValidConfig.Validate()
	suite.Error(err, "Expected error for negative port")
	suite.Contains(err.Error(), "backend listen port must be greater than 0", "Expected error to contain 'backend listen port must be greater than 0' message")
}

func (suite *BackendConfigValidationTestSuite) TestValidPort_Success() {
	suite.ValidConfig.ListenPort = 8080
	err := suite.ValidConfig.Validate()
	suite.NoError(err, "Expected no error for valid port")
}

func (suite *BackendConfigValidationTestSuite) TestInvalidPort_SuccessTooHigh() {
	suite.ValidConfig.ListenPort = 65535
	err := suite.ValidConfig.Validate()
	suite.NoError(err, "Expected no error for valid port")
}

type NatsConfigValidationTestSuite struct {
	suite.Suite
	ValidConfig *config.NatsConfig
}

func TestNatsConfigValidation(t *testing.T) {
	suite.Run(t, new(NatsConfigValidationTestSuite))
}

func (suite *NatsConfigValidationTestSuite) SetupTest() {
	suite.ValidConfig = &config.NatsConfig{
		Url: "nats://localhost:4222",
	}
}

func (suite *NatsConfigValidationTestSuite) TestValidUrl() {
	cfg := *suite.ValidConfig
	cfg.Url = "nats://localhost:4222"
	suite.NoError(cfg.Validate())
}

func (suite *NatsConfigValidationTestSuite) TestEmptyUrl() {
	cfg := *suite.ValidConfig
	cfg.Url = ""
	suite.Error(cfg.Validate())
}
