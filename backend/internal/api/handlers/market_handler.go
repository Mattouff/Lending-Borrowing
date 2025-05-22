package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/dto"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/repository"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// MarketHandler manages market data API endpoints
type MarketHandler struct {
	lendingService     service.LendingService
	borrowingService   service.BorrowingService
	collateralService  service.CollateralService
	userRepository     repository.UserRepository
	positionRepository repository.PositionRepository
}

// NewMarketHandler creates a new market data handler
func NewMarketHandler(
	lendingService service.LendingService,
	borrowingService service.BorrowingService,
	collateralService service.CollateralService,
	userRepository repository.UserRepository,
	positionRepository repository.PositionRepository,
) *MarketHandler {
	return &MarketHandler{
		lendingService:     lendingService,
		borrowingService:   borrowingService,
		collateralService:  collateralService,
		userRepository:     userRepository,
		positionRepository: positionRepository,
	}
}

// GetMarketOverview godoc
// @Summary Get market overview
// @Description Get overview of the market including TVL, borrowed amounts, and rates
// @Tags market
// @Accept json
// @Produce json
// @Success 200 {object} dto.MarketOverviewResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /market/overview [get]
func (h *MarketHandler) GetMarketOverview(c *fiber.Ctx) error {
	// Get total deposited amount
	totalDeposited, err := h.lendingService.GetTotalDeposited(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get total deposited: "+err.Error())
	}

	// Get total borrowed amount
	totalBorrowed, err := h.borrowingService.GetTotalBorrowed(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get total borrowed: "+err.Error())
	}

	// Get lending interest rate
	lendingRate, err := h.lendingService.GetCurrentInterestRate(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get lending interest rate: "+err.Error())
	}

	// Get borrowing interest rate
	borrowingRate, err := h.borrowingService.GetCurrentInterestRate(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get borrowing interest rate: "+err.Error())
	}

	// Get active users count (users who have logged in within the last 30 days)
	activeUsersCount, err := h.getActiveUsersCount(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get active users count: "+err.Error())
	}

	// Get active positions count
	activePositionsCount, err := h.getActivePositionsCount(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get active positions count: "+err.Error())
	}

	// Return the market overview
	return c.Status(fiber.StatusOK).JSON(dto.MarketOverviewResponse{
		TotalValueLocked:    totalDeposited.String(),
		TotalBorrowed:       totalBorrowed.String(),
		ActiveUsers:         activeUsersCount,
		ActivePositions:     activePositionsCount,
		AverageLendingAPY:   lendingRate.String(),
		AverageBorrowingAPY: borrowingRate.String(),
	})
}

// getActiveUsersCount returns the count of users who have logged in within the last 30 days
func (h *MarketHandler) getActiveUsersCount(ctx context.Context) (int64, error) {
	// In a real implementation, you'd query the database for users with last_login within the last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	// Using a filter to count users who have logged in recently
	filter := map[string]any{
		"last_login_after": thirtyDaysAgo,
	}

	return h.userRepository.CountWithFilter(ctx, filter)
}

// getActivePositionsCount returns the count of active positions
func (h *MarketHandler) getActivePositionsCount(ctx context.Context) (int64, error) {
	// Using a filter to count positions with status "active"
	filter := map[string]any{
		"status": models.StatusActive,
	}

	// Count positions that match the filter
	return h.positionRepository.Count(ctx, filter)
}

// GetTokensMarketData godoc
// @Summary Get tokens market data
// @Description Get detailed market data for all supported tokens
// @Tags market
// @Accept json
// @Produce json
// @Success 200 {object} dto.TokensMarketResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /market/tokens [get]
func (h *MarketHandler) GetTokensMarketData(c *fiber.Ctx) error {
	// In a real implementation, you'd get token data from your services or a dedicated token service
	// This is just a placeholder implementation

	// Return a placeholder response
	return c.Status(fiber.StatusOK).JSON(dto.TokensMarketResponse{
		Tokens: []dto.TokenMarketData{
			{
				Token: dto.TokenMetadata{
					Address:  "0x...", // Replace with actual token address
					Symbol:   "ETH",
					Name:     "Ethereum",
					Decimals: 18,
					PriceUSD: "1000.00",
					LogoURI:  "https://example.com/eth-logo.png",
				},
				TotalSupply:        "1000000000000000000000",
				TotalDeposited:     "500000000000000000000",
				TotalBorrowed:      "200000000000000000000",
				AvailableLiquidity: "300000000000000000000",
				LendingAPY:         "500000000000000000",  // 0.5%
				BorrowingAPY:       "1000000000000000000", // 1%
				CollateralFactor:   "750000000000000000",  // 75%
			},
		},
	})
}
