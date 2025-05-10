package dto

import "time"

// PositionStatus for DTO
type PositionStatus string

const (
	PositionStatusActive     PositionStatus = "active"
	PositionStatusLiquidated PositionStatus = "liquidated"
	PositionStatusClosed     PositionStatus = "closed"
)

// PositionResponse represents a position in API responses
type PositionResponse struct {
	ID               uint           `json:"id"`
	UserID           uint           `json:"userId"`
	CollateralAmount string         `json:"collateralAmount"`
	CollateralToken  string         `json:"collateralToken"`
	BorrowedAmount   string         `json:"borrowedAmount"`
	BorrowedToken    string         `json:"borrowedToken"`
	InterestRate     string         `json:"interestRate"`
	Status           PositionStatus `json:"status"`
	HealthFactor     string         `json:"healthFactor"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}

// PositionListResponse represents a list of positions for API responses
type PositionListResponse struct {
	Positions []PositionResponse `json:"positions"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"pageSize"`
	TotalPage int                `json:"totalPage"`
}
