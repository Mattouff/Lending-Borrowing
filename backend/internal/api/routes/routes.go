package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/config"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(app *fiber.App, services *Services, cfg *config.Config) {
	api := app.Group("/api/v1")

	// Setup individual route groups
	SetupUserRoutes(api, services.UserService, cfg)
}
