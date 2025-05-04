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

// Deposit godoc
// @Summary Deposit tokens to lending pool
// @Description Make a deposit into the lending pool to earn interest
// @Tags lending
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TransactionRequest true "Deposit amount"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /lending/deposit [post]
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

// Withdraw godoc
// @Summary Withdraw tokens from lending pool
// @Description Withdraw deposited tokens from the lending pool
// @Tags lending
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TransactionRequest true "Withdraw amount"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /lending/withdraw [post]
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

// GetLendingBalance godoc
// @Summary Get lending balance
// @Description Get user's current balance in the lending pool
// @Tags lending
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /lending/balance [get]
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

// GetLendingInfo godoc
// @Summary Get lending information
// @Description Get detailed lending information for the user including interest earned
// @Tags lending
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.LendingInfoResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /lending/info [get]
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

// GetTransactionHistory godoc
// @Summary Get lending transaction history
// @Description Get paginated history of user's lending transactions
// @Tags lending
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} dto.TransactionListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /lending/transactions [get]
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

// GetPoolInfo godoc
// @Summary Get lending pool information
// @Description Get information about the lending pool including total deposited and interest rate
// @Tags lending
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /lending/pool-info [get]
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
