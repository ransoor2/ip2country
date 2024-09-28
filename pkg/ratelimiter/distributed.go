package ratelimiter

type DistributedRateLimiter struct {
}

func NewDistributedRateLimiter() *DistributedRateLimiter {
	return &DistributedRateLimiter{}
}

func (rl *DistributedRateLimiter) Allow(_ string) bool {
	// Implement distributed rate limiting logic here
	return true
}
