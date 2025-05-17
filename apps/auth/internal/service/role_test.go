package service

import (
	"auth/mocks"
	"github.com/Ruletk/GoMarketplace/pkg/logging"
	"github.com/stretchr/testify/suite"
	"testing"
)

// Test Suite
type RoleServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.MockRoleRepository
	service  RoleService
}

func TestRoleService(t *testing.T) {
	suite.Run(t, new(RoleServiceTestSuite))
}

// Настройка перед каждым тестом
func (suite *RoleServiceTestSuite) SetupTest() {
	logging.InitTestLogger()
	suite.mockRepo = new(mocks.MockRoleRepository)
	suite.service = &roleService{roleRepository: suite.mockRepo}
}
