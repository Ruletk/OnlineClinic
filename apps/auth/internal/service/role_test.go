package service

import (
	repositorymock "auth/mock/repository"
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RoleServiceTestSuite struct {
	suite.Suite
	mockRepo *repositorymock.MockRoleRepository
	service  RoleService
}

func TestRoleService(t *testing.T) {
	suite.Run(t, new(RoleServiceTestSuite))
}

func (suite *RoleServiceTestSuite) SetupTest() {
	logging.InitLogger(config.Config{
		Logger: config.LoggerConfig{
			LoggerName: "test_role",
			TestMode:   true,
		},
	})
	suite.mockRepo = repositorymock.NewMockRoleRepository(suite.T())
	suite.service = &roleService{roleRepository: suite.mockRepo}
}
