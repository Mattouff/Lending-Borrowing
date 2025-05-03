package dto

import "time"

// UserRole for DTO types
type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

// UserRegistrationRequest represents the data needed to register a new user
type UserRegistrationRequest struct {
	Address  string `json:"address" validate:"required,eth_addr"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
}

// UserAuthRequest represents data needed for user authentication
type UserAuthRequest struct {
	Address   string `json:"address" validate:"required,eth_addr"`
	Signature string `json:"signature" validate:"required"`
	Message   string `json:"message" validate:"required"`
}

// UserUpdateRequest represents data that can be updated for a user
type UserUpdateRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID        uint      `json:"id"`
	Address   string    `json:"address"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	Verified  bool      `json:"verified"`
	LastLogin time.Time `json:"lastLogin,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserListResponse represents a list of users for API responses
type UserListResponse struct {
	Users     []UserResponse `json:"users"`
	Total     int64          `json:"total"`
	Page      int            `json:"page"`
	PageSize  int            `json:"pageSize"`
	TotalPage int            `json:"totalPage"`
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
