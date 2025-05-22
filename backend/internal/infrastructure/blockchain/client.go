package blockchain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

// Network represents supported blockchain networks
type Network string

const (
	// Ethereum mainnet
	Mainnet Network = "mainnet"
	// Ethereum testnet (Hoodi)
	Hoodi Network = "hoodi"
	// Ethereum testnet (Sepolia)
	Sepolia Network = "sepolia"
	// Local development network
	Local Network = "local"
)

// ChainConfig contains network configuration details
type ChainConfig struct {
	ChainID       *big.Int
	NetworkName   Network
	RpcURL        string
	BlockExplorer string
	Contracts     map[string]common.Address
}

// EthClient manages Ethereum client connections and provides access to contract wrappers
type EthClient struct {
	client      *ethclient.Client
	config      *ChainConfig
	initialized bool
	mu          sync.Mutex
}

var (
	// Single instance of EthClient (singleton)
	instance *EthClient
	once     sync.Once
)

// GetInstance returns the singleton instance of EthClient
func GetInstance() *EthClient {
	once.Do(func() {
		instance = &EthClient{
			initialized: false,
			config:      nil,
			client:      nil,
		}
	})
	return instance
}

// Initialize sets up the Ethereum client with the provided RPC URL
func (ec *EthClient) Initialize(rpcURL string, networkName Network, contracts map[string]common.Address) error {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if ec.initialized {
		return errors.New("ethereum client is already initialized")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	// Get chain ID to confirm connection
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Set block explorer URL based on network
	var blockExplorer string
	switch networkName {
	case Mainnet:
		blockExplorer = "https://etherscan.io"
		if !isMainnetChainID(chainID) {
			return fmt.Errorf("connected to wrong network: expected mainnet, got chain ID %s", chainID.String())
		}
	case Hoodi:
		blockExplorer = "https://hoodi.etherscan.io/"
		if chainID.Cmp(big.NewInt(560048)) != 0 {
			return fmt.Errorf("connected to wrong network: expected Hoodi (chain ID 560048), got %s", chainID.String())
		}
	case Sepolia:
		blockExplorer = "https://sepolia.etherscan.io"
		if chainID.Cmp(big.NewInt(11155111)) != 0 {
			return fmt.Errorf("connected to wrong network: expected Sepolia (chain ID 11155111), got %s", chainID.String())
		}
	case Local:
		blockExplorer = ""
		// Allow any chain ID for local development
	default:
		return fmt.Errorf("unsupported network: %s", networkName)
	}

	ec.client = client
	ec.config = &ChainConfig{
		ChainID:       chainID,
		NetworkName:   networkName,
		RpcURL:        rpcURL,
		BlockExplorer: blockExplorer,
		Contracts:     contracts,
	}
	ec.initialized = true

	return nil
}

// isMainnetChainID checks if a chain ID corresponds to Ethereum mainnet (1) or equivalent L2s
func isMainnetChainID(id *big.Int) bool {
	return id.Cmp(params.MainnetChainConfig.ChainID) == 0
}

// GetClient returns the ethclient.Client instance
func (ec *EthClient) GetClient() (*ethclient.Client, error) {
	if !ec.initialized {
		return nil, errors.New("ethereum client not initialized")
	}
	return ec.client, nil
}

// GetConfig returns the current chain configuration
func (ec *EthClient) GetConfig() (*ChainConfig, error) {
	if !ec.initialized {
		return nil, errors.New("ethereum client not initialized")
	}
	return ec.config, nil
}

// GetContractAddress returns the address for a named contract
func (ec *EthClient) GetContractAddress(name string) (common.Address, error) {
	if !ec.initialized {
		return common.Address{}, errors.New("ethereum client not initialized")
	}

	address, exists := ec.config.Contracts[name]
	if !exists {
		return common.Address{}, fmt.Errorf("contract not found: %s", name)
	}

	return address, nil
}

// Close closes the client connection
func (ec *EthClient) Close() {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if ec.client != nil {
		ec.client.Close()
		ec.initialized = false
	}
}
