package config

import (
	"github.com/Ruletk/OnlineClinic/pkg/config/logging"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type DefaultConfigTestSuite struct {
	suite.Suite
}

func TestDefaultConfig(t *testing.T) {
	suite.Run(t, new(DefaultConfigTestSuite))
}

func (suite *DefaultConfigTestSuite) SetupTest() {
	// clear the environment variables before each test
	os.Clearenv()
}

func (suite *DefaultConfigTestSuite) TestGetDefaultEnv_Success() {
	_ = os.Setenv("APP_PORT", "9000")
	value := GetEnvWithDefault("APP_PORT", "8000")
	suite.Equal("9000", value, "Expected APP_PORT to be 9000")
}

func (suite *DefaultConfigTestSuite) TestGetDefaultEnv_NoEnv() {
	value := GetEnvWithDefault("NON_EXISTENT_ENV", "default_value")
	suite.Equal("default_value", value, "Expected default value to be returned")
}

func (suite *DefaultConfigTestSuite) TestGetDefaultEnv_EmptyValue() {
	_ = os.Setenv("EMPTY_ENV", "")
	value := GetEnvWithDefault("EMPTY_ENV", "default_value")
	suite.Equal("default_value", value, "Expected default value to be returned for empty env")
}

func (suite *DefaultConfigTestSuite) TestDefaultConfig_NoEnv() {
	config, err := GetDefaultConfiguration()
	suite.NoError(err, "Expected no error when getting default configuration")
	suite.Equal(8080, config.Backend.ListenPort, "Expected default APP_PORT to be 8080")
	suite.Equal("0.0.0.0", config.Backend.ListenAddress, "Expected default APP_HOST to be 0.0.0.0")
	suite.Equal("localhost", config.Database.Host, "Expected default DB_HOST to be localhost")
	suite.Equal(5432, config.Database.Port, "Expected default DB_PORT to be 5432")
	suite.Equal("postgres", config.Database.User, "Expected default DB_USER to be postgres")
	suite.Equal("postgres", config.Database.Password, "Expected default DB_PASSWORD to be postgres")
	suite.Equal("postgres", config.Database.DBName, "Expected default DB_NAME to be postgres")
	suite.Equal("disable", config.Database.SSLMode, "Expected default DB_SSL_MODE to be disable")
	suite.Equal("utf8", config.Database.Charset, "Expected default DB_CHARSET to be utf8")
	suite.Equal("default", config.Logger.LoggerName, "Expected default LOGGER_NAME to be default")
	suite.Equal(logging.Info, config.Logger.Level, "Expected default LOGGER_LEVEL to be info")
	suite.Equal(logging.JSON, config.Logger.Format, "Expected default LOGGER_FORMAT to be json")
	suite.Equal(true, config.Logger.EnableCaller, "Expected default LOGGER_ENABLE_CALLER to be true")
}

func (suite *DefaultConfigTestSuite) TestDefaultConfig_AppPort_Success() {
	_ = os.Setenv("APP_PORT", "9000")
	config, err := GetDefaultConfiguration()
	suite.NoError(err, "Expected no error when getting default configuration")
	suite.Equal(9000, config.Backend.ListenPort, "Expected APP_PORT to be 9000")
}

func (suite *DefaultConfigTestSuite) TestDefaultConfig_AppHost_Success() {
	_ = os.Setenv("APP_HOST", "127.0.0.1")
	config, err := GetDefaultConfiguration()
	suite.NoError(err, "Expected no error when getting default configuration")
	suite.Equal("127.0.0.1", config.Backend.ListenAddress, "Expected APP_HOST to be 127.0.0.1")
}

func (suite *DefaultConfigTestSuite) TestDefaultConfig_DBPort_Invalid() {
	_ = os.Setenv("DB_PORT", "invalid")
	_, err := GetDefaultConfiguration()
	suite.Error(err, "Expected error when DB_PORT is invalid")
}

func (suite *DefaultConfigTestSuite) TestDefaultConfig_AppPort_Invalid() {
	_ = os.Setenv("APP_PORT", "invalid")
	_, err := GetDefaultConfiguration()
	suite.Error(err, "Expected error when APP_PORT is invalid")
}

func (suite *DefaultConfigTestSuite) TestDefaultConfig_EnableCaller_Invalid() {
	_ = os.Setenv("LOGGER_ENABLE_CALLER", "invalid")
	_, err := GetDefaultConfiguration()
	suite.Error(err, "Expected error when LOGGER_ENABLE_CALLER is invalid")
}

func (suite *DefaultConfigTestSuite) TestDefaultConfig_BackendConfig_Invalid() {
	_ = os.Setenv("APP_HOST", "invalid address")
	_, err := GetDefaultConfiguration()
	suite.Error(err, "Expected error when APP_HOST is invalid")
}

func (suite *DefaultConfigTestSuite) TestDefaultConfig_LoggerConfig_Invalid() {
	_ = os.Setenv("LOGGER_LEVEL", "invalid")
	_, err := GetDefaultConfiguration()
	suite.Error(err, "Expected error when LOGGER_LEVEL is invalid")
}

func (suite *DefaultConfigTestSuite) TestDefaultConfig_DatabaseConfig_Invalid() {
	_ = os.Setenv("DB_PORT", "70000")
	_, err := GetDefaultConfiguration()
	suite.Error(err, "Expected error when DB_HOST is invalid")
}
