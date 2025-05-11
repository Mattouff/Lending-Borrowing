package service

import (
	"context"
	"time"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	// GenerateToken creates a JWT token for a user
	GenerateToken(ctx context.Context, user *models.User) (string, time.Time, error)

	// ValidateToken validates a JWT token and returns the user
	ValidateToken(ctx context.Context, tokenString string) (*models.User, error)

	// InvalidateToken removes a specific token from valid tokens
	InvalidateToken(ctx context.Context, userID uint, tokenID string) error

	// InvalidateAllUserTokens removes all valid tokens for a user
	InvalidateAllUserTokens(ctx context.Context, userID uint) error
}
