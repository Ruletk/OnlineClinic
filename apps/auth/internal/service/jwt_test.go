package service

import (
	"auth/internal/repository"
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

type JwtServiceTestSuite struct {
	suite.Suite
	secret  string
	algo    jwt.SigningMethod
	service JwtService
}

func TestJwtServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JwtServiceTestSuite))
}

func (suite *JwtServiceTestSuite) SetupTest() {
	cfg := &config.Config{
		Logger: config.LoggerConfig{
			LoggerName: "test_jwt",
			TestMode:   true,
		},
	}
	logging.InitLogger(*cfg)
	suite.secret = "testsecret"
	suite.algo = jwt.SigningMethodHS256
	suite.service = NewJwtService(suite.algo, suite.secret)
}

func (suite *JwtServiceTestSuite) TestGenerateToken() {
	payload := jwt.MapClaims{
		"userId": 12345,
		"role":   "user",
	}
	expires := int64(3600) // 1 hour

	token, err := suite.service.GenerateToken(payload, expires)
	suite.NoError(err)
	suite.NotEmpty(token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(suite.secret), nil
	})
	suite.NoError(err)
	suite.True(parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	suite.True(ok)
	suite.Equal(float64(12345), claims["userId"])
	suite.Equal("user", claims["role"])
	suite.WithinDuration(time.Now().Add(time.Second*time.Duration(expires)), time.Unix(int64(claims["exp"].(float64)), 0), time.Minute)
	suite.WithinDuration(time.Now(), time.Unix(int64(claims["iat"].(float64)), 0), time.Minute)
	suite.WithinDuration(time.Now(), time.Unix(int64(claims["nbf"].(float64)), 0), time.Minute)
}

func (suite *JwtServiceTestSuite) TestGenerateVerificationToken() {
	userId := int64(12345)

	token, err := suite.service.GenerateVerificationToken(userId)
	suite.NoError(err)
	suite.NotEmpty(token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(suite.secret), nil
	})
	suite.NoError(err)
	suite.True(parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	suite.True(ok)
	suite.Equal(float64(userId), claims["userId"])
	suite.Equal("verification", claims["type"])
	suite.WithinDuration(time.Now().Add(time.Hour*24*7), time.Unix(int64(claims["exp"].(float64)), 0), time.Minute)
	suite.WithinDuration(time.Now(), time.Unix(int64(claims["iat"].(float64)), 0), time.Minute)
	suite.WithinDuration(time.Now(), time.Unix(int64(claims["nbf"].(float64)), 0), time.Minute)
}

func (suite *JwtServiceTestSuite) TestGeneratePasswordResetToken() {
	userId := int64(12345)

	token, err := suite.service.GeneratePasswordResetToken(userId)
	suite.NoError(err)
	suite.NotEmpty(token)
	fmt.Println(token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(suite.secret), nil
	})
	suite.NoError(err)
	suite.True(parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	suite.True(ok)
	suite.Equal(float64(userId), claims["userId"])
	suite.Equal("password_reset", claims["type"])
	suite.WithinDuration(time.Now().Add(time.Hour*24), time.Unix(int64(claims["exp"].(float64)), 0), time.Minute)
	suite.WithinDuration(time.Now(), time.Unix(int64(claims["iat"].(float64)), 0), time.Minute)
	suite.WithinDuration(time.Now(), time.Unix(int64(claims["nbf"].(float64)), 0), time.Minute)
}

func (suite *JwtServiceTestSuite) TestParseToken() {
	payload := jwt.MapClaims{
		"userId": 12345,
		"role":   "user",
	}
	expires := int64(3600) // 1 hour

	token, err := suite.service.GenerateToken(payload, expires)
	suite.NoError(err)
	suite.NotEmpty(token)

	claims, err := suite.service.ParseToken(token)
	suite.NoError(err)
	suite.NotNil(claims)
	suite.Equal(float64(12345), claims["userId"])
	suite.Equal("user", claims["role"])
}

func (suite *JwtServiceTestSuite) TestParseToken_InvalidToken() {
	invalidToken := "invalid.token.string"

	claims, err := suite.service.ParseToken(invalidToken)
	suite.Error(err)
	suite.Nil(claims)
}

func (suite *JwtServiceTestSuite) TestParseToken_InvalidClaims() {
	// Create a token with invalid claims
	token := jwt.NewWithClaims(suite.algo, jwt.MapClaims{})
	tokenString, err := token.SignedString([]byte(suite.secret))
	suite.NoError(err)

	// No credentials jwt is still jwt but with no claims
	claims, err := suite.service.ParseToken(tokenString)
	suite.NoError(err)
	suite.NotNil(claims)
}

func (suite *JwtServiceTestSuite) TestIsVerificationToken() {
	userId := int64(12345)

	// Generate a valid verification token
	token, err := suite.service.GenerateVerificationToken(userId)
	suite.NoError(err)
	suite.NotEmpty(token)

	// Test with a valid verification token
	isValid, parsedUserId := suite.service.IsVerificationToken(token)
	suite.True(isValid)
	suite.Equal(userId, parsedUserId)

	// Test with an invalid token
	invalidToken := "invalid.token.string"
	isValid, parsedUserId = suite.service.IsVerificationToken(invalidToken)
	suite.False(isValid)
	suite.Equal(int64(0), parsedUserId)

	// Generate a token with a different type
	payload := jwt.MapClaims{
		"userId": userId,
		"type":   "different_type",
	}
	differentToken, err := suite.service.GenerateToken(payload, 3600)
	suite.NoError(err)
	suite.NotEmpty(differentToken)

	// Test with a token of a different type
	isValid, parsedUserId = suite.service.IsVerificationToken(differentToken)
	suite.False(isValid)
	suite.Equal(int64(12345), parsedUserId)
}

func (suite *JwtServiceTestSuite) TestIsPasswordResetToken() {
	userId := int64(12345)

	// Generate a valid password reset token
	token, err := suite.service.GeneratePasswordResetToken(userId)
	suite.NoError(err)
	suite.NotEmpty(token)

	// Test with a valid password reset token
	isValid, parsedUserId := suite.service.IsPasswordResetToken(token)
	suite.True(isValid)
	suite.Equal(userId, parsedUserId)

	// Test with an invalid token
	invalidToken := "invalid.token.string"
	isValid, parsedUserId = suite.service.IsPasswordResetToken(invalidToken)
	suite.False(isValid)
	suite.Equal(int64(0), parsedUserId)

	// Generate a token with a different type
	payload := jwt.MapClaims{
		"userId": userId,
		"type":   "different_type",
	}
	differentToken, err := suite.service.GenerateToken(payload, 3600)
	suite.NoError(err)
	suite.NotEmpty(differentToken)

	// Test with a token of a different type
	isValid, parsedUserId = suite.service.IsPasswordResetToken(differentToken)
	suite.False(isValid)
	suite.Equal(int64(12345), parsedUserId)
}

func (suite *JwtServiceTestSuite) TestGenerateAccessToken_Success() {
	roleAdmin := repository.Role{ID: 1, Name: "user"}
	roleUser := repository.Role{ID: 2, Name: "admin"}
	user := &repository.Auth{ID: 123, Roles: []repository.Role{roleAdmin, roleUser}}

	token, err := suite.service.GenerateAccessToken(user)
	suite.NoError(err)
	suite.NotEmpty(token)

	claims, err := suite.service.ParseToken(token)
	suite.NoError(err)
	suite.NotNil(claims)
	suite.Equal(float64(123), claims["userId"])
	suite.Equal([]interface{}{"user", "admin"}, claims["roles"])
}
func (suite *JwtServiceTestSuite) TestGenerateAccessToken_NoRoles() {
	user := &repository.Auth{ID: 123}

	token, err := suite.service.GenerateAccessToken(user)
	suite.NoError(err)
	suite.NotEmpty(token)

	claims, err := suite.service.ParseToken(token)
	suite.NoError(err)
	suite.NotNil(claims)
	suite.Equal(float64(123), claims["userId"])
	suite.Nil(claims["roles"])
}
