package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/repository"
)

type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new PostgreSQL implementation of TransactionRepository
func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

// Create inserts a new transaction into the database
func (r *transactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

// FindByID retrieves a transaction by ID
func (r *transactionRepository) FindByID(ctx context.Context, id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	result := r.db.WithContext(ctx).First(&transaction, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction, nil
}

// FindByHash retrieves a transaction by blockchain transaction hash
func (r *transactionRepository) FindByHash(ctx context.Context, hash string) (*models.Transaction, error) {
	var transaction models.Transaction
	result := r.db.WithContext(ctx).Where("hash = ?", hash).First(&transaction)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction, nil
}

// FindByUserID retrieves all transactions for a specific user
func (r *transactionRepository) FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if offset >= 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

// Update updates an existing transaction
func (r *transactionRepository) Update(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

// List retrieves all transactions with optional filtering and pagination
func (r *transactionRepository) List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	query := r.db.WithContext(ctx).Model(&models.Transaction{})

	// Apply filters
	for key, value := range filter {
		// Special handling for array values in filters (e.g., IN queries)
		if key == "type" {
			if types, ok := value.([]models.TransactionType); ok && len(types) > 0 {
				query = query.Where("type IN ?", types)
			} else {
				query = query.Where("type = ?", value)
			}
		} else {
			query = query.Where(key+" = ?", value)
		}
	}

	if offset >= 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

// Count returns the total number of transactions matching the filter
func (r *transactionRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Transaction{})

	// Apply filters
	for key, value := range filter {
		// Special handling for array values in filters (e.g., IN queries)
		if key == "type" {
			if types, ok := value.([]models.TransactionType); ok && len(types) > 0 {
				query = query.Where("type IN ?", types)
			} else {
				query = query.Where("type = ?", value)
			}
		} else {
			query = query.Where(key+" = ?", value)
		}
	}

	err := query.Count(&count).Error
	return count, err
}
