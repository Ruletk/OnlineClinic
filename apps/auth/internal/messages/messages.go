package messages

// AuthRequest represents a login request
type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RoleRequest represents a role request for creation or deletion
type RoleRequest struct {
	Name string `json:"name" binding:"required"`
}

// RoleAssignRequest represents a request to assign or unassign a role to a user
type RoleAssignRequest struct {
	Name   string `json:"name" binding:"required"`
	UserID int64  `json:"user_id" binding:"required"`
}

// TokenRequest represents a token request
type TokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// AuthResponse represents a successful authentication response
type AuthResponse struct {
	Token string `json:"token"`
}

// ApiResponse represents a generic API response
type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

// PasswordChangeRequest represents a request to change a password
type PasswordChangeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordChange represents the new password
type PasswordChange struct {
	NewPassword string `json:"newPassword" binding:"required"`
}

// AuthDataResponse represents the response to a validation request
type AuthDataResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
