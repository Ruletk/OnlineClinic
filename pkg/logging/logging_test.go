package logging_test

import (
	"bytes"
	"encoding/json"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/sirupsen/logrus"
	"testing"

	"github.com/Ruletk/OnlineClinic/pkg/config"
	logconf "github.com/Ruletk/OnlineClinic/pkg/config/logging"
	"github.com/stretchr/testify/suite"
)

// TestSuite содержит общие зависимости для всех тестов.
type LoggingTestSuite struct {
	suite.Suite
	buf *bytes.Buffer
}

// Настраивает тестовое окружение перед каждым тестом.
func (s *LoggingTestSuite) SetupTest() {
	s.buf = new(bytes.Buffer)

	// Minimal configuration for the logger
	cfg := config.Config{
		Logger: config.LoggerConfig{
			Format:       logconf.JSON,
			EnableCaller: true,
			Output:       s.buf,
		},
	}

	logging.InitLogger(cfg)
}

// Возвращает текущий вывод логов в виде строки.
func (s *LoggingTestSuite) getLogOutput() string {
	return s.buf.String()
}

// Разбирает JSON-лог в map.
func (s *LoggingTestSuite) parseLogOutput() map[string]interface{} {
	output := s.getLogOutput()
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	s.Require().NoError(err, "Failed to unmarshal log output")
	return logEntry
}

// Запускает Test Suite.
func TestLoggingSuite(t *testing.T) {
	suite.Run(t, new(LoggingTestSuite))
}

// Проверяет базовое логирование (уровень, сообщение, caller).
func (s *LoggingTestSuite) TestBasicLogging() {
	logging.Logger.Info("Test message")

	logEntry := s.parseLogOutput()

	// Проверяем поля лога
	s.Equal("Test message", logEntry["message"], "Log message mismatch")
	s.Equal("info", logEntry["level"], "Log level mismatch")
	s.Contains(logEntry, "caller", "Caller field should exist")
}

// Проверяет обработку невалидного уровня логирования.
func (s *LoggingTestSuite) TestInvalidLogLevel() {
	// Меняем уровень на невалидный (в SetupTest стоит "info" по умолчанию)
	cfg := config.Config{
		Logger: config.LoggerConfig{
			Level:  "invalid",
			Output: s.buf,
		},
	}

	logging.InitLogger(cfg)
	s.Equal(logrus.InfoLevel, logging.Logger.GetLevel(), "Should fall back to 'info' on invalid level")
}

// Проверяет наличие поля 'loggerName' (если оно должно быть).
func (s *LoggingTestSuite) TestLoggerNameField() {
	logging.Logger.Info("Test with logger name")

	logEntry := s.parseLogOutput()
	s.Contains(logEntry, "loggerName", "Logger name should be in log output")
	s.Equal("default", logEntry["loggerName"], "Default logger name should be 'default'")
}
