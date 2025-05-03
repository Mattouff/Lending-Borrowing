package repository

import (
	"context"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
)

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	// Create inserts a new transaction into the database
	Create(ctx context.Context, transaction *models.Transaction) error

	// FindByID retrieves a transaction by ID
	FindByID(ctx context.Context, id uint) (*models.Transaction, error)

	// FindByHash retrieves a transaction by blockchain transaction hash
	FindByHash(ctx context.Context, hash string) (*models.Transaction, error)

	// FindByUserID retrieves all transactions for a specific user
	FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]*models.Transaction, error)

	// Update updates an existing transaction
	Update(ctx context.Context, transaction *models.Transaction) error

	// List retrieves all transactions with optional filtering and pagination
	List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*models.Transaction, error)

	// Count returns the total number of transactions matching the filter
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}
