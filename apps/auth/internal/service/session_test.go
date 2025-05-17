package service

import (
	"auth/internal/repository"
	"auth/mocks"
	"errors"
	"github.com/Ruletk/GoMarketplace/pkg/logging"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// Test Suite
type SessionServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.MockSessionRepo
	service  sessionService
}

func TestSessionService(t *testing.T) {
	suite.Run(t, new(SessionServiceTestSuite))
}

// Настройка перед каждым тестом
func (suite *SessionServiceTestSuite) SetupTest() {
	logging.InitTestLogger()
	suite.mockRepo = new(mocks.MockSessionRepo)
	suite.service = sessionService{sessionRepo: suite.mockRepo}
}

// Тесты для CreateSession
func (suite *SessionServiceTestSuite) TestCreateSession_Success() {
	expectedUser := &repository.Auth{ID: 123}
	expectedToken := "test_token_123"

	suite.mockRepo.On("Create", mock.MatchedBy(func(s *repository.Session) bool {
		return s.UserID == expectedUser.ID &&
			len(s.SessionKey) == 64 &&
			s.ExpiresAt.After(time.Now())
	})).Return(nil).Run(func(args mock.Arguments) {
		s := args.Get(0).(*repository.Session)
		s.SessionKey = expectedToken
	})

	// Act
	response, err := suite.service.CreateSession(expectedUser)

	// Assert
	suite.NoError(err)
	suite.Equal(expectedToken, response.Token)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SessionServiceTestSuite) TestCreateSession_RepoError() {
	expectedError := errors.New("repo error")
	suite.mockRepo.On("Create", mock.Anything).Return(expectedError)

	_, err := suite.service.CreateSession(&repository.Auth{})

	suite.ErrorIs(err, expectedError)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Тесты для GetSession
func (suite *SessionServiceTestSuite) TestGetSession_Success() {
	token := "valid_token"
	expectedSession := &repository.Session{SessionKey: token}

	suite.mockRepo.On("Get", token).Return(expectedSession, nil)
	suite.mockRepo.On("UpdateLastUsed", token).Return(nil)

	result, err := suite.service.GetSession(token)

	suite.NoError(err)
	suite.Equal(*expectedSession, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SessionServiceTestSuite) TestGetSession_NotFound() {
	token := "invalid_token"
	expectedError := errors.New("not found")

	suite.mockRepo.On("Get", token).Return(nil, expectedError)

	_, err := suite.service.GetSession(token)

	suite.ErrorIs(err, expectedError)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Тесты для GetUserID
func (suite *SessionServiceTestSuite) TestGetUserID_Success() {
	token := "valid_token"
	expectedUser := &repository.Auth{ID: 1}
	mockSession := &repository.Session{User: expectedUser}

	suite.mockRepo.On("Get", token).Return(mockSession, nil)
	suite.mockRepo.On("UpdateLastUsed", token).Return(nil)

	userID, err := suite.service.GetUserID(token)

	suite.NoError(err)
	suite.Equal(expectedUser.ID, userID)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Тесты для DeleteSession
func (suite *SessionServiceTestSuite) TestDeleteSession_Success() {
	token := "valid_token"
	mockSession := &repository.Session{
		SessionKey: token,
		ExpiresAt:  time.Now().Add(1 * time.Hour),
	}

	suite.mockRepo.On("Get", token).Return(mockSession, nil)
	suite.mockRepo.On("Delete", token).Return(nil)

	err := suite.service.DeleteSession(token)

	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SessionServiceTestSuite) TestDeleteSession_Expired() {
	token := "expired_token"
	mockSession := &repository.Session{
		SessionKey: token,
		ExpiresAt:  time.Now().Add(-1 * time.Hour),
	}

	suite.mockRepo.On("Get", token).Return(mockSession, nil)

	err := suite.service.DeleteSession(token)

	suite.ErrorIs(err, gorm.ErrRecordNotFound)
	suite.mockRepo.AssertNotCalled(suite.T(), "Delete")
}

// Тесты для HardDeleteSessions
func (suite *SessionServiceTestSuite) TestHardDeleteSessions_Success() {
	suite.mockRepo.On("HardDeleteAllExpired").Return(nil)

	err := suite.service.HardDeleteSessions()

	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Тесты для DeleteInactiveSessions
func (suite *SessionServiceTestSuite) TestDeleteInactiveSessions_Success() {
	suite.mockRepo.On("HardDeleteAllInactive").Return(nil)

	err := suite.service.DeleteInactiveSessions()

	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
}
