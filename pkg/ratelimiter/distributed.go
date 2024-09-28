package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ransoor2/ip2country/config"
	"github.com/ransoor2/ip2country/pkg/logger"
)

const (
	global               = "global"
	rateLimiterKeyPrefix = "rate_limiter:"
)

type DistributedRateLimiter struct {
	log                  logger.Interface
	client               *redis.Client
	globalBucketCapacity int
	bucketCapacity       int
	interval             time.Duration
}

func NewDistributedRateLimiter(cfg config.RateLimiter, l logger.Interface) *DistributedRateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		l.Fatal("Failed to connect to Redis: %v", err)
	}

	return &DistributedRateLimiter{
		log:                  l,
		client:               client,
		globalBucketCapacity: cfg.MaxRequests,
		bucketCapacity:       cfg.UserRequests,
		interval:             cfg.Interval,
	}
}
func (rl *DistributedRateLimiter) Allow(ctx context.Context, clientIP string) bool {
	globalKey := rateLimiterKeyPrefix + global
	clientKey := rateLimiterKeyPrefix + clientIP

	rl.log.Debug("Rate limiting keys: global=%s, client=%s", globalKey, clientKey)

	// Increment global bucket
	globalCurrent, err := rl.client.Incr(ctx, globalKey).Result()
	if err != nil {
		return false
	}

	if globalCurrent == 1 {
		_, err = rl.client.Expire(ctx, globalKey, rl.interval).Result()
		if err != nil {
			return false
		}
	}

	// Increment client bucket
	clientCurrent, err := rl.client.Incr(ctx, clientKey).Result()
	if err != nil {
		return false
	}

	if clientCurrent == 1 {
		_, err = rl.client.Expire(ctx, clientKey, rl.interval).Result()
		if err != nil {
			return false
		}
	}

	// Check if both global and client limits are within the allowed capacity
	return globalCurrent <= int64(rl.globalBucketCapacity) && clientCurrent <= int64(rl.bucketCapacity)
}
