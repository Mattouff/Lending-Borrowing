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

// Register handles user registration
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.UserRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate the request
	if req.Address == "" || req.Username == "" || req.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Address, username, and email are required")
	}

	// Create the user
	user, err := h.userService.Register(c.Context(), req.Address, req.Username, req.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to register user: "+err.Error())
	}

	// Convert to response DTO
	response := dto.UserResponse{
		ID:        user.ID,
		Address:   user.Address,
		Username:  user.Username,
		Email:     user.Email,
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

// Authenticate handles user authentication
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
		Email:     user.Email,
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

// GetProfile retrieves the user's profile
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
		Email:     user.Email,
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

// UpdateProfile updates the user's profile
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

	if req.Email != "" {
		user.Email = req.Email
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
		Email:     user.Email,
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

// GetUserByID retrieves a user by ID (admin only)
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
		Email:     user.Email,
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

// GetUserByAddress retrieves a user by Ethereum address (admin only)
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
		Email:     user.Email,
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

// ListUsers lists all users (admin only)
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
			Email:     user.Email,
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

// VerifyUser marks a user as verified (admin only)
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

// NonceMessage generates a nonce message for a user to sign
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
