package repository

import (
	"github.com/Mattouff/Lending-Borrowing/internal/models"

	"gorm.io/gorm"
)

// UserRepository handles database operations for User model
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// Create adds a new user to the database
func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

// FindByWalletAddress finds a user by wallet address
func (r *UserRepository) FindByWalletAddress(walletAddress string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("wallet_address = ?", walletAddress).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user in the database
func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

// UpdateNonce updates a user's nonce
func (r *UserRepository) UpdateNonce(walletAddress string, nonce string) error {
	return r.DB.Model(&models.User{}).Where("wallet_address = ?", walletAddress).Update("nonce", nonce).Error
}

// Delete removes a user from the database
func (r *UserRepository) Delete(walletAddress string) error {
	return r.DB.Where("wallet_address = ?", walletAddress).Delete(&models.User{}).Error
}
