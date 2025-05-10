package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"slices"
)

// AuthClaims extends JWT standard claims with user information
type AuthClaims struct {
	jwt.RegisteredClaims
	UserID  uint            `json:"id"`
	Address string          `json:"address"`
	Role    models.UserRole `json:"role"`
}

// Authentication middleware to verify JWT tokens
func Authentication(cfg *config.Config) fiber.Handler {
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

		// Parse and validate the token
		claims := &AuthClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			// Validate the signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			// Return the secret key
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
		}

		if !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		// Store user information in the context
		c.Locals("userID", claims.UserID)
		c.Locals("address", claims.Address)
		c.Locals("role", claims.Role)

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
