package logging

import (
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/config/logging"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

const Version = "1.0.0"

var Logger *logrus.Logger

// InitLogger initializes the logger with the provided configuration.
// It sets the log level, format, and output destination based on the configuration.
// It also adds a hook to include the logger name in log entries.
// WARNING: If TestMode is enabled, the logger will ignore Format, EnableCaller, and Output.
func InitLogger(config config.Config) {
	Logger = logrus.New()

	if config.Logger.TestMode {
		Logger.SetLevel(logrus.PanicLevel)
	} else if config.Logger.Format == logging.Text {
		setupLoggerFormatterText(Logger)
	} else if config.Logger.Format == logging.JSON {
		setupLoggerFormatterJSON(Logger)
	} else {
		setupLoggerFormatterText(Logger)
		Logger.Warn("Invalid logger format, defaulting to 'text'")
	}

	if config.Logger.EnableCaller && !config.Logger.TestMode {
		Logger.SetReportCaller(true)
	}

	level, err := logrus.ParseLevel(string(config.Logger.Level))
	if err != nil {
		Logger.Warn("Invalid log level, defaulting to 'info'")
		level = logrus.InfoLevel
	}

	Logger.SetLevel(level)

	if config.Logger.LoggerName != "" {
		Logger.AddHook(&LoggerNameHook{LoggerName: config.Logger.LoggerName})
	} else {
		Logger.AddHook(&LoggerNameHook{LoggerName: "default"})
		Logger.Warn("Logger name is empty, defaulting to 'default'")
	}

	if config.Logger.TestMode {
		Logger.SetOutput(io.Discard)
	} else if config.Logger.Output != nil {
		Logger.SetOutput(config.Logger.Output)
	} else {
		Logger.SetOutput(os.Stdout)
		Logger.Warn("Logger output is nil, defaulting to stdout")
	}
}

// setupLoggerFormatterText sets up the logger to use a text formatter.
// Support function for InitLogger.
func setupLoggerFormatterText(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyFunc:  "caller",
			logrus.FieldKeyMsg:   "message",
		},
	})
}

// setupLoggerFormatterJSON sets up the logger to use a JSON formatter.
// Support function for InitLogger.
func setupLoggerFormatterJSON(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "caller",
		},
	})
}
