// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/ransoor2/ip2country/config"
	v1 "github.com/ransoor2/ip2country/internal/controller/http/v1"
	"github.com/ransoor2/ip2country/internal/ip2country"
	"github.com/ransoor2/ip2country/internal/repositories/disk"
	"github.com/ransoor2/ip2country/internal/repositories/mongo"
	"github.com/ransoor2/ip2country/pkg/cache"
	"github.com/ransoor2/ip2country/pkg/httpserver"
	"github.com/ransoor2/ip2country/pkg/logger"
	"github.com/ransoor2/ip2country/pkg/ratelimiter"
)

// Constants for repository types
const (
	RepoTypeMongo = "mongo"
	RepoTypeDisk  = "disk"
)

// Constants for rate limiter types
const (
	RateLimiterTypeLocal       = "local"
	RateLimiterTypeDistributed = "distributed"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Cache
	cacheInst, err := cache.New(cfg.Cache.Size)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - cacheInst.New: %w", err))
	}

	// Repository
	repo, err := initializeRepository(cfg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - initializeRepository: %w", err))
	}

	// RateLimiter
	rateLimiter, err := getRateLimiter(cfg, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - getRateLimiter: %w", err))
	}

	// Use case
	ip2CountryService := ip2country.New(repo, l, cacheInst)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, ip2CountryService, rateLimiter)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

func initializeRepository(cfg *config.Config) (ip2country.Repository, error) {
	switch cfg.Repository.Type {
	case RepoTypeMongo:
		return mongo.New(cfg.MongoRepository.URI, cfg.MongoRepository.DB, cfg.MongoRepository.Collection)
	case RepoTypeDisk:
		return disk.New(cfg.DiskRepository.RelativePath)
	default:
		return nil, fmt.Errorf("app - initializeRepository - unknown repository type: %s", cfg.Repository.Type)
	}
}

func getRateLimiter(cfg *config.Config, l logger.Interface) (v1.RateLimiter, error) {
	switch cfg.RateLimiter.Type {
	case RateLimiterTypeLocal:
		return ratelimiter.NewLocalRateLimiter(cfg.RateLimiter, l), nil
	case RateLimiterTypeDistributed:
		return ratelimiter.NewDistributedRateLimiter(cfg.RateLimiter, l), nil
	default:
		return nil, fmt.Errorf("unknown rate limiter type: %s", cfg.RateLimiter.Type)
	}
}
