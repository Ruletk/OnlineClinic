package service

import (
	"auth/internal/messages"
	"auth/internal/repository"
	natsmock "auth/mock/nats"
	repositorymock "auth/mock/repository"
	servicemock "auth/mock/service"
	"errors"
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type AuthServiceTestSuite struct {
	suite.Suite
	authRepo       *repositorymock.MockAuthRepository
	sessionService *servicemock.MockSessionService
	jwtService     *servicemock.MockJwtService
	natsPublisher  *natsmock.MockPublisher
	service        AuthService
	storage        *repositorymock.MockStorage
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (suite *AuthServiceTestSuite) SetupTest() {
	logging.InitLogger(config.Config{
		Logger: config.LoggerConfig{
			LoggerName: "test_auth",
			TestMode:   true,
		},
	})

	suite.authRepo = repositorymock.NewMockAuthRepository(suite.T())
	suite.sessionService = servicemock.NewMockSessionService(suite.T())
	suite.jwtService = servicemock.NewMockJwtService(suite.T())
	suite.natsPublisher = natsmock.NewMockPublisher(suite.T())
	suite.storage = repositorymock.NewMockStorage(suite.T())
	suite.service = NewAuthService(
		suite.authRepo, suite.sessionService, suite.jwtService, suite.natsPublisher, suite.storage,
	)
}

func (suite *AuthServiceTestSuite) TestLogin_Success() {
	req := &messages.AuthRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	user := &repository.Auth{
		ID:    1,
		Email: "test@example.com",
	}
	user.PasswordHash = user.GeneratePasswordHash(req.Password)

	suite.authRepo.On("GetByEmail", req.Email).Return(user, nil)
	suite.sessionService.On("CreateSession", user).Return(messages.AuthResponse{Token: "sessiontoken"}, nil)
	suite.authRepo.On("ComparePassword", req.Password).Return(true)

	resp, token, err := suite.service.Login(req)

	suite.NoError(err)
	suite.Equal(200, resp.Code)
	suite.Equal("Successfully authenticated", resp.Message)
	suite.Equal("sessiontoken", token)
}

func (suite *AuthServiceTestSuite) TestLogin_UserNotFound() {
	req := &messages.AuthRequest{
		Email:    "notfound@example.com",
		Password: "password",
	}

	suite.authRepo.On("GetByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)

	resp, token, err := suite.service.Login(req)

	suite.Error(err)
	suite.Nil(resp)
	suite.Empty(token)
}

func (suite *AuthServiceTestSuite) TestLogin_InvalidPassword() {
	req := &messages.AuthRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	user := &repository.Auth{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: "$2a$10$7QJ8K8Q8Q8Q8Q8Q8Q8Q8QO", // hashed password
	}

	suite.authRepo.On("GetByEmail", req.Email).Return(user, nil)
	suite.authRepo.On("ComparePassword", req.Password).Return(false)

	resp, token, err := suite.service.Login(req)

	suite.Error(err)
	suite.Nil(resp)
	suite.Empty(token)
}

func (suite *AuthServiceTestSuite) TestLogin_SessionCreationFailure() {
	req := &messages.AuthRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	user := &repository.Auth{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: "$2a$10$7QJ8K8Q8Q8Q8Q8Q8Q8Q8QO", // hashed password
	}

	suite.authRepo.On("GetByEmail", req.Email).Return(user, nil)
	suite.authRepo.On("ComparePassword", req.Password).Return(true)
	suite.sessionService.On("CreateSession", user.ID).Return(messages.AuthResponse{}, errors.New("session creation failed"))

	resp, token, err := suite.service.Login(req)

	suite.Error(err)
	suite.Nil(resp)
	suite.Empty(token)
}

func (suite *AuthServiceTestSuite) TestRegister_Success() {
	req := &messages.AuthRequest{
		Email:    "newuser@example.com",
		Password: "password",
	}

	suite.authRepo.On("GetByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)
	suite.authRepo.On("Create", mock.AnythingOfType("*repository.Auth")).Return(nil)
	suite.jwtService.On("GenerateVerificationToken", mock.AnythingOfType("int64")).Return("verificationtoken", nil)

	resp, err := suite.service.Register(req)

	suite.NoError(err)
	suite.Equal(201, resp.Code)
	suite.Equal("Successfully registered", resp.Message)
}

func (suite *AuthServiceTestSuite) TestRegister_UserAlreadyExists() {
	req := &messages.AuthRequest{
		Email:    "existinguser@example.com",
		Password: "password",
	}

	suite.authRepo.On("GetByEmail", req.Email).Return(&repository.Auth{}, nil)

	resp, err := suite.service.Register(req)

	suite.Error(err)
	suite.Nil(resp)
	suite.Equal(gorm.ErrDuplicatedKey, err)
}

func (suite *AuthServiceTestSuite) TestRegister_UnexpectedError() {
	req := &messages.AuthRequest{
		Email:    "newuser@example.com",
		Password: "password",
	}

	suite.authRepo.On("GetByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)
	suite.authRepo.On("Create", mock.AnythingOfType("*repository.Auth")).Return(errors.New("unexpected error"))

	resp, err := suite.service.Register(req)

	suite.Error(err)
	suite.Nil(resp)
	suite.Equal("unexpected error", err.Error())
}

func (suite *AuthServiceTestSuite) TestRegister_VerificationEmailFailure() {
	req := &messages.AuthRequest{
		Email:    "newuser@example.com",
		Password: "password",
	}

	suite.authRepo.On("GetByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)
	suite.authRepo.On("Create", mock.AnythingOfType("*repository.Auth")).Return(nil)
	suite.jwtService.On("GenerateVerificationToken", mock.AnythingOfType("int64")).Return("verificationtoken", nil)

	resp, err := suite.service.Register(req)

	suite.Error(err)
	suite.Nil(resp)
	suite.Equal("email send failure", err.Error())
}

func (suite *AuthServiceTestSuite) TestLogout_Success() {
	token := "validtoken"

	suite.sessionService.On("DeleteSession", token).Return(nil)

	err := suite.service.Logout(token)

	suite.NoError(err)
	suite.sessionService.AssertCalled(suite.T(), "DeleteSession", token)
}

func (suite *AuthServiceTestSuite) TestLogout_SessionDeletionFailure() {
	token := "validtoken"
	expectedError := errors.New("session deletion failed")

	suite.sessionService.On("DeleteSession", token).Return(expectedError)

	err := suite.service.Logout(token)

	suite.Error(err)
	suite.Equal(expectedError, err)
	suite.sessionService.AssertCalled(suite.T(), "DeleteSession", token)
}

func (suite *AuthServiceTestSuite) TestSendVerificationEmail_Success() {
	email := "test@example.com"
	user := &repository.Auth{
		ID:     1,
		Email:  email,
		Active: false,
	}

	suite.authRepo.On("GetByEmail", email).Return(user, nil)
	suite.jwtService.On("GenerateVerificationToken", user.ID).Return("verificationtoken", nil)

	err := suite.service.SendVerificationEmail(email)

	suite.NoError(err)
	suite.authRepo.AssertCalled(suite.T(), "GetByEmail", email)
	suite.jwtService.AssertCalled(suite.T(), "GenerateVerificationToken", user.ID)
}

func (suite *AuthServiceTestSuite) TestSendVerificationEmail_UserNotFound() {
	email := "notfound@example.com"

	suite.authRepo.On("GetByEmail", email).Return(nil, gorm.ErrRecordNotFound)

	err := suite.service.SendVerificationEmail(email)

	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)
	suite.authRepo.AssertCalled(suite.T(), "GetByEmail", email)
}

func (suite *AuthServiceTestSuite) TestSendVerificationEmail_UserAlreadyVerified() {
	email := "test@example.com"
	user := &repository.Auth{
		ID:     1,
		Email:  email,
		Active: true,
	}

	suite.authRepo.On("GetByEmail", email).Return(user, nil)

	err := suite.service.SendVerificationEmail(email)

	suite.NoError(err)
	suite.authRepo.AssertCalled(suite.T(), "GetByEmail", email)
}

func (suite *AuthServiceTestSuite) TestSendVerificationEmail_TokenGenerationFailure() {
	email := "test@example.com"
	user := &repository.Auth{
		ID:     1,
		Email:  email,
		Active: false,
	}
	expectedError := errors.New("token generation failed")

	suite.authRepo.On("GetByEmail", email).Return(user, nil)
	suite.jwtService.On("GenerateVerificationToken", user.ID).Return("", expectedError)

	err := suite.service.SendVerificationEmail(email)

	suite.Error(err)
	suite.Equal(expectedError, err)
	suite.authRepo.AssertCalled(suite.T(), "GetByEmail", email)
	suite.jwtService.AssertCalled(suite.T(), "GenerateVerificationToken", user.ID)
}

func (suite *AuthServiceTestSuite) TestSendVerificationEmail_EmailSendingFailure() {
	email := "test@example.com"
	user := &repository.Auth{
		ID:     1,
		Email:  email,
		Active: false,
	}
	expectedError := errors.New("email sending failed")

	suite.authRepo.On("GetByEmail", email).Return(user, nil)
	suite.jwtService.On("GenerateVerificationToken", user.ID).Return("verificationtoken", nil)

	err := suite.service.SendVerificationEmail(email)

	suite.Error(err)
	suite.Equal(expectedError, err)
	suite.authRepo.AssertCalled(suite.T(), "GetByEmail", email)
	suite.jwtService.AssertCalled(suite.T(), "GenerateVerificationToken", user.ID)
}

func (suite *AuthServiceTestSuite) TestChangePassword_Success() {
	req := &messages.PasswordChangeRequest{
		Email: "test@example.com",
	}
	user := &repository.Auth{
		ID:    1,
		Email: req.Email,
	}

	suite.authRepo.On("GetByEmail", req.Email).Return(user, nil)
	suite.jwtService.On("GeneratePasswordResetToken", user.ID).Return("resettoken", nil)

	err := suite.service.RequestChangePassword(req)

	suite.NoError(err)
	suite.authRepo.AssertCalled(suite.T(), "GetByEmail", req.Email)
	suite.jwtService.AssertCalled(suite.T(), "GeneratePasswordResetToken", user.ID)
}

func (suite *AuthServiceTestSuite) TestChangePassword_UserNotFound() {
	req := &messages.PasswordChangeRequest{
		Email: "notfound@example.com",
	}

	suite.authRepo.On("GetByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)

	err := suite.service.RequestChangePassword(req)

	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)
	suite.authRepo.AssertCalled(suite.T(), "GetByEmail", req.Email)
}

func (suite *AuthServiceTestSuite) TestChangePassword_TokenGenerationFailure() {
	req := &messages.PasswordChangeRequest{
		Email: "test@example.com",
	}
	user := &repository.Auth{
		ID:    1,
		Email: req.Email,
	}
	expectedError := errors.New("token generation failed")

	suite.authRepo.On("GetByEmail", req.Email).Return(user, nil)
	suite.jwtService.On("GeneratePasswordResetToken", user.ID).Return("", expectedError)

	err := suite.service.RequestChangePassword(req)

	suite.Error(err)
	suite.Equal(expectedError, err)
	suite.authRepo.AssertCalled(suite.T(), "GetByEmail", req.Email)
	suite.jwtService.AssertCalled(suite.T(), "GeneratePasswordResetToken", user.ID)
}

func (suite *AuthServiceTestSuite) TestChangePassword_EmailSendingFailure() {
	req := &messages.PasswordChangeRequest{
		Email: "test@example.com",
	}
	user := &repository.Auth{
		ID:    1,
		Email: req.Email,
	}
	expectedError := errors.New("email sending failed")

	suite.authRepo.On("GetByEmail", req.Email).Return(user, nil)
	suite.jwtService.On("GeneratePasswordResetToken", user.ID).Return("resettoken", nil)
	err := suite.service.RequestChangePassword(req)

	suite.Error(err)
	suite.Equal(expectedError, err)
	suite.authRepo.AssertCalled(suite.T(), "GetByEmail", req.Email)
	suite.jwtService.AssertCalled(suite.T(), "GeneratePasswordResetToken", user.ID)
}

func (suite *AuthServiceTestSuite) TestResetPassword_Success() {
	req := &messages.PasswordChange{
		NewPassword: "newpassword",
	}
	token := "validtoken"
	userID := int64(1)
	user := &repository.Auth{
		ID: userID,
	}

	suite.jwtService.On("IsPasswordResetToken", token).Return(true, userID)
	suite.authRepo.On("GetByID", userID).Return(user, nil)
	suite.authRepo.On("Update", mock.AnythingOfType("*repository.Auth")).Return(nil)
	suite.jwtService.On("DeleteToken", token).Return(nil)

	err := suite.service.ChangePassword(req, token)

	suite.NoError(err)
	suite.jwtService.AssertCalled(suite.T(), "IsPasswordResetToken", token)
	suite.authRepo.AssertCalled(suite.T(), "GetByID", userID)
	suite.authRepo.AssertCalled(suite.T(), "Update", mock.AnythingOfType("*repository.Auth"))
	suite.jwtService.AssertCalled(suite.T(), "DeleteToken", token)
}

func (suite *AuthServiceTestSuite) TestResetPassword_InvalidToken() {
	req := &messages.PasswordChange{
		NewPassword: "newpassword",
	}
	token := "invalidtoken"

	suite.jwtService.On("IsPasswordResetToken", token).Return(false, int64(0))

	err := suite.service.ChangePassword(req, token)

	suite.Error(err)
	suite.Equal(jwt.ErrTokenInvalidClaims, err)
	suite.jwtService.AssertCalled(suite.T(), "IsPasswordResetToken", token)
}

func (suite *AuthServiceTestSuite) TestResetPassword_UserNotFound() {
	req := &messages.PasswordChange{
		NewPassword: "newpassword",
	}
	token := "validtoken"
	userID := int64(1)

	suite.jwtService.On("IsPasswordResetToken", token).Return(true, userID)
	suite.authRepo.On("GetByID", userID).Return(nil, gorm.ErrRecordNotFound)

	err := suite.service.ChangePassword(req, token)

	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)
	suite.jwtService.AssertCalled(suite.T(), "IsPasswordResetToken", token)
	suite.authRepo.AssertCalled(suite.T(), "GetByID", userID)
}

func (suite *AuthServiceTestSuite) TestResetPassword_UpdateFailure() {
	req := &messages.PasswordChange{
		NewPassword: "newpassword",
	}
	token := "validtoken"
	userID := int64(1)
	user := &repository.Auth{
		ID: userID,
	}
	expectedError := errors.New("update failed")

	suite.jwtService.On("IsPasswordResetToken", token).Return(true, userID)
	suite.authRepo.On("GetByID", userID).Return(user, nil)
	suite.authRepo.On("Update", mock.AnythingOfType("*repository.Auth")).Return(expectedError)

	err := suite.service.ChangePassword(req, token)

	suite.Error(err)
	suite.Equal(expectedError, err)
	suite.jwtService.AssertCalled(suite.T(), "IsPasswordResetToken", token)
	suite.authRepo.AssertCalled(suite.T(), "GetByID", userID)
	suite.authRepo.AssertCalled(suite.T(), "Update", mock.AnythingOfType("*repository.Auth"))
}

func (suite *AuthServiceTestSuite) TestResetPassword_DeleteTokenFailure() {
	req := &messages.PasswordChange{
		NewPassword: "newpassword",
	}
	token := "validtoken"
	userID := int64(1)
	user := &repository.Auth{
		ID: userID,
	}
	expectedError := errors.New("delete token failed")

	suite.jwtService.On("IsPasswordResetToken", token).Return(true, userID)
	suite.authRepo.On("GetByID", userID).Return(user, nil)
	suite.authRepo.On("Update", mock.AnythingOfType("*repository.Auth")).Return(nil)
	suite.jwtService.On("DeleteToken", token).Return(expectedError)

	err := suite.service.ChangePassword(req, token)

	suite.Error(err)
	suite.Equal(expectedError, err)
	suite.jwtService.AssertCalled(suite.T(), "IsPasswordResetToken", token)
	suite.authRepo.AssertCalled(suite.T(), "GetByID", userID)
	suite.authRepo.AssertCalled(suite.T(), "Update", mock.AnythingOfType("*repository.Auth"))
	suite.jwtService.AssertCalled(suite.T(), "DeleteToken", token)
}

func (suite *AuthServiceTestSuite) TestGetUserData_Success() {
	userID := int64(1)
	user := &repository.Auth{
		ID:    userID,
		Email: "test@example.com",
	}

	suite.authRepo.On("GetByID", userID).Return(user, nil)

	resp, err := suite.service.GetUserData(userID)

	suite.NoError(err)
	suite.NotNil(resp)
	suite.Equal(userID, resp.ID)
	suite.Equal(user.Email, resp.Email)
	suite.authRepo.AssertCalled(suite.T(), "GetByID", userID)
}

func (suite *AuthServiceTestSuite) TestGetUserData_UserNotFound() {
	userID := int64(1)

	suite.authRepo.On("GetByID", userID).Return(nil, gorm.ErrRecordNotFound)

	resp, err := suite.service.GetUserData(userID)

	suite.Error(err)
	suite.Nil(resp)
	suite.Equal(gorm.ErrRecordNotFound, err)
	suite.authRepo.AssertCalled(suite.T(), "GetByID", userID)
}

func (suite *AuthServiceTestSuite) TestRefresh_Success() {
	token := "validrefreshtoken"
	user := &repository.Auth{ID: 1}
	newToken := "newaccesstoken"

	suite.sessionService.On("GetSession", token).Return(repository.Session{SessionKey: token, User: user}, nil)
	suite.jwtService.On("GenerateAccessToken", user).Return(newToken, nil)

	refreshedToken, err := suite.service.Refresh(token)

	suite.NoError(err)
	suite.Equal(newToken, refreshedToken)
	suite.sessionService.AssertCalled(suite.T(), "GetSession", token)
	suite.jwtService.AssertCalled(suite.T(), "GenerateAccessToken", user)
}

func (suite *AuthServiceTestSuite) TestRefresh_GetUserIDFailure() {
	token := "invalidrefreshtoken"
	expectedError := errors.New("failed to get user ID")

	suite.sessionService.On("GetSession", token).Return(repository.Session{}, expectedError)

	refreshedToken, err := suite.service.Refresh(token)

	suite.Error(err)
	suite.Empty(refreshedToken)
	suite.Equal(expectedError, err)
	suite.sessionService.AssertCalled(suite.T(), "GetSession", token)
}

func (suite *AuthServiceTestSuite) TestRefresh_GenerateAccessTokenFailure() {
	token := "validrefreshtoken"
	user := &repository.Auth{ID: 1}
	expectedError := errors.New("failed to generate access token")

	suite.sessionService.On("GetSession", token).Return(repository.Session{SessionKey: token, User: user}, nil)
	suite.jwtService.On("GenerateAccessToken", user).Return("", expectedError)

	refreshedToken, err := suite.service.Refresh(token)

	suite.Error(err)
	suite.Empty(refreshedToken)
	suite.Equal(expectedError, err)
	suite.sessionService.AssertCalled(suite.T(), "GetSession", token)
	suite.jwtService.AssertCalled(suite.T(), "GenerateAccessToken", user)
}
