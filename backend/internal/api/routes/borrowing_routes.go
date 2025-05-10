package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/handlers"
	"github.com/Mattouff/Lending-Borrowing/internal/api/middleware"
	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// SetupBorrowingRoutes configures the routes for borrowing operations
func SetupBorrowingRoutes(router fiber.Router, borrowingService service.BorrowingService, cfg *config.Config) {
	// Create handler
	borrowingHandler := handlers.NewBorrowingHandler(borrowingService)

	// Borrowing routes
	borrowingRouter := router.Group("/borrowing")

	// Public routes
	borrowingRouter.Get("/stats", borrowingHandler.GetBorrowingStats)

	// Protected routes (require authentication)
	borrowingRouter.Use(middleware.Authentication(cfg))
	borrowingRouter.Post("/borrow", borrowingHandler.Borrow)
	borrowingRouter.Post("/repay", borrowingHandler.Repay)
	borrowingRouter.Get("/balance", borrowingHandler.GetBorrowedAmount)
	borrowingRouter.Get("/info", borrowingHandler.GetBorrowingInfo)
	borrowingRouter.Get("/transactions", borrowingHandler.GetTransactionHistory)
}
