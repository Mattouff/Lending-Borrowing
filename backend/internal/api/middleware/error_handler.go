package middleware

import (
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/Mattouff/Lending-Borrowing/internal/api/dto"
)

// ErrorHandler is a middleware that handles all errors in the application
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Default status code
		statusCode := fiber.StatusInternalServerError
		errorMsg := "Internal Server Error"

		// Check if the error is a Fiber error
		var fiberError *fiber.Error
		if errors.As(err, &fiberError) {
			statusCode = fiberError.Code
			errorMsg = fiberError.Message
		} else {
			// Handle specific error types
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				statusCode = fiber.StatusNotFound
				errorMsg = "Resource not found"
			case errors.Is(err, gorm.ErrDuplicatedKey):
				statusCode = fiber.StatusConflict
				errorMsg = "Resource already exists"
			}
		}

		// Check for validation errors
		var validateErr validator.ValidationErrors
		if errors.As(err, &validateErr) {
			statusCode = fiber.StatusBadRequest
			errorMsg = "Validation error"

			// Create a more user-friendly error message
			validationErrors := make(map[string]string)
			for _, e := range validateErr {
				validationErrors[e.Field()] = getValidationErrorMessage(e)
			}

			return c.Status(statusCode).JSON(fiber.Map{
				"error":       errorMsg,
				"validations": validationErrors,
			})
		}

		// Log the error if it's a server error
		if statusCode >= 500 {
			log.Printf("Server error: %v", err)
		}

		// Return a JSON response with the error
		return c.Status(statusCode).JSON(dto.ErrorResponse{
			Error: errorMsg,
		})
	}
}

// getValidationErrorMessage returns a user-friendly error message for a validation error
func getValidationErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value does not satisfy minimum length"
	case "max":
		return "Value exceeds maximum length"
	case "email":
		return "Invalid email format"
	case "eth_addr":
		return "Invalid Ethereum address"
	default:
		return "Invalid value"
	}
}
