package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/Mattouff/Lending-Borrowing/internal/api/middleware"
	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/repository"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
	valkey "github.com/Mattouff/Lending-Borrowing/pkg/cache"
)

// authService implements the AuthService interface
type authService struct {
	cfg          *config.Config
	userRepo     repository.UserRepository
	valkeyClient *valkey.Client
}

// NewAuthService creates a new authentication service
func NewAuthService(cfg *config.Config, userRepo repository.UserRepository, valkeyClient *valkey.Client) service.AuthService {
	return &authService{
		cfg:          cfg,
		userRepo:     userRepo,
		valkeyClient: valkeyClient,
	}
}

// GenerateToken creates a JWT token for a user
func (s *authService) GenerateToken(ctx context.Context, user *models.User) (string, time.Time, error) {
	// Generate unique token ID
	tokenID := uuid.New().String()

	// Set expiration time
	expireTime := time.Duration(s.cfg.JWT.ExpireTime) * time.Minute
	expirationTime := time.Now().Add(expireTime)

	// Create claims
	claims := &middleware.AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID, // JWT ID (jti)
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:  user.ID,
		Address: user.Address,
		Role:    user.Role,
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate the JWT
	tokenString, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	// Store token in Valkey
	err = s.valkeyClient.StoreValidToken(ctx, user.ID, tokenID, expireTime)
	if err != nil {
		return "", time.Time{}, err
	}

	// Update user's last login time
	now := time.Now()
	user.LastLogin = &now
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		// Non-critical error, we can still return the token
		// but should log this error
	}

	return tokenString, expirationTime, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	// Parse the JWT
	claims := &middleware.AuthClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check if token is still valid in Valkey
	valid, err := s.valkeyClient.IsValidToken(ctx, claims.UserID, claims.RegisteredClaims.ID)
	if err != nil {
		// If Valkey is down, we'll fall back to just checking if the user exists and is not deleted
		// Log this error but continue for availability reasons
	} else if !valid {
		return nil, errors.New("token has been revoked")
	}

	// Get user from database
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if user is deleted
	if user.DeletedAt.Valid {
		return nil, errors.New("user account has been deleted")
	}

	return user, nil
}

// InvalidateToken removes a specific token from valid tokens
func (s *authService) InvalidateToken(ctx context.Context, userID uint, tokenID string) error {
	return s.valkeyClient.InvalidateToken(ctx, userID, tokenID)
}

// InvalidateAllUserTokens removes all valid tokens for a user
func (s *authService) InvalidateAllUserTokens(ctx context.Context, userID uint) error {
	return s.valkeyClient.InvalidateAllUserTokens(ctx, userID)
}
