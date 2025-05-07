package database

import (
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresDatabase function creating a new database instance.
// For our project, it's postgres.
// Connection string and parameters are taken from the config object.
func NewPostgresDatabase(config *config.Config) (*gorm.DB, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if err := config.Database.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database config: %w", err)
	}

	url := GetPostgresConnectionString(config)
	db := postgres.Open(url)

	database, err := gorm.Open(db, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return database, nil
}

// GetPostgresConnectionString function creating a connection string for postgres.
// Converts the config object to a string.
func GetPostgresConnectionString(config *config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
	)
}
