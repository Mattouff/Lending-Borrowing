// @title Lending & Borrowing Platform API
// @version 1.0
// @description API for decentralized lending and borrowing platform
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Mattouff/Lending-Borrowing/internal/api/dto"
	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// UserHandler manages user-related API endpoints
type UserHandler struct {
	userService service.UserService
	config      *config.Config
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		config:      cfg,
	}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user with Ethereum address
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.UserRegistrationRequest true "User registration details"
// @Success 201 {object} dto.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.UserRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate the request
	if req.Address == "" || req.Username == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Address and username are required")
	}

	// Create the user
	user, err := h.userService.Register(c.Context(), req.Address, req.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to register user: "+err.Error())
	}

	// Convert to response DTO
	response := dto.UserResponse{
		ID:        user.ID,
		Address:   user.Address,
		Username:  user.Username,
		Role:      dto.UserRole(user.Role),
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.LastLogin != nil {
		response.LastLogin = *user.LastLogin
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    response,
	})
}

// Authenticate godoc
// @Summary Authenticate with signature
// @Description Authenticate user with Ethereum signature
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.UserAuthRequest true "Authentication details"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/auth [post]
func (h *UserHandler) Authenticate(c *fiber.Ctx) error {
	var req dto.UserAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate the request
	if req.Address == "" || req.Signature == "" || req.Message == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Address, signature, and message are required")
	}

	// Verify the signature
	isValid, err := h.userService.VerifySignature(c.Context(), req.Address, req.Message, req.Signature)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to verify signature: "+err.Error())
	}

	if !isValid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid signature")
	}

	// Get the user
	user, err := h.userService.GetByAddress(c.Context(), req.Address)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user: "+err.Error())
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// Generate token
	token, err := h.userService.GenerateAuthToken(c.Context(), user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token: "+err.Error())
	}

	// Convert to response DTO
	userResponse := dto.UserResponse{
		ID:        user.ID,
		Address:   user.Address,
		Username:  user.Username,
		Role:      dto.UserRole(user.Role),
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.LastLogin != nil {
		userResponse.LastLogin = *user.LastLogin
	}

	return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
		User:  userResponse,
		Token: token,
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieves the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=dto.UserResponse}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by authentication middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	// Get the user
	user, err := h.userService.GetByID(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user: "+err.Error())
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// Convert to response DTO
	response := dto.UserResponse{
		ID:        user.ID,
		Address:   user.Address,
		Username:  user.Username,
		Role:      dto.UserRole(user.Role),
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.LastLogin != nil {
		response.LastLogin = *user.LastLogin
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Data:    response,
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UserUpdateRequest true "User profile update details"
// @Success 200 {object} dto.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by authentication middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	var req dto.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Get the current user
	user, err := h.userService.GetByID(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user: "+err.Error())
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// Update fields if provided
	if req.Username != "" {
		user.Username = req.Username
	}

	user.UpdatedAt = time.Now()

	// Save the updated user
	if err := h.userService.Update(c.Context(), user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update profile: "+err.Error())
	}

	// Convert to response DTO
	response := dto.UserResponse{
		ID:        user.ID,
		Address:   user.Address,
		Username:  user.Username,
		Role:      dto.UserRole(user.Role),
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.LastLogin != nil {
		response.LastLogin = *user.LastLogin
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    response,
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieves a user by ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} dto.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/admin/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	// Get ID from URL
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	// Get the user
	user, err := h.userService.GetByID(c.Context(), uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user: "+err.Error())
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// Convert to response DTO
	response := dto.UserResponse{
		ID:        user.ID,
		Address:   user.Address,
		Username:  user.Username,
		Role:      dto.UserRole(user.Role),
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.LastLogin != nil {
		response.LastLogin = *user.LastLogin
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetUserByAddress godoc
// @Summary Get user by address
// @Description Retrieves a user by Ethereum address (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param address path string true "Ethereum address"
// @Success 200 {object} dto.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/admin/address/{address} [get]
func (h *UserHandler) GetUserByAddress(c *fiber.Ctx) error {
	// Get address from URL
	address := c.Params("address")
	if address == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Address is required")
	}

	// Get the user
	user, err := h.userService.GetByAddress(c.Context(), address)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user: "+err.Error())
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// Convert to response DTO
	response := dto.UserResponse{
		ID:        user.ID,
		Address:   user.Address,
		Username:  user.Username,
		Role:      dto.UserRole(user.Role),
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.LastLogin != nil {
		response.LastLogin = *user.LastLogin
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Data:    response,
	})
}

// ListUsers godoc
// @Summary List all users
// @Description Lists all users with pagination (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} dto.UserListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/admin [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get users
	users, err := h.userService.ListUsers(c.Context(), offset, pageSize)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to list users: "+err.Error())
	}

	// Convert to response DTOs
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponse := dto.UserResponse{
			ID:        user.ID,
			Address:   user.Address,
			Username:  user.Username,
			Role:      dto.UserRole(user.Role),
			Verified:  user.Verified,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		if user.LastLogin != nil {
			userResponse.LastLogin = *user.LastLogin
		}

		userResponses[i] = userResponse
	}

	// Get total count
	total, err := h.userService.Count(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to count users: "+err.Error())
	}

	// Calculate total pages
	totalPages := (int(total) + pageSize - 1) / pageSize

	return c.Status(fiber.StatusOK).JSON(dto.UserListResponse{
		Users:     userResponses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPages,
	})
}

// VerifyUser godoc
// @Summary Verify user
// @Description Marks a user as verified (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/admin/{id}/verify [put]
func (h *UserHandler) VerifyUser(c *fiber.Ctx) error {
	// Get ID from URL
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	// Get the user
	user, err := h.userService.GetByID(c.Context(), uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user: "+err.Error())
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// Update verification status
	user.Verified = true
	user.UpdatedAt = time.Now()

	if err := h.userService.Update(c.Context(), user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to verify user: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "User verified successfully",
	})
}

// NonceMessage godoc
// @Summary Get nonce message
// @Description Gets a nonce message for the user to sign for authentication
// @Tags users
// @Accept json
// @Produce json
// @Param address path string true "Ethereum address"
// @Success 200 {object} object{message=string,nonce=string}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/nonce/{address} [get]
func (h *UserHandler) NonceMessage(c *fiber.Ctx) error {
	// Get address from URL
	address := c.Params("address")
	if address == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Address is required")
	}

	// Get or create the user
	user, err := h.userService.GetByAddress(c.Context(), address)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user: "+err.Error())
	}

	var nonce string
	if user == nil {
		// User doesn't exist yet, create a temporary nonce
		nonce = h.userService.GenerateNonce(address)
	} else {
		nonce = user.Nonce
		if nonce == "" {
			// Generate a new nonce if none exists
			nonce = h.userService.GenerateNonce(address)
			user.Nonce = nonce
			if err := h.userService.Update(c.Context(), user); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user nonce: "+err.Error())
			}
		}
	}

	// Create the message to be signed
	message := "Sign this message to authenticate with our platform. Nonce: " + nonce

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"nonce":   nonce,
	})
}
