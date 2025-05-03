package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/swaggo/fiber-swagger"
	"gorm.io/gorm"

	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/Mattouff/Lending-Borrowing/internal/infrastructure/blockchain"
	"github.com/Mattouff/Lending-Borrowing/pkg/database"
)

// @title Lending & Borrowing Platform API
// @version 1.0
// @description API for decentralized lending and borrowing platform
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load environment variables and configuration
	if err := config.LoadEnv(""); err != nil {
		log.Printf("Warning: %v", err)
		log.Println("Continuing without .env file")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.Connect(cfg.Database.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.MigrateDB(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize Ethereum client
	ethClient := blockchain.GetInstance()
	if err := ethClient.Initialize(
		cfg.Blockchain.RpcURL,
		cfg.Blockchain.NetworkName,
		cfg.Blockchain.ContractAddresses,
	); err != nil {
		log.Fatalf("Failed to initialize blockchain client: %v", err)
	}
	defer ethClient.Close()

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
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(app.Listen(serverAddr))
}

func setupRoutes(api fiber.Router, db *gorm.DB) {
	//routes.SetupUserRoutes(api, db)
	// Add more routes as needed
}
