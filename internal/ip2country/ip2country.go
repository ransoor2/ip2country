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

func New(repo Repository, logger logger.Interface, cache Cache) *IP2Country {
	return &IP2Country{
		repo:   repo,
		logger: logger,
		cache:  cache,
	}
}

func (I *IP2Country) IP2CountryNCity(ctx context.Context, ip string) (string, string, error) {
	// Check cache first
	if cachedValue, found := I.cache.Get(ip); found {
		if result, ok := cachedValue.([2]string); ok {
			return result[0], result[1], nil
		}
	}

	// Fetch from repository if not in cache
	country, city, err := I.repo.CountryNCityByIP(ctx, ip)
	if err != nil {
		I.logger.Error("error finding country", "http - v1 - findCountry", "error", err)
		return "", "", err
	}

	// Store result in cache
	I.cache.Set(ip, [2]string{country, city}, 10*time.Minute)

	return country, city, nil
}
