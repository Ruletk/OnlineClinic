package database_test

import (
	clinicConfig "github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/database"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DatabaseTestSuite struct {
	suite.Suite
	Config *clinicConfig.Config
}

func TestDatabase(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (suite *DatabaseTestSuite) SetupTest() {
	suite.Config = &clinicConfig.Config{
		Database: clinicConfig.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			DBName:   "test_db",
			SSLMode:  "disable",
			Charset:  "utf8",
		},
	}
}

func (suite *DatabaseTestSuite) TestConfigConversion() {
	connectionString := database.GetPostgresConnectionString(suite.Config)
	expectedString := "host=localhost port=5432 user=postgres password=password dbname=test_db sslmode=disable"

	suite.Equal(expectedString, connectionString, "The connection string should match the expected format")
}

func (suite *DatabaseTestSuite) TestNewDatabase_Config_NilConfig() {
	_, err := database.NewPostgresDatabase(nil)
	suite.Error(err, "Expected an error when creating a new database with nil config")
	suite.ErrorContains(err, "config cannot be nil", "Expected error to contain 'config cannot be nil' message")
}

func (suite *DatabaseTestSuite) TestNewDatabase_NewDatabase_Error() {
	_, err := database.NewPostgresDatabase(suite.Config)
	suite.Error(err, "Expecting an error when creating a new database with valid config")
	suite.ErrorContains(err, "failed to open database", "Expected error to contain 'failed to open database' message")
	suite.NotContains(err.Error(), "invalid database config", "Expected error to not contain 'invalid database config' message")
}
