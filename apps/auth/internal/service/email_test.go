package service

import (
	"auth/mocks"
	"errors"
	"github.com/Ruletk/GoMarketplace/pkg/logging"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gomail.v2"
	"testing"
)

type EmailServiceTestSuite struct {
	suite.Suite
	mockDialer *mocks.MockDialer
	service    EmailService
}

func TestEmailService(t *testing.T) {
	suite.Run(t, new(EmailServiceTestSuite))
}

func (suite *EmailServiceTestSuite) SetupTest() {
	logging.InitTestLogger()
	suite.mockDialer = new(mocks.MockDialer)
	suite.service = NewEmailService(suite.mockDialer)
}

func (suite *EmailServiceTestSuite) TestEmailService_SendVerificationEmail() {
	email := "test@example.com"
	token := "test-token"

	// Настраиваем ожидаемый вызов
	suite.mockDialer.On("DialAndSend", mock.MatchedBy(func(messages []*gomail.Message) bool {
		msg := messages[0]
		to := msg.GetHeader("To")
		subject := msg.GetHeader("Subject")
		return to[0] == email &&
			subject[0] == "Verify your email"
	})).Return(nil)

	err := suite.service.SendVerificationEmail(email, token)

	suite.NoError(err)
	suite.mockDialer.AssertExpectations(suite.T())
}

func (suite *EmailServiceTestSuite) TestEmailService_SendPasswordResetEmail() {
	email := "user@example.com"
	token := "reset-token"

	suite.mockDialer.On("DialAndSend", mock.MatchedBy(func(messages []*gomail.Message) bool {
		msg := messages[0]
		to := msg.GetHeader("To")
		subject := msg.GetHeader("Subject")
		return to[0] == email &&
			subject[0] == "Password reset"
	})).Return(nil)

	err := suite.service.SendPasswordResetEmail(email, token)

	suite.NoError(err)
	suite.mockDialer.AssertExpectations(suite.T())
}

func (suite *EmailServiceTestSuite) TestEmailService_SendEmailErrorHandling() {
	expectedError := errors.New("smtp error")

	suite.mockDialer.On("DialAndSend", mock.Anything).Return(expectedError)

	err := suite.service.SendVerificationEmail("test@error.com", "token")

	suite.Error(err)
	suite.Equal(expectedError, err)
}

func (suite *EmailServiceTestSuite) TestNewEmailService() {
	dialer := gomail.NewDialer("smtp.example.com", 587, "user", "pass")
	service := NewEmailService(dialer)

	suite.IsType(&emailService{}, service)
	suite.Equal(dialer, service.(*emailService).dialer)
}
