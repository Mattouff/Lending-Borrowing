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

// LendingHandler manages lending-related API endpoints
type LendingHandler struct {
	lendingService service.LendingService
}

// NewLendingHandler creates a new lending handler
func NewLendingHandler(lendingService service.LendingService) *LendingHandler {
	return &LendingHandler{
		lendingService: lendingService,
	}
}

// Deposit handles token deposits to the lending pool
func (h *LendingHandler) Deposit(c *fiber.Ctx) error {
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

	// Call the lending service to make the deposit
	txHash, err := h.lendingService.Deposit(c.Context(), common.HexToAddress(address), amount)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process deposit: "+err.Error())
	}

	// Return the transaction hash
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Deposit transaction submitted",
		Data: fiber.Map{
			"transactionHash": txHash,
		},
	})
}

// Withdraw handles token withdrawals from the lending pool
func (h *LendingHandler) Withdraw(c *fiber.Ctx) error {
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

	// Call the lending service to make the withdrawal
	txHash, err := h.lendingService.Withdraw(c.Context(), common.HexToAddress(address), amount)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process withdrawal: "+err.Error())
	}

	// Return the transaction hash
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Withdrawal transaction submitted",
		Data: fiber.Map{
			"transactionHash": txHash,
		},
	})
}

// GetLendingBalance returns the user's balance in the lending pool
func (h *LendingHandler) GetLendingBalance(c *fiber.Ctx) error {
	// Extract the user address from the authentication middleware
	address, ok := c.Locals("address").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Call the lending service to get the balance
	balance, err := h.lendingService.GetUserBalance(c.Context(), common.HexToAddress(address))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get lending balance: "+err.Error())
	}

	// Return the balance
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Data: fiber.Map{
			"balance": balance.String(),
		},
	})
}

// GetLendingInfo returns detailed lending information for the user
func (h *LendingHandler) GetLendingInfo(c *fiber.Ctx) error {
	// Extract the user address from the authentication middleware
	address, ok := c.Locals("address").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Get the user's balance
	balance, err := h.lendingService.GetUserBalance(c.Context(), common.HexToAddress(address))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get lending balance: "+err.Error())
	}

	// Get the current interest rate
	interestRate, err := h.lendingService.GetCurrentInterestRate(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get interest rate: "+err.Error())
	}

	// Get the interest earned by the user
	interestEarned, err := h.lendingService.GetUserInterestEarned(c.Context(), common.HexToAddress(address))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get interest earned: "+err.Error())
	}

	// Return the lending information
	return c.Status(fiber.StatusOK).JSON(dto.LendingInfoResponse{
		TotalDeposited:      balance.String(),
		InterestEarned:      interestEarned.String(),
		CurrentInterestRate: interestRate.String(),
	})
}

// GetTransactionHistory returns the lending transaction history for a user
func (h *LendingHandler) GetTransactionHistory(c *fiber.Ctx) error {
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
	transactions, err := h.lendingService.GetUserTransactionHistory(
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
		"type": []models.TransactionType{models.TransactionDeposit, models.TransactionWithdraw},
	}

	// Get total count from service
	total, err := h.lendingService.CountUserTransactions(c.Context(), ethAddress, filter)
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

// GetPoolInfo returns information about the lending pool
func (h *LendingHandler) GetPoolInfo(c *fiber.Ctx) error {
	// Get the total deposited amount
	totalDeposited, err := h.lendingService.GetTotalDeposited(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get total deposited: "+err.Error())
	}

	// Get the current interest rate
	interestRate, err := h.lendingService.GetCurrentInterestRate(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get interest rate: "+err.Error())
	}

	// Return the pool information
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Data: fiber.Map{
			"totalDeposited": totalDeposited.String(),
			"interestRate":   interestRate.String(),
		},
	})
}
