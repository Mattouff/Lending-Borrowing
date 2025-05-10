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

// Liquidate godoc
// @Summary Liquidate under-collateralized position
// @Description Liquidate a position that has fallen below minimum health factor
// @Tags liquidation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TransactionLiquidationRequest true "Liquidation parameters"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /liquidation/liquidate [post]
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

// GetLiquidatablePositions godoc
// @Summary Get liquidatable positions
// @Description Get list of positions that can be liquidated
// @Tags liquidation
// @Accept json
// @Produce json
// @Success 200 {object} dto.LiquidatablePositionsResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /liquidation/positions [get]
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

// GetLiquidationHistory godoc
// @Summary Get liquidation history
// @Description Get paginated history of liquidation events
// @Tags liquidation
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} dto.TransactionListResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /liquidation/history [get]
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

// GetLiquidationBonus godoc
// @Summary Get liquidation bonus
// @Description Get current liquidation bonus percentage
// @Tags liquidation
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /liquidation/bonus [get]
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
