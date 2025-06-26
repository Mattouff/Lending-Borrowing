package handlers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/dto"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// CollateralHandler manages collateral-related API endpoints
type CollateralHandler struct {
	collateralService service.CollateralService
}

// NewCollateralHandler creates a new collateral handler
func NewCollateralHandler(collateralService service.CollateralService) *CollateralHandler {
	return &CollateralHandler{
		collateralService: collateralService,
	}
}

// Note: Collateral deposit and withdrawal methods have been removed.
// Users should interact directly with the collateral contract using their wallets
// for true decentralization. The backend now only provides read-only information.

// GetCollateralBalance godoc
// @Summary Get collateral balance
// @Description Get user's current collateral balance
// @Tags collateral
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /collateral/balance [get]
func (h *CollateralHandler) GetCollateralBalance(c *fiber.Ctx) error {
	// Extract the user address from the authentication middleware
	address, ok := c.Locals("address").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Call the collateral service to get the collateral balance
	balance, err := h.collateralService.GetCollateralBalance(c.Context(), common.HexToAddress(address))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get collateral balance: "+err.Error())
	}

	// Return the collateral balance
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Data: fiber.Map{
			"collateralBalance": balance.String(),
		},
	})
}

// GetCollateralInfo godoc
// @Summary Get collateral information
// @Description Get detailed information about user's collateral position
// @Tags collateral
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.CollateralInfoResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /collateral/info [get]
func (h *CollateralHandler) GetCollateralInfo(c *fiber.Ctx) error {
	// Extract the user address from the authentication middleware
	address, ok := c.Locals("address").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Get the user's collateral balance
	balance, err := h.collateralService.GetCollateralBalance(c.Context(), common.HexToAddress(address))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get collateral balance: "+err.Error())
	}

	// Get the collateral ratio
	ratio, err := h.collateralService.GetCollateralRatio(c.Context(), common.HexToAddress(address))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get collateral ratio: "+err.Error())
	}

	// Get the minimum collateral ratio
	minRatio, err := h.collateralService.GetMinCollateralRatio(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get minimum collateral ratio: "+err.Error())
	}

	// Get the maximum borrowable amount
	maxBorrowable, err := h.collateralService.GetMaxBorrowableAmount(c.Context(), common.HexToAddress(address))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get max borrowable amount: "+err.Error())
	}

	// Check if the position is at risk
	isAtRisk, err := h.collateralService.IsAtRisk(c.Context(), common.HexToAddress(address))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check risk status: "+err.Error())
	}

	// Return the collateral information
	return c.Status(fiber.StatusOK).JSON(dto.CollateralInfoResponse{
		TotalCollateral:    balance.String(),
		CollateralRatio:    ratio.String(),
		MinCollateralRatio: minRatio.String(),
		MaxBorrowable:      maxBorrowable.String(),
		IsAtRisk:           isAtRisk,
	})
}
