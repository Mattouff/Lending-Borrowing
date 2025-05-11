package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/handlers"
	"github.com/Mattouff/Lending-Borrowing/internal/api/middleware"
	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// SetupLendingRoutes configures the routes for lending operations
func SetupLendingRoutes(router fiber.Router, lendingService service.LendingService, authService service.AuthService, cfg *config.Config) {
	// Create handler
	lendingHandler := handlers.NewLendingHandler(lendingService)

	// Lending routes
	lendingRouter := router.Group("/lending")

	// Public routes
	lendingRouter.Get("/pool-info", lendingHandler.GetPoolInfo)

	// Protected routes (require authentication)
	lendingRouter.Use(middleware.Authentication(cfg, authService))
	lendingRouter.Post("/deposit", lendingHandler.Deposit)
	lendingRouter.Post("/withdraw", lendingHandler.Withdraw)
	lendingRouter.Get("/balance", lendingHandler.GetLendingBalance)
	lendingRouter.Get("/info", lendingHandler.GetLendingInfo)
	lendingRouter.Get("/transactions", lendingHandler.GetTransactionHistory)
}
