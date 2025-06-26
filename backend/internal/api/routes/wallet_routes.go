package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/handlers"
)

// SetupWalletRoutes sets up the wallet-based transaction routes
func SetupWalletRoutes(app *fiber.App, walletHandler *handlers.WalletHandler, authMiddleware fiber.Handler) {
	// Wallet group with authentication
	wallet := app.Group("/api/v1/wallet", authMiddleware)

	// Lending wallet routes
	lending := wallet.Group("/lending")
	lending.Post("/prepare-deposit", walletHandler.PrepareDeposit)
	lending.Post("/prepare-withdraw", walletHandler.PrepareWithdraw)

	// TODO: Borrowing wallet routes (to be implemented)
	// borrowing := wallet.Group("/borrowing")
	// borrowing.Post("/prepare-borrow", walletHandler.PrepareBorrow)
	// borrowing.Post("/prepare-repay", walletHandler.PrepareRepay)

	// TODO: Collateral wallet routes (to be implemented)
	// collateral := wallet.Group("/collateral")
	// collateral.Post("/prepare-deposit", walletHandler.PrepareCollateralDeposit)
	// collateral.Post("/prepare-withdraw", walletHandler.PrepareCollateralWithdraw)

	// General wallet routes
	wallet.Post("/submit-transaction", walletHandler.SubmitTransaction)
	wallet.Get("/transaction/:hash/status", walletHandler.GetTransactionStatus)
}
