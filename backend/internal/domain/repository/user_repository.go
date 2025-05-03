package repository

import (
	"context"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create inserts a new user into the database
	Create(ctx context.Context, user *models.User) error

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id uint) (*models.User, error)

	// FindByAddress retrieves a user by Ethereum address
	FindByAddress(ctx context.Context, address string) (*models.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *models.User) error

	// Delete marks a user as deleted (soft delete)
	Delete(ctx context.Context, id uint) error

	// List retrieves all users with optional pagination
	List(ctx context.Context, offset, limit int) ([]*models.User, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)
}
