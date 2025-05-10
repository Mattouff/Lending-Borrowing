package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/handlers"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/repository"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// SetupMarketRoutes configures the routes for market data
func SetupMarketRoutes(
	router fiber.Router,
	lendingService service.LendingService,
	borrowingService service.BorrowingService,
	collateralService service.CollateralService,
	userRepository repository.UserRepository,
	positionRepository repository.PositionRepository,
) {
	// Create handler with all required dependencies
	marketHandler := handlers.NewMarketHandler(
		lendingService,
		borrowingService,
		collateralService,
		userRepository,
		positionRepository,
	)

	// Market routes
	marketRouter := router.Group("/market")

	// Public routes
	marketRouter.Get("/overview", marketHandler.GetMarketOverview)
	marketRouter.Get("/tokens", marketHandler.GetTokensMarketData)
}
