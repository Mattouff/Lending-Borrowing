package service

import (
	"context"
	"time"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
)

// UserService defines the interface for user business logic
type UserService interface {
	// Register creates a new user
	Register(ctx context.Context, address, username string) (*models.User, error)

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uint) (*models.User, error)

	// GetByAddress retrieves a user by Ethereum address
	GetByAddress(ctx context.Context, address string) (*models.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *models.User) error

	// Delete marks a user as deleted (soft delete)
	Delete(ctx context.Context, id uint) error

	// VerifySignature verifies that a signature was created by the user's address
	VerifySignature(ctx context.Context, address, message, signature string) (bool, error)

	// GenerateAuthToken generates an authentication token for a user
	GenerateAuthToken(ctx context.Context, user *models.User) (string, time.Time, error)

	// ListUsers retrieves all users with optional pagination
	ListUsers(ctx context.Context, offset, limit int) ([]*models.User, error)

	// GenerateNonce generates a new nonce for signature verification
	GenerateNonce(address string) string

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)

	// CountWithFilter counts users that match the given filter criteria
	CountWithFilter(ctx context.Context, filter map[string]any) (int64, error)
}
