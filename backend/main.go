package main

import (
	"fmt"
	"log"

	_ "github.com/Mattouff/Lending-Borrowing/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"github.com/Mattouff/Lending-Borrowing/internal/api/middleware"
	"github.com/Mattouff/Lending-Borrowing/internal/api/routes"
	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/Mattouff/Lending-Borrowing/internal/infrastructure/blockchain"
	"github.com/Mattouff/Lending-Borrowing/internal/infrastructure/persistence/postgres"
	"github.com/Mattouff/Lending-Borrowing/internal/service"
	valkey "github.com/Mattouff/Lending-Borrowing/pkg/cache"
	"github.com/Mattouff/Lending-Borrowing/pkg/database"
)

// @title Lending & Borrowing Platform API
// @version 1.0
// @description API for decentralized lending and borrowing platform
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	// Initialize Valkey client
	valkeyClient, err := valkey.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Valkey client: %v", err)
	}
	defer valkeyClient.Close()

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

	// Create repository factory
	repoFactory := postgres.NewRepositoryFactory(db)

	// Initialize repositories
	userRepo := repoFactory.GetUserRepository()
	transactionRepo := repoFactory.GetTransactionRepository()
	positionRepo := repoFactory.GetPositionRepository()

	// Initialize services
	authService := service.NewAuthService(
		cfg,
		userRepo,
		valkeyClient,
	)

	userService := service.NewUserService(userRepo, cfg, authService)

	// Initialize collateral service first since borrowing service depends on it
	collateralService, err := service.NewCollateralService(
		transactionRepo,
		userRepo,
		positionRepo,
	)
	if err != nil {
		log.Fatalf("Failed to create collateral service: %v", err)
	}

	// Initialize other services
	lendingService, err := service.NewLendingService(
		transactionRepo,
		userRepo,
	)
	if err != nil {
		log.Fatalf("Failed to create lending service: %v", err)
	}

	borrowingService, err := service.NewBorrowingService(
		transactionRepo,
		userRepo,
		positionRepo,
		collateralService,
	)
	if err != nil {
		log.Fatalf("Failed to create borrowing service: %v", err)
	}

	liquidationService, err := service.NewLiquidationService(
		transactionRepo,
		userRepo,
		positionRepo,
		collateralService,
	)
	if err != nil {
		log.Fatalf("Failed to create liquidation service: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
	})

	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Create services container to pass to routes
	services := &routes.Services{
		UserService:        userService,
		LendingService:     lendingService,
		BorrowingService:   borrowingService,
		CollateralService:  collateralService,
		LiquidationService: liquidationService,
		AuthService:        authService,
		ValkeyClient:       valkeyClient,
	}

	// Create repositories container to pass to routes
	repositories := &routes.Repositories{
		UserRepository:        userRepo,
		PositionRepository:    positionRepo,
		TransactionRepository: transactionRepo,
	}

	// Setup routes with both services and repositories
	routes.SetupRoutes(app, services, repositories, cfg)

	// Setup Swagger documentation route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(app.Listen(serverAddr))
}
