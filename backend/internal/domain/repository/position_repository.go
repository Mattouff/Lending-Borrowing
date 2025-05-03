package repository

import (
	"context"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
)

// PositionRepository defines the interface for lending/borrowing position data access
type PositionRepository interface {
	// Create inserts a new position into the database
	Create(ctx context.Context, position *models.Position) error

	// FindByID retrieves a position by ID
	FindByID(ctx context.Context, id uint) (*models.Position, error)

	// FindByUserID retrieves all positions for a specific user
	FindByUserID(ctx context.Context, userID uint) ([]*models.Position, error)

	// FindActiveByUserID retrieves all active positions for a specific user
	FindActiveByUserID(ctx context.Context, userID uint) ([]*models.Position, error)

	// Update updates an existing position
	Update(ctx context.Context, position *models.Position) error

	// UpdateStatus updates the status of a position
	UpdateStatus(ctx context.Context, id uint, status models.PositionStatus) error

	// FindAtRisk finds all positions at risk of liquidation
	FindAtRisk(ctx context.Context, healthFactorThreshold string) ([]*models.Position, error)

	// List retrieves all positions with optional filtering and pagination
	List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*models.Position, error)

	// Count returns the total number of positions matching the filter
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}
