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

// Borrow handles borrowing tokens against collateral
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

// Repay handles repaying borrowed tokens
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

// GetBorrowedAmount returns the amount borrowed by the user
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

// GetBorrowingInfo returns detailed borrowing information for the user
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

// GetTransactionHistory returns the borrowing transaction history for a user
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

// GetBorrowingStats returns statistics about borrowing on the platform
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
