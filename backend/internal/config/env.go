package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv(filePath string) error {
	if filePath == "" {
		filePath = ".env"
	}

	// Check if file exists before attempting to load
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("env file does not exist: %s", filePath)
	}

	return godotenv.Load(filePath)
}

// GetEnv retrieves an environment variable or returns a default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvBool retrieves a boolean environment variable or returns a default
func GetEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return b
}

// GetEnvInt retrieves an integer environment variable or returns a default
func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return i
}

// GetEnvArray retrieves a comma-separated environment variable as a string array
func GetEnvArray(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return strings.Split(value, ",")
}

// GetRequiredEnv returns the value of the environment variable or panics if not set
func GetRequiredEnv(key string) string {
    value, exists := os.LookupEnv(key)
    if !exists || value == "" {
        panic(fmt.Sprintf("Required environment variable %s is not set", key))
    }
    return value
}
