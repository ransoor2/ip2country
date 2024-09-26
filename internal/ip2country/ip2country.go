package ip2country

import (
	"context"
	"time"

	"github.com/ransoor2/ip2country/pkg/logger"
)

type Repository interface {
	CountryNCityByIP(context.Context, string) (string, string, error)
}

type Cache interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) (interface{}, bool)
}

type IP2Country struct {
	repo   Repository
	logger logger.Interface
	cache  Cache
}

func New(repo Repository, l logger.Interface, cache Cache) *IP2Country {
	return &IP2Country{
		repo:   repo,
		logger: l,
		cache:  cache,
	}
}

func (c *IP2Country) IP2CountryNCity(ctx context.Context, ip string) (country, city string, err error) {
	// Check cache first
	if cachedValue, found := c.cache.Get(ip); found {
		if result, ok := cachedValue.([2]string); ok {
			return result[0], result[1], nil
		}
	}

	// Fetch from repository if not in cache
	country, city, err = c.repo.CountryNCityByIP(ctx, ip)
	if err != nil {
		c.logger.Error("error finding country", "http - v1 - findCountry", "error", err)
		return "", "", err
	}

	// Store result in cache
	c.cache.Set(ip, [2]string{country, city}, 10*time.Minute)

	return country, city, nil
}
