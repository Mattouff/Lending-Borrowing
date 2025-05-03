package routes

import (
    "github.com/gofiber/fiber/v2"

    "github.com/Mattouff/Lending-Borrowing/internal/config"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(app *fiber.App, services *Services, repositories *Repositories, cfg *config.Config) {
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