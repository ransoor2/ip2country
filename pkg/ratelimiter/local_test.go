package ratelimiter

import (
	"testing"
	"time"

	"github.com/ransoor2/ip2country/config"
	"github.com/stretchr/testify/assert"
)

func TestLocalRateLimiter(t *testing.T) {
	cfg := config.RateLimiter{
		MaxRequests:   10,
		UserRequests:  5,
		Interval:      time.Second,
		CleanInterval: time.Second * 10,
		BucketTTL:     time.Second * 5,
	}
	rl := NewLocalRateLimiter(cfg)

	clientIP := "192.168.1.1"

	// Test allowing requests within the limit
	for i := 0; i < 5; i++ {
		assert.True(t, rl.Allow(clientIP), "Request should be allowed")
	}

	// Test exceeding the limit
	assert.False(t, rl.Allow(clientIP), "Request should be denied")

	// Wait for the refill rate duration and test again
	time.Sleep(time.Second)

	for i := 0; i < 5; i++ {
		assert.True(t, rl.Allow(clientIP), "Request should be allowed")
	}
}

func TestGlobalRateLimiter(t *testing.T) {
	cfg := config.RateLimiter{
		MaxRequests:   10,
		UserRequests:  5,
		Interval:      time.Second,
		CleanInterval: time.Second * 10,
		BucketTTL:     time.Second * 5,
	}
	rl := NewLocalRateLimiter(cfg)

	clientIP1 := "192.168.1.1"
	clientIP2 := "192.168.1.2"

	// Test allowing requests within the global limit
	for i := 0; i < 5; i++ {
		assert.True(t, rl.Allow(clientIP1), "Request should be allowed for clientIP1")
		assert.True(t, rl.Allow(clientIP2), "Request should be allowed for clientIP2")
	}

	// Test exceeding the global limit
	assert.False(t, rl.Allow(clientIP1), "Request should be denied for clientIP1")
	assert.False(t, rl.Allow(clientIP2), "Request should be denied for clientIP2")

	// Wait for the refill rate duration and test again
	time.Sleep(time.Second)

	// Check the global bucket capacity
	assert.True(t, rl.Allow(clientIP1), "Request should be allowed for clientIP1 after refill")
	assert.True(t, rl.Allow(clientIP2), "Request should be allowed for clientIP2 after refill")
}

func TestBucketCleanup(t *testing.T) {
	cfg := config.RateLimiter{
		MaxRequests:   10,
		UserRequests:  5,
		Interval:      time.Second,
		CleanInterval: time.Second,
		BucketTTL:     time.Second * 2,
	}
	rl := NewLocalRateLimiter(cfg)

	clientIP := "192.168.1.1"

	// Allow a request to create the bucket
	assert.True(t, rl.Allow(clientIP), "Request should be allowed")

	// Ensure the bucket exists
	rl.mu.Lock()
	_, exists := rl.buckets[clientIP]
	rl.mu.Unlock()
	assert.True(t, exists, "Bucket should exist")

	// Wait for the bucket TTL to expire
	time.Sleep(time.Second * 3)

	// Ensure the bucket has been cleaned up
	rl.mu.Lock()
	_, exists = rl.buckets[clientIP]
	rl.mu.Unlock()
	assert.False(t, exists, "Bucket should be cleaned up")
}
