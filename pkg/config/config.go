package config

import (
	_var "github.com/Ruletk/OnlineClinic/pkg/config/logging"
	"io"
)

type Config struct {
	Database DatabaseConfig
	Backend  BackendConfig
	Logger   LoggerConfig
}

type DatabaseConfig struct {
	Host     string // Database host, ip address or domain name
	Port     int    // Database port, usually 5432 for postgres
	User     string // Database user to authenticate
	Password string // Database password to authenticate
	DBName   string // Database name to use this database
	SSLMode  string // sslmode can be "disable", "require", "verify-ca", or "verify-full"
	Charset  string
}

type BackendConfig struct {
	ListenAddress string // Address to listen on, ip address or domain name. Use `0.0.0.0` for all interfaces
	ListenPort    int    // Port to listen on.
}

type LoggerConfig struct {
	Level        _var.LoggerLevel
	Format       _var.LoggerFormat
	EnableCaller bool
	Output       io.Writer
	TestMode     bool
	LoggerName   string
}
