package api

import (
	"auth/internal/messages"
	"auth/internal/service"
	"errors"
	"github.com/Ruletk/GoMarketplace/pkg/logging"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type AuthAPI struct {
	authService    service.AuthService
	sessionService service.SessionService
	roleService    service.RoleService
}

func NewAuthAPI(authService service.AuthService, sessionService service.SessionService, roleService service.RoleService) *AuthAPI {
	return &AuthAPI{authService: authService, sessionService: sessionService, roleService: roleService}
}

func (api *AuthAPI) RegisterRoutes(router *gin.RouterGroup) {
	logging.Logger.Info("Registering public routes")
	router.GET("/verify/:token", api.Verify)
	router.GET("/refresh", api.Refresh)

	logging.Logger.Info("Registering public only routes")
	router.POST("/login", api.Login)
	router.POST("/register", api.Register)
	router.POST("/change-password", api.ChangePassword)
	router.POST("/change-password/:token", api.ChangePasswordWithToken)

	logging.Logger.Info("Registering private routes")
	router.GET("/logout", api.Logout) // Required authentication
	router.DELETE("/admin/sessions/hard-delete", api.HardDeleteSessions)
	router.DELETE("/admin/sessions/delete-inactive", api.DeleteInactiveSessions)
	router.POST("/role", api.CreateRole)
	router.DELETE("/role", api.DeleteRole)
	router.POST("/role/assign", api.AssignUserToRole)
	router.POST("/role/unassign", api.UnassignUserFromRole)
}

func (api *AuthAPI) Login(c *gin.Context) {
	logging.Logger.Info("Logging in user")

	// Parse the request
	var req messages.AuthRequest
	err := c.ShouldBindJSON(&req)

	// Check if the request is valid
	if err != nil {
		logging.Logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, messages.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Invalid request",
		})
		return
	}

	// Authenticate the user
	resp, token, err := api.authService.Login(&req)
	if errors.Is(err, service.ErrInvalidCredentials) {
		logging.Logger.WithError(err).Error("Wrong email or password")
		c.JSON(http.StatusUnauthorized, messages.ApiResponse{
			Code:    http.StatusUnauthorized,
			Type:    "error",
			Message: "Wrong email or password",
		})
		return
	} else if err != nil {
		logging.Logger.WithError(err).Error("Internal server error")
		c.JSON(http.StatusInternalServerError, messages.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "error",
			Message: "Internal server error. Details: " + err.Error(),
		})
		return
	}

	if resp != nil {
		logging.Logger.Info("User logged in successfully")
		c.SetCookie("token", token, 31536000, "/", "", false, true)
		c.JSON(http.StatusOK, resp)
		return
	}

	logging.Logger.Error("Wrong email or password")
	// Return an error if the user is not authenticated
	c.JSON(http.StatusUnauthorized, messages.ApiResponse{
		Code:    401,
		Type:    "error",
		Message: "Wrong email or password",
	})

}

func (api *AuthAPI) Register(c *gin.Context) {
	logging.Logger.Info("Registering user")

	var req messages.AuthRequest
	err := c.ShouldBindJSON(&req)

	// Check if the request is valid
	if err != nil {
		logging.Logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, messages.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Invalid request",
		})
		return
	}

	logging.Logger.Info("Registering user with email: " + req.Email)

	// Register the user
	resp, err := api.authService.Register(&req)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		logging.Logger.WithError(err).Error("User with this email already registered")
		c.JSON(http.StatusConflict, messages.ApiResponse{
			Code:    http.StatusConflict,
			Type:    "error",
			Message: "User with this email already registered",
		})
		return
	} else if err != nil {
		logging.Logger.WithError(err).Error("Internal server error")
		c.JSON(http.StatusInternalServerError, messages.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "error",
			Message: "Internal server error. Details: " + err.Error(),
		})
		return
	}

	if resp != nil {
		logging.Logger.Info("User registered successfully")
		c.JSON(resp.Code, resp)
		return
	}

	logging.Logger.Error("User with this email already registered")
	// WTF, this is a duplicate of the error above
	// TODO: Refactor this
	// Return an error if the user is already registered
	c.JSON(http.StatusConflict, messages.ApiResponse{
		Code:    http.StatusConflict,
		Type:    "error",
		Message: "User with this email already registered",
	})
}

func (api *AuthAPI) Logout(c *gin.Context) {
	logging.Logger.Info("Logging out user")

	token, err := c.Cookie("token")
	if err != nil {
		logging.Logger.Error("No user provided")
		c.JSON(http.StatusBadRequest, messages.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "No token provided",
		})
		return
	}

	// Logout the user
	_ = api.authService.Logout(token)

	logging.Logger.Info("User logged out successfully, token: " + token)

	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, messages.ApiResponse{
		Code:    http.StatusOK,
		Type:    "success",
		Message: "Successfully logged out",
	})
}

func (api *AuthAPI) ChangePassword(c *gin.Context) {
	logging.Logger.Info("Changing password")

	var req messages.PasswordChangeRequest
	err := c.ShouldBindJSON(&req)

	// Check if the request is valid
	if err != nil {
		logging.Logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, messages.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Invalid request or email",
		})
		return
	}

	// Send an email with a token to the user
	logging.Logger.Info("Changing password for email: " + req.Email)
	err = api.authService.ChangePassword(&req)
	if err == nil {
		logging.Logger.Info("Password change sent for email: " + req.Email)
		domain := string([]rune(req.Email)[strings.Index(req.Email, "@")+1:])
		c.JSON(http.StatusOK, messages.ApiResponse{
			Code:    http.StatusOK,
			Type:    "success",
			Message: "Password change request sent successfully. Check your email for further instructions. Domain: " + domain,
		})
		return
	}

	// Return an error if the email is invalid
	logging.Logger.WithError(err).Error("Invalid request or email")
	c.JSON(http.StatusBadRequest, messages.ApiResponse{
		Code:    http.StatusBadRequest,
		Type:    "error",
		Message: "Invalid request or email",
	})
}

func (api *AuthAPI) ChangePasswordWithToken(c *gin.Context) {
	logging.Logger.Info("Changing password with token")

	token := c.Param("token")
	var req messages.PasswordChange
	err := c.ShouldBindJSON(&req)

	// Check if the request is valid
	if err != nil || token == "" {
		logging.Logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, messages.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Invalid request",
		})
		return
	}

	// Change the password
	err = api.authService.ResetPassword(&req, token)
	if err == nil {
		logging.Logger.Info("Password changed successfully")
		c.JSON(http.StatusOK, messages.ApiResponse{
			Code:    http.StatusOK,
			Type:    "success",
			Message: "Password changed successfully",
		})
		return
	}

	// Return an error if the token is invalid
	logging.Logger.WithError(err).Error("Invalid token")
	c.JSON(http.StatusUnauthorized, messages.ApiResponse{
		Code:    http.StatusUnauthorized,
		Type:    "error",
		Message: "Invalid token",
	})
}

func (api *AuthAPI) Verify(c *gin.Context) {
	//TODO: Refactor this to send X-Access-Token header
	token := c.Param("token")
	// Check if the token is valid
	if token == "" {
		c.JSON(http.StatusBadRequest, messages.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Invalid request",
		})
		return
	}

	// Verify the token
	err := api.authService.VerifyUser(token)
	if err == nil {
		c.JSON(http.StatusOK, messages.ApiResponse{
			Code:    http.StatusOK,
			Type:    "success",
			Message: "User verified successfully",
		})
		return
	}

	// Return an error if the token is invalid
	c.JSON(http.StatusUnauthorized, messages.ApiResponse{
		Code:    http.StatusUnauthorized,
		Type:    "error",
		Message: "Invalid token",
	})
}

func (api *AuthAPI) HardDeleteSessions(c *gin.Context) {
	// TODO: Add admin check
	logging.Logger.Info("Starting delete all expired sessions...")
	err := api.sessionService.HardDeleteSessions()
	if err == nil {
		logging.Logger.Info("Sessions deleted successfully")
		c.JSON(http.StatusOK, messages.ApiResponse{
			Code:    http.StatusOK,
			Type:    "success",
			Message: "Sessions deleted successfully",
		})
		return
	}

	logging.Logger.WithError(err).Error("Internal server error")
	c.JSON(http.StatusInternalServerError, messages.ApiResponse{
		Code:    http.StatusInternalServerError,
		Type:    "error",
		Message: "Internal server error. Details: " + err.Error(),
	})
}

func (api *AuthAPI) DeleteInactiveSessions(c *gin.Context) {
	// TODO: Add admin check
	logging.Logger.Info("Starting delete all inactive sessions...")

	err := api.sessionService.DeleteInactiveSessions()

	if err == nil {
		logging.Logger.Info("Sessions deleted successfully")
		c.JSON(http.StatusOK, messages.ApiResponse{
			Code:    http.StatusOK,
			Type:    "success",
			Message: "Sessions deleted successfully",
		})
		return
	}

	logging.Logger.WithError(err).Error("Internal server error")

	c.JSON(http.StatusInternalServerError, messages.ApiResponse{
		Code:    http.StatusInternalServerError,
		Type:    "error",
		Message: "Internal server error. Details: " + err.Error(),
	})
}

func (api *AuthAPI) Refresh(c *gin.Context) {
	logging.Logger.Info("Generating Access Token")

	refreshToken, err := c.Cookie("token")
	if err != nil {
		logging.Logger.Info("No refresh token provided. Skipping")
		c.JSON(http.StatusOK, messages.ApiResponse{
			Code:    http.StatusOK,
			Type:    "info",
			Message: "No refresh token provided. Skipping",
		})
		return
	}

	token, err := api.authService.Refresh(refreshToken)
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to refresh token")
		c.JSON(http.StatusUnauthorized, messages.ApiResponse{
			Code:    http.StatusUnauthorized,
			Type:    "error",
			Message: "Invalid token",
		})
		return
	}

	c.Header("X-Access-Token", token)
	c.JSON(http.StatusOK, messages.ApiResponse{
		Code:    http.StatusOK,
		Type:    "success",
		Message: "Token refreshed successfully",
	})
}

func (api *AuthAPI) CreateRole(c *gin.Context) {
	var req messages.RoleRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		logging.Logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, messages.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Invalid request",
		})
		return
	}

	err = api.roleService.Create(req.Name)
	if errors.Is(err, gorm.ErrDuplicatedKey) || (err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")) {
		logging.Logger.WithError(err).Error("Role already exists")
		c.JSON(http.StatusConflict, messages.ApiResponse{
			Code:    http.StatusConflict,
			Type:    "error",
			Message: "Role already exists",
		})
		return
	} else if err != nil {
		logging.Logger.WithError(err).Error("Failed to create role")
		c.JSON(http.StatusInternalServerError, messages.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "error",
			Message: "Internal server error. Details: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, messages.ApiResponse{
		Code:    http.StatusOK,
		Type:    "success",
		Message: "Role created successfully",
	})
}

func (api *AuthAPI) DeleteRole(c *gin.Context) {
	var req messages.RoleRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		logging.Logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, messages.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Invalid request",
		})
		return
	}

	err = api.roleService.Delete(req.Name)
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to delete role")
		c.JSON(http.StatusInternalServerError, messages.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "error",
			Message: "Internal server error. Details: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, messages.ApiResponse{
		Code:    http.StatusOK,
		Type:    "success",
		Message: "Role deleted successfully",
	})

}

func (api *AuthAPI) AssignUserToRole(c *gin.Context) {
	var req messages.RoleAssignRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logging.Logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, messages.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Invalid request",
		})
		return
	}

	role, err := api.roleService.Get(req.Name)
	if err != nil {
		logging.Logger.Info("Role not found: ", req.Name)
		c.JSON(http.StatusNotFound, messages.ApiResponse{
			Code:    http.StatusNotFound,
			Type:    "error",
			Message: "Role not found",
		})
		return
	}

	err = api.authService.AddRoleToUser(req.UserID, role)
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to assign role to user")
		c.JSON(http.StatusInternalServerError, messages.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "error",
			Message: "Internal server error. Details: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, messages.ApiResponse{
		Code:    http.StatusOK,
		Type:    "success",
		Message: "Role assigned successfully",
	})
}

func (api *AuthAPI) UnassignUserFromRole(c *gin.Context) {
	var req messages.RoleAssignRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logging.Logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, messages.ApiResponse{
			Code:    http.StatusBadRequest,
			Type:    "error",
			Message: "Invalid request",
		})
		return
	}

	role, err := api.roleService.Get(req.Name)
	if err != nil {
		logging.Logger.Info("Role not found: ", req.Name)
		c.JSON(http.StatusNotFound, messages.ApiResponse{
			Code:    http.StatusNotFound,
			Type:    "error",
			Message: "Role not found",
		})
		return
	}

	err = api.authService.RemoveRoleFromUser(req.UserID, role)
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to unassign role from user")
		c.JSON(http.StatusInternalServerError, messages.ApiResponse{
			Code:    http.StatusInternalServerError,
			Type:    "error",
			Message: "Internal server error. Details: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, messages.ApiResponse{
		Code:    http.StatusOK,
		Type:    "success",
		Message: "Role unassigned successfully",
	})
}
