package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// Create a temporary YAML configuration file
	yamlContent := `
app:
  name: "TestApp"
  version: "1.0.0"
http:
  port: "8080"
logger:
  log_level: "debug"
cache:
  size: 100
repository:
  type: "mongo"
mongoRepository:
  uri: "mongodb://localhost:27017"
  db: "testdb"
  collection: "testcollection"
rateLimiter:
  type: "local"
  maxRequests: 100
  userRequests: 5
  interval: 1s
  redisAddr: 'localhost:6379'
`
	tmpFile, err := os.CreateTemp("", "config-*.yml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)

	// Set environment variables
	os.Setenv("APP_NAME", "EnvApp")
	os.Setenv("APP_VERSION", "2.0.0")
	os.Setenv("HTTP_PORT", "9090")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("CACHE_SIZE", "200")
	os.Setenv("REPOSITORY_TYPE", "mongo")
	os.Setenv("MONGO_REPOSITORY_URI", "mongodb://envhost:27017")
	os.Setenv("MONGO_REPOSITORY_DB", "envdb")
	os.Setenv("MONGO_REPOSITORY_COLLECTION", "envcollection")
	os.Setenv("RATE_LIMITER_TYPE", "distributed")
	os.Setenv("RATE_LIMITER_MAX_REQUESTS", "150")
	os.Setenv("RATE_LIMITER_USER_REQUESTS", "10")
	os.Setenv("RATE_LIMITER_INTERVAL", "2s")
	os.Setenv("RATE_LIMITER_REDIS_ADDR", "localhost:6379")
	defer os.Clearenv()

	// Load configuration
	cfg, err := NewConfig(tmpFile.Name())
	assert.NoError(t, err)

	// Validate configuration values
	assert.Equal(t, "EnvApp", cfg.App.Name)
	assert.Equal(t, "2.0.0", cfg.App.Version)
	assert.Equal(t, "9090", cfg.HTTP.Port)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, 200, cfg.Cache.Size)
	assert.Equal(t, "mongo", cfg.Repository.Type)
	assert.Equal(t, "mongodb://envhost:27017", cfg.MongoRepository.URI)
	assert.Equal(t, "envdb", cfg.MongoRepository.DB)
	assert.Equal(t, "envcollection", cfg.MongoRepository.Collection)
	assert.Equal(t, "distributed", cfg.RateLimiter.Type)
	assert.Equal(t, 150, cfg.RateLimiter.MaxRequests)
	assert.Equal(t, 10, cfg.RateLimiter.UserRequests)
	assert.Equal(t, 2*time.Second, cfg.RateLimiter.Interval)
	assert.Equal(t, "localhost:6379", cfg.RateLimiter.RedisAddr)
}

func TestInvalidRepositoryType(t *testing.T) {
	// Create a temporary YAML configuration file
	yamlContent := `
app:
  name: "TestApp"
  version: "1.0.0"
http:
  port: "8080"
logger:
  log_level: "debug"
cache:
  size: 100
repository:
  type: "mongo"
mongoRepository:
  uri: "mongodb://localhost:27017"
  db: "testdb"
  collection: "testcollection"
rateLimiter:
  type: "local"
  maxRequests: 100
  userRequests: 5
  interval: 1s
  redisAddr: 'localhost:6379'
`
	tmpFile, err := os.CreateTemp("", "config-*.yml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)

	// Set environment variables with an invalid repository type
	os.Setenv("REPOSITORY_TYPE", "invalid_type")
	defer os.Clearenv()

	// Load configuration
	_, err = NewConfig(tmpFile.Name())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation error")
}
