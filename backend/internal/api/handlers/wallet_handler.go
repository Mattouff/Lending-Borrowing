package handlers

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/dto"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// WalletHandler manages wallet-based transaction endpoints
type WalletHandler struct {
	lendingService    service.LendingService
	borrowingService  service.BorrowingService
	collateralService service.CollateralService
}

// NewWalletHandler creates a new wallet handler
func NewWalletHandler(
	lendingService service.LendingService,
	borrowingService service.BorrowingService,
	collateralService service.CollateralService,
) *WalletHandler {
	return &WalletHandler{
		lendingService:    lendingService,
		borrowingService:  borrowingService,
		collateralService: collateralService,
	}
}

// PrepareDeposit godoc
// @Summary Prepare deposit transaction
// @Description Get contract interaction data for a deposit transaction
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TransactionRequest true "Deposit amount"
// @Success 200 {object} dto.APIResponse{data=service.ContractInteractionData}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /wallet/lending/prepare-deposit [post]
func (h *WalletHandler) PrepareDeposit(c *fiber.Ctx) error {
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

	// Prepare the transaction
	contractData, err := h.lendingService.PrepareDeposit(c.Context(), common.HexToAddress(address), amount)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to prepare deposit: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Deposit transaction prepared",
		Data:    contractData,
	})
}

// PrepareWithdraw godoc
// @Summary Prepare withdraw transaction
// @Description Get contract interaction data for a withdraw transaction
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TransactionRequest true "Withdraw amount"
// @Success 200 {object} dto.APIResponse{data=service.ContractInteractionData}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /wallet/lending/prepare-withdraw [post]
func (h *WalletHandler) PrepareWithdraw(c *fiber.Ctx) error {
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

	// Prepare the transaction
	contractData, err := h.lendingService.PrepareWithdraw(c.Context(), common.HexToAddress(address), amount)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to prepare withdraw: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Withdraw transaction prepared",
		Data:    contractData,
	})
}

// SubmitTransaction godoc
// @Summary Submit wallet-signed transaction
// @Description Submit a wallet-signed transaction for processing
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.WalletTransactionRequest true "Wallet transaction details"
// @Success 200 {object} dto.APIResponse{data=dto.WalletTransactionResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /wallet/submit-transaction [post]
func (h *WalletHandler) SubmitTransaction(c *fiber.Ctx) error {
	// Extract the user address from the authentication middleware
	userAddress, ok := c.Locals("address").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Parse the request body
	var req dto.WalletTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate that the user addresses match
	if req.UserAddress != userAddress {
		return fiber.NewError(fiber.StatusBadRequest, "User address mismatch")
	}

	// Convert the amount from string to *big.Int
	amount, success := new(big.Int).SetString(req.Amount, 10)
	if !success {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid amount format")
	}

	// Convert DTO transaction type to model transaction type
	var txType models.TransactionType
	switch req.Type {
	case dto.TransactionTypeDeposit:
		txType = models.TransactionDeposit
	case dto.TransactionTypeWithdraw:
		txType = models.TransactionWithdraw
	case dto.TransactionTypeBorrow:
		txType = models.TransactionBorrow
	case dto.TransactionTypeRepay:
		txType = models.TransactionRepay
	case dto.TransactionTypeLiquidate:
		txType = models.TransactionLiquidate
	default:
		return fiber.NewError(fiber.StatusBadRequest, "Invalid transaction type")
	}

	// Process the transaction based on transaction type
	var transaction *models.Transaction
	var err error

	switch txType {
	case models.TransactionDeposit, models.TransactionWithdraw:
		transaction, err = h.lendingService.ProcessWalletTransaction(
			c.Context(),
			req.TransactionHash,
			common.HexToAddress(userAddress),
			txType,
			amount,
		)
	case models.TransactionBorrow, models.TransactionRepay:
		// TODO: Implement borrowing service wallet transaction processing
		return fiber.NewError(fiber.StatusNotImplemented, "Borrowing wallet transactions not yet implemented")
	default:
		return fiber.NewError(fiber.StatusBadRequest, "Unsupported transaction type")
	}

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process transaction: "+err.Error())
	}

	response := dto.WalletTransactionResponse{
		TransactionID:   transaction.ID,
		TransactionHash: transaction.Hash,
		Status:          string(transaction.Status),
		Message:         "Transaction submitted successfully",
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Transaction processed successfully",
		Data:    response,
	})
}

// GetTransactionStatus godoc
// @Summary Get transaction status
// @Description Get the current status of a transaction
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param hash path string true "Transaction hash"
// @Success 200 {object} dto.APIResponse{data=map[string]interface{}}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /wallet/transaction/{hash}/status [get]
func (h *WalletHandler) GetTransactionStatus(c *fiber.Ctx) error {
	// Extract the user address from the authentication middleware
	userAddress, ok := c.Locals("address").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Get transaction hash from path
	txHash := c.Params("hash")
	if txHash == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Transaction hash is required")
	}

	// TODO: Implement transaction status checking from blockchain
	// For now, return a placeholder response
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Data: fiber.Map{
			"hash":   txHash,
			"status": "pending",
			"user":   userAddress,
		},
	})
}
