package config

import (
	"errors"
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config/logging"
	"strings"
)

// Validate checks if the DatabaseConfig is valid.
// Returns a joined error with all validation failures, or nil if valid.
func (c *DatabaseConfig) Validate() error {
	var errs []error

	if c.Host == "" {
		errs = append(errs, fmt.Errorf("database host cannot be empty"))
	}
	if c.Port <= 0 {
		errs = append(errs, fmt.Errorf("database port must be greater than 0"))
	}
	if c.Port > 65535 {
		errs = append(errs, fmt.Errorf("database port must be less than 65536"))
	}
	if c.User == "" {
		errs = append(errs, fmt.Errorf("database user cannot be empty"))
	}
	if c.Password == "" {
		errs = append(errs, fmt.Errorf("database password cannot be empty"))
	}
	if c.DBName == "" {
		errs = append(errs, fmt.Errorf("database name cannot be empty"))
	}
	if c.SSLMode == "" {
		errs = append(errs, fmt.Errorf("database sslmode cannot be empty"))
	}
	if c.Charset == "" {
		errs = append(errs, fmt.Errorf("database charset cannot be empty"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (c *BackendConfig) Validate() error {
	var errs []error

	if c.ListenAddress == "" {
		errs = append(errs, fmt.Errorf("backend listen address cannot be empty"))
	}
	if len(strings.Fields(c.ListenAddress)) > 1 {
		errs = append(errs, fmt.Errorf("backend listen address cannot contain spaces"))
	}
	if c.ListenPort <= 0 {
		errs = append(errs, fmt.Errorf("backend listen port must be greater than 0"))
	}
	if c.ListenPort > 65535 {
		errs = append(errs, fmt.Errorf("backend listen port must be less than 65536"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (c *LoggerConfig) Validate() error {
	var errs []error

	if c.Level == "" {
		errs = append(errs, fmt.Errorf("logger level cannot be empty"))
	}
	if c.Format == "" {
		errs = append(errs, fmt.Errorf("logger format cannot be empty"))
	}
	if c.LoggerName == "" {
		errs = append(errs, fmt.Errorf("logger name cannot be empty"))
	}

	normalizedLevel := strings.ToLower(string(c.Level))
	normalizedFormat := strings.ToLower(string(c.Format))

	switch logging.LoggerLevel(normalizedLevel) {
	case logging.Debug, logging.Info, logging.Warn, logging.Error, logging.Fatal:
	default:
		errs = append(errs, fmt.Errorf("invalid level: %s", c.Level))
	}

	switch logging.LoggerFormat(normalizedFormat) {
	case logging.JSON, logging.Text:
	default:
		errs = append(errs, fmt.Errorf("invalid format: %s", c.Format))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (c NatsConfig) Validate() error {
	var errs []error

	if c.Url == "" {
		errs = append(errs, fmt.Errorf("nats url cannot be empty"))
	}
	if !strings.HasPrefix(c.Url, "nats://") {
		errs = append(errs, fmt.Errorf("nats url must start with 'nats://'"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (c RedisConfig) Validate() error {
	var errs []error

	if c.Host == "" {
		errs = append(errs, fmt.Errorf("redis host cannot be empty"))
	}
	if c.Port <= 0 {
		errs = append(errs, fmt.Errorf("redis port must be greater than 0"))
	}
	if c.Port > 65535 {
		errs = append(errs, fmt.Errorf("redis port must be less than 65536"))
	}
	if c.DB < 0 {
		errs = append(errs, fmt.Errorf("redis db number cannot be negative"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
