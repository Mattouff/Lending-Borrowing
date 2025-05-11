package valkey

import (
    "fmt"
    "sync"
    "time"
)

// tokenCacheItem represents a cached token validation result
type tokenCacheItem struct {
    valid      bool
    expiration time.Time
}

// TokenCache provides a simple in-memory cache for token validation
type TokenCache struct {
    items map[string]tokenCacheItem
    mu    sync.RWMutex
    ttl   time.Duration
}

// NewTokenCache creates a new token cache with the specified TTL
func NewTokenCache(ttl time.Duration) *TokenCache {
    cache := &TokenCache{
        items: make(map[string]tokenCacheItem),
        ttl:   ttl,
    }
    
    // Start a background goroutine to clean expired items
    go cache.cleanupLoop()
    
    return cache
}

// Get retrieves a token validation result from the cache
func (c *TokenCache) Get(userID uint, tokenID string) (bool, bool) {
    key := formatCacheKey(userID, tokenID)
    
    c.mu.RLock()
    item, found := c.items[key]
    c.mu.RUnlock()
    
    if !found {
        return false, false
    }
    
    // Check if item has expired
    if time.Now().After(item.expiration) {
        c.mu.Lock()
        delete(c.items, key)
        c.mu.Unlock()
        return false, false
    }
    
    return item.valid, true
}

// Set stores a token validation result in the cache
func (c *TokenCache) Set(userID uint, tokenID string, valid bool) {
    key := formatCacheKey(userID, tokenID)
    
    c.mu.Lock()
    c.items[key] = tokenCacheItem{
        valid:      valid,
        expiration: time.Now().Add(c.ttl),
    }
    c.mu.Unlock()
}

// Invalidate removes a specific token from the cache
func (c *TokenCache) Invalidate(userID uint, tokenID string) {
    key := formatCacheKey(userID, tokenID)
    
    c.mu.Lock()
    delete(c.items, key)
    c.mu.Unlock()
}

// InvalidateUser removes all tokens for a user from the cache
func (c *TokenCache) InvalidateUser(userID uint) {
    prefix := formatUserPrefix(userID)
    
    c.mu.Lock()
    for key := range c.items {
        if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
            delete(c.items, key)
        }
    }
    c.mu.Unlock()
}

// cleanupLoop periodically removes expired items from the cache
func (c *TokenCache) cleanupLoop() {
    ticker := time.NewTicker(c.ttl / 2)
    defer ticker.Stop()
    
    for range ticker.C {
        c.cleanup()
    }
}

// cleanup removes expired items from the cache
func (c *TokenCache) cleanup() {
    now := time.Now()
    
    c.mu.Lock()
    for key, item := range c.items {
        if now.After(item.expiration) {
            delete(c.items, key)
        }
    }
    c.mu.Unlock()
}

// formatCacheKey formats a cache key for a token
func formatCacheKey(userID uint, tokenID string) string {
    return fmt.Sprintf("%d:%s", userID, tokenID)
}

// formatUserPrefix formats a prefix for all tokens of a user
func formatUserPrefix(userID uint) string {
    return fmt.Sprintf("%d:", userID)
}