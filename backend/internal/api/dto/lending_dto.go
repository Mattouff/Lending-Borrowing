package dto

// LendingInfoResponse represents lending information for a user
type LendingInfoResponse struct {
	TotalDeposited      string `json:"totalDeposited"`
	InterestEarned      string `json:"interestEarned"`
	CurrentInterestRate string `json:"currentInterestRate"`
}

// BorrowingInfoResponse represents borrowing information for a user
type BorrowingInfoResponse struct {
	TotalBorrowed       string `json:"totalBorrowed"`
	InterestAccrued     string `json:"interestAccrued"`
	CurrentInterestRate string `json:"currentInterestRate"`
}

// CollateralInfoResponse represents collateral information for a user
type CollateralInfoResponse struct {
	TotalCollateral    string `json:"totalCollateral"`
	CollateralRatio    string `json:"collateralRatio"`
	MinCollateralRatio string `json:"minCollateralRatio"`
	MaxBorrowable      string `json:"maxBorrowable"`
	IsAtRisk           bool   `json:"isAtRisk"`
}

// LiquidatablePositionResponse represents a position that can be liquidated
type LiquidatablePositionResponse struct {
	PositionID       uint   `json:"positionId"`
	UserAddress      string `json:"userAddress"`
	CollateralAmount string `json:"collateralAmount"`
	CollateralToken  string `json:"collateralToken"`
	BorrowedAmount   string `json:"borrowedAmount"`
	BorrowedToken    string `json:"borrowedToken"`
	HealthFactor     string `json:"healthFactor"`
	LiquidationBonus string `json:"liquidationBonus"`
}

// LiquidatablePositionsResponse represents positions that can be liquidated
type LiquidatablePositionsResponse struct {
	Positions []LiquidatablePositionResponse `json:"positions"`
}
