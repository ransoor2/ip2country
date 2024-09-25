// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/ransoor2/ip2country/internal/ip2country"
	"github.com/ransoor2/ip2country/internal/repositories/disk_repository"
	"github.com/ransoor2/ip2country/internal/repositories/mongo_repository"
	"github.com/ransoor2/ip2country/pkg/cache"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/ransoor2/ip2country/config"
	v1 "github.com/ransoor2/ip2country/internal/controller/http/v1"
	"github.com/ransoor2/ip2country/pkg/httpserver"
	"github.com/ransoor2/ip2country/pkg/logger"
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

	// Use case
	ip2CountryService := ip2country.New(repo, l, cacheInst)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, ip2CountryService)
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
	case "mongo":
		return mongo_repository.New(cfg.MongoRepository.URI, cfg.MongoRepository.DB, cfg.MongoRepository.Collection)
	case "disk":
		return disk_repository.New(cfg.DiskRepository.RelativePath)
	default:
		return nil, fmt.Errorf("app - initializeRepository - unknown repository type: %s", cfg.Repository.Type)
	}
}
