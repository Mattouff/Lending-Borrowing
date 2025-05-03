package postgres

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/repository"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new PostgreSQL implementation of UserRepository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create inserts a new user into the database
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID retrieves a user by ID
func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByAddress retrieves a user by Ethereum address
func (r *userRepository) FindByAddress(ctx context.Context, address string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("address = ?", address).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete marks a user as deleted (soft delete)
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// List retrieves all users with optional pagination
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	query := r.db.WithContext(ctx).Model(&models.User{})

	if offset >= 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error
	return count, err
}

// Add this method to your userRepository implementation:

// CountWithFilter counts users that match the given filter criteria
func (r *userRepository) CountWithFilter(ctx context.Context, filter map[string]any) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.User{})

	// Apply filters
	for key, value := range filter {
		// Special handling for time-based filters
		if key == "last_login_after" {
			if timeValue, ok := value.(time.Time); ok {
				query = query.Where("last_login >= ?", timeValue)
			}
		} else if key == "last_login_before" {
			if timeValue, ok := value.(time.Time); ok {
				query = query.Where("last_login <= ?", timeValue)
			}
		} else {
			query = query.Where(key+" = ?", value)
		}
	}

	err := query.Count(&count).Error
	return count, err
}
