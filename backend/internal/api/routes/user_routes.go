package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/handlers"
	"github.com/Mattouff/Lending-Borrowing/internal/api/middleware"
	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// SetupUserRoutes configures the routes for user management
func SetupUserRoutes(router fiber.Router, userService service.UserService, authService service.AuthService, cfg *config.Config) {
	// Create handler
	userHandler := handlers.NewUserHandler(userService, authService, cfg)

	// User routes
	userRouter := router.Group("/users")

	// Public routes
	userRouter.Post("/register", userHandler.Register)
	userRouter.Post("/auth", userHandler.Authenticate)
	userRouter.Get("/nonce/:address", userHandler.NonceMessage)

	// Protected routes (require authentication)
	userRouter.Use(middleware.Authentication(cfg, authService))
	userRouter.Get("/profile", userHandler.GetProfile)
	userRouter.Put("/profile", userHandler.UpdateProfile)
	userRouter.Delete("/account", userHandler.DeleteAccount)

	// Admin only routes
	adminRouter := userRouter.Group("/admin")
	adminRouter.Use(middleware.RoleAuthorization(models.RoleAdmin))
	adminRouter.Get("/", userHandler.ListUsers)
	adminRouter.Get("/:id", userHandler.GetUserByID)
	adminRouter.Get("/address/:address", userHandler.GetUserByAddress)
	adminRouter.Put("/:id/verify", userHandler.VerifyUser)
	adminRouter.Delete("/:id", userHandler.DeleteUser)
}
