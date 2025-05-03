package handlers

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/dto"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// LiquidationHandler manages liquidation-related API endpoints
type LiquidationHandler struct {
	liquidationService service.LiquidationService
}

// NewLiquidationHandler creates a new liquidation handler
func NewLiquidationHandler(liquidationService service.LiquidationService) *LiquidationHandler {
	return &LiquidationHandler{
		liquidationService: liquidationService,
	}
}

// Liquidate handles liquidating an under-collateralized position
func (h *LiquidationHandler) Liquidate(c *fiber.Ctx) error {
	// Extract the liquidator address from the authentication middleware
	liquidatorAddress, ok := c.Locals("address").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Parse the request body
	var req dto.TransactionLiquidationRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Convert the amount from string to *big.Int
	amount, success := new(big.Int).SetString(req.Amount, 10)
	if !success {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid amount format")
	}

	// Call the liquidation service to liquidate the position
	txHash, err := h.liquidationService.Liquidate(
		c.Context(),
		common.HexToAddress(liquidatorAddress),
		common.HexToAddress(req.BorrowerAddress),
		amount,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process liquidation: "+err.Error())
	}

	// Return the transaction hash
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Liquidation transaction submitted",
		Data: fiber.Map{
			"transactionHash": txHash,
		},
	})
}

// GetLiquidatablePositions returns positions that can be liquidated
func (h *LiquidationHandler) GetLiquidatablePositions(c *fiber.Ctx) error {
	// Call the liquidation service to get liquidatable positions
	positions, err := h.liquidationService.GetLiquidatablePositions(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get liquidatable positions: "+err.Error())
	}

	// Get liquidation bonus
	bonus, err := h.liquidationService.GetLiquidationBonus(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get liquidation bonus: "+err.Error())
	}

	// Convert positions to DTOs
	posResponseList := make([]dto.LiquidatablePositionResponse, len(positions))
	for i, pos := range positions {
		posResponseList[i] = dto.LiquidatablePositionResponse{
			PositionID:       pos.ID,
			UserAddress:      "0x...", // This would need to be fetched from the user service
			CollateralAmount: pos.CollateralAmount,
			CollateralToken:  pos.CollateralToken,
			BorrowedAmount:   pos.BorrowedAmount,
			BorrowedToken:    pos.BorrowedToken,
			HealthFactor:     pos.HealthFactor,
			LiquidationBonus: bonus.String(),
		}
	}

	// Return the liquidatable positions
	return c.Status(fiber.StatusOK).JSON(dto.LiquidatablePositionsResponse{
		Positions: posResponseList,
	})
}

// GetLiquidationHistory returns the history of liquidation events
func (h *LiquidationHandler) GetLiquidationHistory(c *fiber.Ctx) error {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get liquidation history
	transactions, err := h.liquidationService.GetLiquidationHistory(c.Context(), offset, pageSize)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get liquidation history: "+err.Error())
	}

	// Convert transactions to DTOs
	txResponses := make([]dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		txResponses[i] = dto.TransactionResponse{
			ID:           tx.ID,
			UserID:       tx.UserID,
			Type:         dto.TransactionType(tx.Type),
			Status:       dto.TransactionStatus(tx.Status),
			Hash:         tx.Hash,
			Amount:       tx.Amount,
			TokenAddress: tx.TokenAddress,
			BlockNumber:  tx.BlockNumber,
			GasUsed:      tx.GasUsed,
			GasPrice:     tx.GasPrice,
			CreatedAt:    tx.CreatedAt,
			UpdatedAt:    tx.UpdatedAt,
		}
	}


	// Get total count from service
	total, err := h.liquidationService.CountLiquidations(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to count liquidations: "+err.Error())
	}

	// Calculate total pages
	totalPages := (int(total) + pageSize - 1) / pageSize

	return c.Status(fiber.StatusOK).JSON(dto.TransactionListResponse{
		Transactions: txResponses,
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
		TotalPage:    totalPages,
	})
}

// GetLiquidationBonus returns the current liquidation bonus
func (h *LiquidationHandler) GetLiquidationBonus(c *fiber.Ctx) error {
	// Call the liquidation service to get the liquidation bonus
	bonus, err := h.liquidationService.GetLiquidationBonus(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get liquidation bonus: "+err.Error())
	}

	// Return the liquidation bonus
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Data: fiber.Map{
			"liquidationBonus": bonus.String(),
		},
	})
}
