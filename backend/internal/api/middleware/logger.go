package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Logger returns a configured logging middleware
func Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} ${method} ${path} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	})
}

// RequestLogger middleware logs detailed information about each request
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Store start time
		start := time.Now()

		// Get user info if authenticated
		var userInfo string
		if address, ok := c.Locals("address").(string); ok && address != "" {
			userInfo = " | User: " + address
		}

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log the request using standard Go logger
		log.Printf(
			"%s | %s %s%s | Status: %d | Duration: %v | IP: %s | User-Agent: %s",
			c.Context().RemoteAddr(),
			c.Method(),
			c.Path(),
			userInfo,
			c.Response().StatusCode(),
			duration,
			c.IP(),
			c.Get("User-Agent"),
		)

		return err
	}
}
