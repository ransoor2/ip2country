package ratelimiter

import (
	"sync"
	"time"

	"github.com/ransoor2/ip2country/config"
)

type LocalRateLimiter struct {
	buckets              map[string]*Bucket
	globalBucket         *Bucket
	mu                   sync.Mutex
	bucketCapacity       int
	globalBucketCapacity int
	refillRate           time.Duration
	cleanupInterval      time.Duration
	bucketTTL            time.Duration
}

type Bucket struct {
	capacity  int
	remaining int
	lastCheck time.Time
}

func NewLocalRateLimiter(cfg config.RateLimiter) *LocalRateLimiter {
	rl := &LocalRateLimiter{
		buckets:              make(map[string]*Bucket),
		globalBucket:         &Bucket{capacity: cfg.MaxRequests, remaining: cfg.MaxRequests, lastCheck: time.Now()},
		bucketCapacity:       cfg.UserRequests,
		globalBucketCapacity: cfg.MaxRequests,
		refillRate:           cfg.Interval,
		cleanupInterval:      cfg.CleanInterval,
		bucketTTL:            cfg.BucketTTL,
	}

	go rl.cleanupBuckets()

	return rl
}

func (rl *LocalRateLimiter) Allow(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Check if the clientIP bucket exists, if not create a new one
	bucket, exists := rl.buckets[clientIP]
	if !exists {
		bucket = &Bucket{capacity: rl.bucketCapacity, remaining: rl.bucketCapacity, lastCheck: time.Now()}
		rl.buckets[clientIP] = bucket
	}

	now := time.Now()

	// Refill clientIP bucket
	elapsed := now.Sub(bucket.lastCheck)
	if elapsed >= rl.refillRate {
		bucket.remaining = bucket.capacity
		bucket.lastCheck = now
	}

	// Refill global bucket
	globalElapsed := now.Sub(rl.globalBucket.lastCheck)
	if globalElapsed >= rl.refillRate {
		rl.globalBucket.remaining = rl.globalBucket.capacity
		rl.globalBucket.lastCheck = now
	}

	// Check if there are enough tokens in both the clientIP and global bucket
	if bucket.remaining > 0 && rl.globalBucket.remaining > 0 {
		bucket.remaining--
		rl.globalBucket.remaining--
		return true
	}

	return false
}

func (rl *LocalRateLimiter) cleanupBuckets() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, bucket := range rl.buckets {
			if now.Sub(bucket.lastCheck) > rl.bucketTTL {
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}
