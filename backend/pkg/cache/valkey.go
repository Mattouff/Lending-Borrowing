package valkey

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Mattouff/Lending-Borrowing/internal/config"
	"github.com/valkey-io/valkey-go"
)

// Client wraps Valkey client functionality
type Client struct {
	client valkey.Client
	cache  *TokenCache
}

// NewClient creates a new Valkey client
func NewClient(cfg *config.Config) (*Client, error) {
	// Create client options
	options := valkey.ClientOption{
		InitAddress: []string{fmt.Sprintf("%s:%s", cfg.Valkey.Host, cfg.Valkey.Port)},
	}

	// Add password if provided
	if cfg.Valkey.Password != "" {
		options.Password = cfg.Valkey.Password
	}

	// Create the client
	client, err := valkey.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create valkey client: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Do(ctx, client.B().Ping().Build()).Error()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to valkey: %w", err)
	}

	// If DB is specified, select it
	if cfg.Valkey.DB > 0 {
		err = client.Do(ctx, client.B().Select().Index(int64(cfg.Valkey.DB)).Build()).Error()
		if err != nil {
			return nil, fmt.Errorf("failed to select database: %w", err)
		}
	}

	// Create in-memory cache with 10 second TTL
	cache := NewTokenCache(10 * time.Second)

	log.Println("Valkey client connected successfully")
	return &Client{
		client: client,
		cache:  cache,
	}, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	c.client.Close()
	return nil
}

// StoreValidToken adds a token to a user's valid tokens set
func (c *Client) StoreValidToken(ctx context.Context, userID uint, tokenID string, expiration time.Duration) error {
	key := formatValidTokenKey(userID)

	// Add token to set
	err := c.client.Do(ctx, c.client.B().Sadd().Key(key).Member(tokenID).Build()).Error()
	if err != nil {
		return err
	}

	// Set expiration on the key
	err = c.client.Do(ctx, c.client.B().Expire().Key(key).Seconds(int64(expiration.Seconds())).Build()).Error()
	if err != nil {
		return err
	}

	// Update cache
	c.cache.Set(userID, tokenID, true)

	return nil
}

// IsValidToken checks if a token is valid for a user
func (c *Client) IsValidToken(ctx context.Context, userID uint, tokenID string) (bool, error) {
	// Check cache first
	if valid, found := c.cache.Get(userID, tokenID); found {
		return valid, nil
	}

	// Not in cache, check Valkey
	key := formatValidTokenKey(userID)
	result, err := c.client.Do(ctx, c.client.B().Sismember().Key(key).Member(tokenID).Build()).AsInt64()

	// Update cache if no error
	if err == nil {
		valid := result == 1
		c.cache.Set(userID, tokenID, valid)
		return valid, nil
	}

	return false, err
}

// InvalidateAllUserTokens removes all valid tokens for a user
func (c *Client) InvalidateAllUserTokens(ctx context.Context, userID uint) error {
	key := formatValidTokenKey(userID)

	// Clear from Valkey
	err := c.client.Do(ctx, c.client.B().Del().Key(key).Build()).Error()

	// Clear from cache (even if Valkey operation failed)
	c.cache.InvalidateUser(userID)

	return err
}

// InvalidateToken removes a specific token from valid tokens
func (c *Client) InvalidateToken(ctx context.Context, userID uint, tokenID string) error {
	key := formatValidTokenKey(userID)

	// Remove from Valkey
	err := c.client.Do(ctx, c.client.B().Srem().Key(key).Member(tokenID).Build()).Error()

	// Remove from cache
	c.cache.Invalidate(userID, tokenID)

	return err
}

// Helper function to format Valkey key for valid tokens
func formatValidTokenKey(userID uint) string {
	return fmt.Sprintf("valid_token:%d", userID)
}

// Ping checks if the Valkey server is reachable
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Do(ctx, c.client.B().Ping().Build()).Error()
}
