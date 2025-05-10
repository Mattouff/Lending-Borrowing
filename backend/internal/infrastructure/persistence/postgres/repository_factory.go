package postgres

import (
	"sync"

	"gorm.io/gorm"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/repository"
)

// RepositoryFactory provides a centralized way to create and access repositories
type RepositoryFactory struct {
	db              *gorm.DB
	userRepository  repository.UserRepository
	transactionRepo repository.TransactionRepository
	positionRepo    repository.PositionRepository

	userOnce        sync.Once
	transactionOnce sync.Once
	positionOnce    sync.Once
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(db *gorm.DB) *RepositoryFactory {
	return &RepositoryFactory{
		db: db,
	}
}

// GetUserRepository returns a singleton instance of UserRepository
func (f *RepositoryFactory) GetUserRepository() repository.UserRepository {
	f.userOnce.Do(func() {
		f.userRepository = NewUserRepository(f.db)
	})
	return f.userRepository
}

// GetTransactionRepository returns a singleton instance of TransactionRepository
func (f *RepositoryFactory) GetTransactionRepository() repository.TransactionRepository {
	f.transactionOnce.Do(func() {
		f.transactionRepo = NewTransactionRepository(f.db)
	})
	return f.transactionRepo
}

// GetPositionRepository returns a singleton instance of PositionRepository
func (f *RepositoryFactory) GetPositionRepository() repository.PositionRepository {
	f.positionOnce.Do(func() {
		f.positionRepo = NewPositionRepository(f.db)
	})
	return f.positionRepo
}
