package service

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"fmt"

	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/repository"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

type userService struct {
	userRepo    repository.UserRepository
	cfg         *config.Config
	authService service.AuthService
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, cfg *config.Config, authService service.AuthService) service.UserService {
	return &userService{
		userRepo:    userRepo,
		cfg:         cfg,
		authService: authService,
	}
}

// Register creates a new user
func (s *userService) Register(ctx context.Context, address, username string) (*models.User, error) {
	// Validate the Ethereum address
	if !common.IsHexAddress(address) {
		return nil, errors.New("invalid ethereum address")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.FindByAddress(ctx, address)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this address already exists")
	}

	// Create a random nonce for the user
	nonce := crypto.Keccak256Hash([]byte(time.Now().String())).Hex()

	user := &models.User{
		Address:   address,
		Username:  username,
		Role:      models.RoleUser,
		Verified:  false,
		Nonce:     nonce,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// GetByAddress retrieves a user by Ethereum address
func (s *userService) GetByAddress(ctx context.Context, address string) (*models.User, error) {
	return s.userRepo.FindByAddress(ctx, address)
}

// Update updates an existing user
func (s *userService) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
}

// VerifySignature verifies that a signature was created by the user's address
func (s *userService) VerifySignature(ctx context.Context, address, message, signature string) (bool, error) {
	// Convert the signature to the appropriate format
	sig := common.FromHex(signature)
	if len(sig) != 65 {
		return false, errors.New("invalid signature length")
	}

	// Ethereum message prefix
	prefix := "\x19Ethereum Signed Message:\n" + fmt.Sprint(len(message)) + message
	prefixedHash := crypto.Keccak256Hash([]byte(prefix))

	// Check if V needs adjustment for Ethereum's quirks
	if sig[64] > 26 {
		sig[64] -= 27
	}

	// Recover the public key from the signature
	pubKey, err := crypto.Ecrecover(prefixedHash.Bytes(), sig)
	if err != nil {
		return false, err
	}

	// Convert the public key to an Ethereum address
	recoveredAddr := common.BytesToAddress(crypto.Keccak256(pubKey[1:])[12:])

	// Compare with the provided address
	return recoveredAddr.Hex() == common.HexToAddress(address).Hex(), nil
}

// ListUsers retrieves all users with optional pagination
func (s *userService) ListUsers(ctx context.Context, offset, limit int) ([]*models.User, error) {
	return s.userRepo.List(ctx, offset, limit)
}

// GenerateNonce generates a new nonce for signature verification
func (s *userService) GenerateNonce(address string) string {
	// Create a random nonce by hashing the address with a timestamp
	timestamp := time.Now().String()
	data := address + timestamp
	return crypto.Keccak256Hash([]byte(data)).Hex()
}

// Count returns the total number of users
func (s *userService) Count(ctx context.Context) (int64, error) {
	return s.userRepo.Count(ctx)
}

// CountWithFilter counts users that match the given filter criteria
func (s *userService) CountWithFilter(ctx context.Context, filter map[string]any) (int64, error) {
	return s.userRepo.CountWithFilter(ctx, filter)
}

// Delete marks a user as deleted (soft delete) and invalidates all tokens
func (s *userService) Delete(ctx context.Context, id uint) error {
	// First invalidate all user tokens
	if err := s.authService.InvalidateAllUserTokens(ctx, id); err != nil {
		// Log the error but continue with deletion
		fmt.Printf("Failed to invalidate tokens for user %d: %v\n", id, err)
	}

	// Then perform the soft delete
	return s.userRepo.Delete(ctx, id)
}

// GenerateAuthToken generates an authentication token for a user
func (s *userService) GenerateAuthToken(ctx context.Context, user *models.User) (string, time.Time, error) {
	return s.authService.GenerateToken(ctx, user)
}
