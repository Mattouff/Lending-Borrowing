package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/handlers"
	"github.com/Mattouff/Lending-Borrowing/internal/api/middleware"
	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// SetupLiquidationRoutes configures the routes for liquidation operations
func SetupLiquidationRoutes(router fiber.Router, liquidationService service.LiquidationService, authService service.AuthService, cfg *config.Config) {
	// Create handler
	liquidationHandler := handlers.NewLiquidationHandler(liquidationService)

	// Liquidation routes
	liquidationRouter := router.Group("/liquidation")

	// Public routes
	liquidationRouter.Get("/positions", liquidationHandler.GetLiquidatablePositions)
	liquidationRouter.Get("/history", liquidationHandler.GetLiquidationHistory)
	liquidationRouter.Get("/bonus", liquidationHandler.GetLiquidationBonus)

	// Protected routes (require authentication)
	liquidationRouter.Use(middleware.Authentication(cfg, authService))
	liquidationRouter.Post("/liquidate", liquidationHandler.Liquidate)
}
