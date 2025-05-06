package config

import (
	"errors"
	"fmt"
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
