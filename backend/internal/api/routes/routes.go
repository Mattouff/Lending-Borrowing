package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/config"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(app *fiber.App, services *Services, repositories *Repositories, cfg *config.Config) {
	// Health check endpoint (no auth needed)
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	api := app.Group("/api/v1")

	// Setup individual route groups
	SetupUserRoutes(api, services.UserService, cfg)
	SetupLendingRoutes(api, services.LendingService, cfg)
	SetupBorrowingRoutes(api, services.BorrowingService, cfg)
	SetupCollateralRoutes(api, services.CollateralService, cfg)
	SetupLiquidationRoutes(api, services.LiquidationService, cfg)

	// Setup market routes (uses multiple services and repositories)
	SetupMarketRoutes(
		api,
		services.LendingService,
		services.BorrowingService,
		services.CollateralService,
		repositories.UserRepository,
		repositories.PositionRepository,
	)
}
