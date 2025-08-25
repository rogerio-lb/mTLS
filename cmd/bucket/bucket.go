package main

import (
	"fmt"
	"sync"
	"time"
)

// TokenBucket represents a token bucket for rate limiting
type TokenBucket struct {
	capacity   int           // Maximum number of tokens
	tokens     int           // Current number of tokens
	refillRate time.Duration // Time between token refills
	lastRefill time.Time     // Last time tokens were added
	mu         sync.Mutex    // Mutex for thread safety
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity, // Start with full capacity
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// refill adds tokens based on elapsed time
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	// Calculate how many tokens to add based on elapsed time
	tokensToAdd := int(elapsed / tb.refillRate)

	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}
}

// CanConsume checks if a token can be consumed without actually consuming it
func (tb *TokenBucket) CanConsume() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	return tb.tokens > 0
}

// TryConsume attempts to consume a token, returns true if successful
func (tb *TokenBucket) TryConsume() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// GetTokens returns the current number of tokens (for debugging)
func (tb *TokenBucket) GetTokens() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	return tb.tokens
}

// ValidateMessage validates if a message can be processed
func ValidateMessage(tb *TokenBucket, messageID string) bool {
	if tb.TryConsume() {
		fmt.Printf("Message %s: ALLOWED (tokens remaining: %d)\n", messageID, tb.GetTokens())
		return true
	}

	fmt.Printf("Message %s: REJECTED (no tokens available)\n", messageID)
	return false
}

func main() {
	// Create a token bucket with capacity of 5 tokens, refilling 1 token per second
	bucket := NewTokenBucket(5, time.Second)

	fmt.Println("Token Bucket Rate Limiter Demo")
	fmt.Printf("Capacity: %d tokens, Refill rate: 1 token/second\n\n", bucket.capacity)

	// Simulate processing messages
	messages := []string{"msg1", "msg2", "msg3", "msg4", "msg5", "msg6", "msg7"}

	// Process messages rapidly (should exhaust tokens)
	fmt.Println("Processing messages rapidly:")
	for _, msg := range messages {
		ValidateMessage(bucket, msg)
	}

	fmt.Printf("\nWaiting 3 seconds for token refill...\n")
	time.Sleep(3 * time.Second)

	// Try processing more messages after waiting
	fmt.Println("Processing messages after waiting:")
	for i := 8; i <= 10; i++ {
		ValidateMessage(bucket, fmt.Sprintf("msg%d", i))
	}
}
