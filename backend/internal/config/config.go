package config

import (
	"fmt"
	"strings"

	"github.com/Mattouff/Lending-Borrowing/internal/infrastructure/blockchain"
	"github.com/ethereum/go-ethereum/common"
)

// Config holds all configuration for the application
type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	Blockchain BlockchainConfig
	JWT        JWTConfig
	Server     ServerConfig
}

// AppConfig holds application-wide configuration
type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
	LogLevel    string
}

// DatabaseConfig holds database connection information
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

// BlockchainConfig holds blockchain connection information
type BlockchainConfig struct {
	RpcURL            string
	NetworkName       blockchain.Network
	ChainID           int
	GasLimit          uint64
	GasPrice          int64
	ContractAddresses map[string]common.Address
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	ExpireTime int // In hours
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host string
	Port int
}

// LoadConfig loads all configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load application configuration
	appConfig := AppConfig{
		Name:        GetEnv("APP_NAME", "Lending-Borrowing"),
		Environment: GetEnv("APP_ENV", "development"),
		Debug:       GetEnvBool("APP_DEBUG", true),
		LogLevel:    GetEnv("LOG_LEVEL", "info"),
	}

	// Load database configuration
	dbConfig := DatabaseConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     GetEnvInt("DB_PORT", 5432),
		User:     GetEnv("DB_USER", "postgres"),
		Password: GetEnv("DB_PASSWORD", ""),
		DBName:   GetEnv("DB_NAME", "lending_borrowing"),
		SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		TimeZone: GetEnv("DB_TIMEZONE", "UTC"),
	}

	// Load blockchain configuration
	networkName := blockchain.Network(GetEnv("BLOCKCHAIN_NETWORK", "local"))

	// Parse contract addresses
	contractAddressesStr := GetEnv("CONTRACT_ADDRESSES", "")
	contractAddresses := make(map[string]common.Address)

	if contractAddressesStr != "" {
		pairs := strings.Split(contractAddressesStr, ",")
		for _, pair := range pairs {
			keyVal := strings.Split(pair, "=")
			if len(keyVal) == 2 {
				name := keyVal[0]
				address := keyVal[1]
				if common.IsHexAddress(address) {
					contractAddresses[name] = common.HexToAddress(address)
				} else {
					return nil, fmt.Errorf("invalid contract address for %s: %s", name, address)
				}
			}
		}
	}

	blockchainConfig := BlockchainConfig{
		RpcURL:            GetEnv("BLOCKCHAIN_RPC_URL", "http://localhost:8545"),
		NetworkName:       networkName,
		ChainID:           GetEnvInt("BLOCKCHAIN_CHAIN_ID", 1337),
		GasLimit:          uint64(GetEnvInt("BLOCKCHAIN_GAS_LIMIT", 3000000)),
		GasPrice:          int64(GetEnvInt("BLOCKCHAIN_GAS_PRICE", 20000000000)), // 20 Gwei
		ContractAddresses: contractAddresses,
	}

	// Load JWT configuration
	jwtConfig := JWTConfig{
		Secret:     GetEnv("JWT_SECRET", "your-256-bit-secret"),
		ExpireTime: GetEnvInt("JWT_EXPIRE", 24),
	}

	// Load server configuration
	serverConfig := ServerConfig{
		Host: GetEnv("SERVER_HOST", "localhost"),
		Port: GetEnvInt("SERVER_PORT", 8080),
	}

	return &Config{
		App:        appConfig,
		Database:   dbConfig,
		Blockchain: blockchainConfig,
		JWT:        jwtConfig,
		Server:     serverConfig,
	}, nil
}

// GetDatabaseDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode, c.TimeZone,
	)
}
