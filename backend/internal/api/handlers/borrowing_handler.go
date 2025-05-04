package handlers

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/dto"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// BorrowingHandler manages borrowing-related API endpoints
type BorrowingHandler struct {
	borrowingService service.BorrowingService
}

// NewBorrowingHandler creates a new borrowing handler
func NewBorrowingHandler(borrowingService service.BorrowingService) *BorrowingHandler {
	return &BorrowingHandler{
		borrowingService: borrowingService,
	}
}

// Borrow godoc
// @Summary Borrow tokens against collateral
// @Description Borrow tokens against deposited collateral
// @Tags borrowing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TransactionRequest true "Borrow amount"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /borrowing/borrow [post]
func (h *BorrowingHandler) Borrow(c *fiber.Ctx) error {
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

	// Call the borrowing service to process the borrow request
	txHash, err := h.borrowingService.Borrow(c.Context(), common.HexToAddress(address), amount)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process borrowing: "+err.Error())
	}

	// Return the transaction hash
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Borrow transaction submitted",
		Data: fiber.Map{
			"transactionHash": txHash,
		},
	})
}

// Repay godoc
// @Summary Repay borrowed tokens
// @Description Repay tokens previously borrowed from the protocol
// @Tags borrowing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TransactionRequest true "Repay amount"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /borrowing/repay [post]
func (h *BorrowingHandler) Repay(c *fiber.Ctx) error {
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

	// Call the borrowing service to process the repay request
	txHash, err := h.borrowingService.Repay(c.Context(), common.HexToAddress(address), amount)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process repayment: "+err.Error())
	}

	// Return the transaction hash
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Repay transaction submitted",
		Data: fiber.Map{
			"transactionHash": txHash,
		},
	})
}

// GetBorrowedAmount godoc
// @Summary Get borrowed amount
// @Description Get the total amount borrowed by the authenticated user
// @Tags borrowing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /borrowing/balance [get]
func (h *BorrowingHandler) GetBorrowedAmount(c *fiber.Ctx) error {
	// Extract the user address from the authentication middleware
	address, ok := c.Locals("address").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Call the borrowing service to get the borrowed amount
	borrowed, err := h.borrowingService.GetBorrowedAmount(c.Context(), common.HexToAddress(address))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get borrowed amount: "+err.Error())
	}

	// Return the borrowed amount
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Data: fiber.Map{
			"borrowedAmount": borrowed.String(),
		},
	})
}

// GetBorrowingInfo godoc
// @Summary Get borrowing information
// @Description Get detailed information about user's borrowing including interest accrued
// @Tags borrowing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.BorrowingInfoResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /borrowing/info [get]
func (h *BorrowingHandler) GetBorrowingInfo(c *fiber.Ctx) error {
	// Extract the user address from the authentication middleware
	address, ok := c.Locals("address").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Get the borrowed amount
	borrowed, err := h.borrowingService.GetBorrowedAmount(c.Context(), common.HexToAddress(address))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get borrowed amount: "+err.Error())
	}

	// Get the current interest rate
	interestRate, err := h.borrowingService.GetCurrentInterestRate(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get interest rate: "+err.Error())
	}

	// Get the interest accrued by the user
	interestAccrued, err := h.borrowingService.GetUserInterestAccrued(c.Context(), common.HexToAddress(address))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get interest accrued: "+err.Error())
	}

	// Return the borrowing information
	return c.Status(fiber.StatusOK).JSON(dto.BorrowingInfoResponse{
		TotalBorrowed:       borrowed.String(),
		InterestAccrued:     interestAccrued.String(),
		CurrentInterestRate: interestRate.String(),
	})
}

// GetTransactionHistory godoc
// @Summary Get borrowing transaction history
// @Description Get paginated history of user's borrowing transactions
// @Tags borrowing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} dto.TransactionListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /borrowing/transactions [get]
func (h *BorrowingHandler) GetTransactionHistory(c *fiber.Ctx) error {
	// Extract the user address from the authentication middleware
	address, ok := c.Locals("address").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

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

	// Get transaction history
	transactions, err := h.borrowingService.GetUserTransactionHistory(
		c.Context(),
		common.HexToAddress(address),
		offset,
		pageSize,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get transaction history: "+err.Error())
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

	// Get total count for pagination
	ethAddress := common.HexToAddress(address)
	filter := map[string]interface{}{
		"type": []models.TransactionType{models.TransactionBorrow, models.TransactionRepay},
	}

	// Get total count from service
	total, err := h.borrowingService.CountUserTransactions(c.Context(), ethAddress, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to count transactions: "+err.Error())
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

// GetBorrowingStats godoc
// @Summary Get borrowing statistics
// @Description Get global statistics about borrowing on the platform
// @Tags borrowing
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /borrowing/stats [get]
func (h *BorrowingHandler) GetBorrowingStats(c *fiber.Ctx) error {
	// Get the total borrowed amount
	totalBorrowed, err := h.borrowingService.GetTotalBorrowed(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get total borrowed: "+err.Error())
	}

	// Get the current interest rate
	interestRate, err := h.borrowingService.GetCurrentInterestRate(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get interest rate: "+err.Error())
	}

	// Return the borrowing statistics
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Data: fiber.Map{
			"totalBorrowed": totalBorrowed.String(),
			"interestRate":  interestRate.String(),
		},
	})
}
