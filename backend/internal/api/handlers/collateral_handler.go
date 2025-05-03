package handlers

import (
	"math/big"

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

// DepositCollateral handles depositing tokens as collateral
func (h *CollateralHandler) DepositCollateral(c *fiber.Ctx) error {
	// Extract the user address from the authentication middleware
	address, ok := c.Locals("address").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Parse the request body
	var req dto.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Convert the amount from string to *big.Int
	amount, success := new(big.Int).SetString(req.Amount, 10)
	if !success {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid amount format")
	}

	// Call the collateral service to deposit collateral
	txHash, err := h.collateralService.DepositCollateral(c.Context(), common.HexToAddress(address), amount)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to deposit collateral: "+err.Error())
	}

	// Return the transaction hash
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Collateral deposit transaction submitted",
		Data: fiber.Map{
			"transactionHash": txHash,
		},
	})
}

// WithdrawCollateral handles withdrawing tokens from collateral
func (h *CollateralHandler) WithdrawCollateral(c *fiber.Ctx) error {
	// Extract the user address from the authentication middleware
	address, ok := c.Locals("address").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Parse the request body
	var req dto.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Convert the amount from string to *big.Int
	amount, success := new(big.Int).SetString(req.Amount, 10)
	if !success {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid amount format")
	}

	// Call the collateral service to withdraw collateral
	txHash, err := h.collateralService.WithdrawCollateral(c.Context(), common.HexToAddress(address), amount)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to withdraw collateral: "+err.Error())
	}

	// Return the transaction hash
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Collateral withdrawal transaction submitted",
		Data: fiber.Map{
			"transactionHash": txHash,
		},
	})
}

// GetCollateralBalance returns the collateral balance of a user
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

// GetCollateralInfo returns detailed collateral information for the user
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
