package controllers

import (
	"github.com/Mattouff/Lending-Borrowing/internal/models"
	"github.com/Mattouff/Lending-Borrowing/internal/repository"
	"github.com/Mattouff/Lending-Borrowing/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserController handles user-related HTTP requests
type UserController struct {
	DB            *gorm.DB
	UserRepo      *repository.UserRepository
	WalletService *services.WalletService
}

// NewUserController creates a new instance of UserController
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		DB:            db,
		UserRepo:      repository.NewUserRepository(db),
		WalletService: services.NewWalletService(),
	}
}

// GetUser retrieves a user by wallet address
// @Summary Get user by wallet address
// @Description Get user details by Ethereum wallet address
// @Tags users
// @Accept json
// @Produce json
// @Param walletAddress path string true "Ethereum wallet address"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users/{walletAddress} [get]
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	walletAddress := ctx.Params("walletAddress")

	user, err := c.UserRepo.FindByWalletAddress(walletAddress)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return ctx.JSON(user)
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Register a new user with an Ethereum wallet address
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User information"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Generate a new nonce for wallet authentication
	nonce, err := c.WalletService.GenerateNonce()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate nonce",
		})
	}
	user.Nonce = nonce

	if err := c.UserRepo.Create(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot create user",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

// ListUsers retrieves all users
// @Summary List all users
// @Description Get a list of all registered users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (c *UserController) ListUsers(ctx *fiber.Ctx) error {
	var users []models.User
	result := c.DB.Find(&users)
	
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	
	return ctx.JSON(users)
}
