package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"slices"

	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// AuthClaims extends JWT standard claims with user information
type AuthClaims struct {
	jwt.RegisteredClaims
	UserID  uint            `json:"id"`
	Address string          `json:"address"`
	Role    models.UserRole `json:"role"`
}

// Authentication middleware to verify JWT tokens
func Authentication(cfg *config.Config, authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Authorization header is missing")
		}

		// Check if the authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
		}

		// Get the token
		tokenString := parts[1]

        // Validate token using auth service
        user, err := authService.ValidateToken(c.Context(), tokenString)
        if err != nil {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
        }

		// Store user information in the context
		c.Locals("userID", user.ID)
		c.Locals("address", user.Address)
		c.Locals("role", user.Role)

		return c.Next()
	}
}

// RoleAuthorization middleware to check if user has required role
func RoleAuthorization(roles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(models.UserRole)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Role information not found")
		}

		// Check if user role is in the allowed roles
		if slices.Contains(roles, userRole) {
				return c.Next()
			}

		return fiber.NewError(fiber.StatusForbidden, "Access denied")
	}
}
