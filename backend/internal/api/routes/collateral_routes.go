package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/handlers"
	"github.com/Mattouff/Lending-Borrowing/internal/api/middleware"
	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// SetupCollateralRoutes configures the routes for collateral management
func SetupCollateralRoutes(router fiber.Router, collateralService service.CollateralService, authService service.AuthService, cfg *config.Config) {
	// Create handler
	collateralHandler := handlers.NewCollateralHandler(collateralService)

	// Collateral routes
	collateralRouter := router.Group("/collateral")

	// Protected routes (require authentication) - only read-only endpoints
	// Collateral deposit/withdraw transactions should be done directly via wallet
	collateralRouter.Use(middleware.Authentication(cfg, authService))
	collateralRouter.Get("/balance", collateralHandler.GetCollateralBalance)
	collateralRouter.Get("/info", collateralHandler.GetCollateralInfo)
}
