package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/handlers"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// SetupLiquidationRoutes configures the routes for liquidation operations
func SetupLiquidationRoutes(router fiber.Router, liquidationService service.LiquidationService) {
	// Create handler
	liquidationHandler := handlers.NewLiquidationHandler(liquidationService)

	// Liquidation routes
	liquidationRouter := router.Group("/liquidation")

	// Public routes
	liquidationRouter.Get("/positions", liquidationHandler.GetLiquidatablePositions)
	liquidationRouter.Get("/history", liquidationHandler.GetLiquidationHistory)
	liquidationRouter.Get("/bonus", liquidationHandler.GetLiquidationBonus)

	// Protected routes (require authentication) - removed liquidate endpoint
	// Liquidators can interact directly with smart contracts for true decentralization
}
