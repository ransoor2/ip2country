// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/ransoor2/ip2country/docs"
	"github.com/ransoor2/ip2country/pkg/logger"
)

type RateLimiter interface {
	Allow(ctx context.Context, clientIP string) bool
}

// NewRouter -.
// Swagger spec:
// @title       IP2CountryNCity API
// @description Translating IP 2 Country
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, ip2CountryService IP2CountryService,
	rateLimiter RateLimiter) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	routerGroup := handler.Group("/v1")
	routerGroup.Use(rateLimiterMiddleware(rateLimiter))

	newIPToCountryRoutes(routerGroup, ip2CountryService, l)

}

func rateLimiterMiddleware(rl RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if !rl.Allow(c.Request.Context(), clientIP) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
