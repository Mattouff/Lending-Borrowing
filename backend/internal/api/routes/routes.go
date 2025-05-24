package routes

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/config"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(app *fiber.App, services *Services, repositories *Repositories, cfg *config.Config) {
	// Health check endpoint (no auth needed)
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Health check for Valkey connection (no auth needed)
	app.Get("/api/health/valkey", func(c *fiber.Ctx) error {
		ctx := c.Context()
		valkeyClient := services.ValkeyClient

		start := time.Now()
		err := valkeyClient.Ping(ctx)
		latency := time.Since(start)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("Valkey connection failed: %v", err),
			})
		}

		return c.JSON(fiber.Map{
			"status":     "ok",
			"latency_ms": latency.Milliseconds(),
			"message":    "Valkey connection successful",
		})
	})

	api := app.Group("/api/v1")

	// Setup individual route groups
	SetupUserRoutes(api, services.UserService, services.AuthService, cfg)
	SetupLendingRoutes(api, services.LendingService, services.AuthService, cfg)
	SetupBorrowingRoutes(api, services.BorrowingService, services.AuthService, cfg)
	SetupCollateralRoutes(api, services.CollateralService, services.AuthService, cfg)
	SetupLiquidationRoutes(api, services.LiquidationService, services.AuthService, cfg)

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
