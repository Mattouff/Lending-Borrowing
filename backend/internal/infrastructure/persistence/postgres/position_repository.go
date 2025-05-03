package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/repository"
)

type positionRepository struct {
	db *gorm.DB
}

// NewPositionRepository creates a new PostgreSQL implementation of PositionRepository
func NewPositionRepository(db *gorm.DB) repository.PositionRepository {
	return &positionRepository{
		db: db,
	}
}

// Create inserts a new position into the database
func (r *positionRepository) Create(ctx context.Context, position *models.Position) error {
	return r.db.WithContext(ctx).Create(position).Error
}

// FindByID retrieves a position by ID
func (r *positionRepository) FindByID(ctx context.Context, id uint) (*models.Position, error) {
	var position models.Position
	result := r.db.WithContext(ctx).First(&position, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &position, nil
}

// FindByUserID retrieves all positions for a specific user
func (r *positionRepository) FindByUserID(ctx context.Context, userID uint) ([]*models.Position, error) {
	var positions []*models.Position
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&positions).Error; err != nil {
		return nil, err
	}
	return positions, nil
}

// FindActiveByUserID retrieves all active positions for a specific user
func (r *positionRepository) FindActiveByUserID(ctx context.Context, userID uint) ([]*models.Position, error) {
	var positions []*models.Position
	if err := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, models.StatusActive).Find(&positions).Error; err != nil {
		return nil, err
	}
	return positions, nil
}

// Update updates an existing position
func (r *positionRepository) Update(ctx context.Context, position *models.Position) error {
	return r.db.WithContext(ctx).Save(position).Error
}

// UpdateStatus updates the status of a position
func (r *positionRepository) UpdateStatus(ctx context.Context, id uint, status models.PositionStatus) error {
	return r.db.WithContext(ctx).Model(&models.Position{}).Where("id = ?", id).Update("status", status).Error
}

// FindAtRisk finds all positions at risk of liquidation
func (r *positionRepository) FindAtRisk(ctx context.Context, healthFactorThreshold string) ([]*models.Position, error) {
	var positions []*models.Position

	// Find positions with a health factor below the threshold and that are still active
	if err := r.db.WithContext(ctx).
		Where("health_factor <= ? AND status = ?", healthFactorThreshold, models.StatusActive).
		Find(&positions).Error; err != nil {
		return nil, err
	}

	return positions, nil
}

// List retrieves all positions with optional filtering and pagination
func (r *positionRepository) List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*models.Position, error) {
	var positions []*models.Position
	query := r.db.WithContext(ctx).Model(&models.Position{})

	// Apply filters
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	if offset >= 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("created_at DESC").Find(&positions).Error; err != nil {
		return nil, err
	}

	return positions, nil
}

// Count returns the total number of positions matching the filter
func (r *positionRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Position{})

	// Apply filters
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	err := query.Count(&count).Error
	return count, err
}
