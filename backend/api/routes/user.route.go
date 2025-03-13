package routes

import (
	"github.com/Mattouff/Lending-Borrowing/api/controllers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupUserRoutes sets up all user-related routes
func SetupUserRoutes(api fiber.Router, db *gorm.DB) {
	// Initialize the controller
	userController := controllers.NewUserController(db)

	// Create a user group
	users := api.Group("/users")

	// Define routes
	users.Get("/", userController.ListUsers)
	users.Post("/", userController.CreateUser)
	users.Get("/:walletAddress", userController.GetUser)
}
