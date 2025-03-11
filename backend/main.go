package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	_ "github.com/swaggo/fiber-swagger"
	"gorm.io/gorm"

	"github.com/Mattouff/Lending-Borrowing/api/routes"
	"github.com/Mattouff/Lending-Borrowing/pkg/database"
)

// @title Lending & Borrowing Platform API
// @version 1.0
// @description API for decentralized lending and borrowing platform
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.MigrateDB(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Set up routes
	api := app.Group("/api/v1")
	setupRoutes(api, db)

	// Start server
	log.Fatal(app.Listen(":8080"))
}

func setupRoutes(api fiber.Router, db *gorm.DB) {
	routes.SetupUserRoutes(api, db)
}
