package test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/ransoor2/ip2country/config"
	v1 "github.com/ransoor2/ip2country/internal/controller/http/v1"
	"github.com/ransoor2/ip2country/internal/ip2country"
	"github.com/ransoor2/ip2country/internal/repositories/disk"
	"github.com/ransoor2/ip2country/pkg/cache"
	"github.com/ransoor2/ip2country/pkg/httpserver"
	"github.com/ransoor2/ip2country/pkg/logger"
)

type APITestSuite struct {
	suite.Suite
	client *http.Client
	server *httpserver.Server
}

func (s *APITestSuite) SetupSuite() {
	s.client = &http.Client{}

	// Configuration
	os.Setenv("DISK_REPOSITORY_RELATIVE_PATH", "data.json")
	cfg, err := config.NewConfig("../config/config.yml")
	assert.NoError(s.T(), err)

	l := logger.New(cfg.Log.Level)
	// Cache
	cacheInst, err := cache.New(cfg.Cache.Size)
	assert.NoError(s.T(), err)

	// Repository
	repo, err := disk.New(cfg.DiskRepository.RelativePath)
	assert.NoError(s.T(), err)

	// Use case
	ip2CountryService := ip2country.New(repo, l, cacheInst)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, ip2CountryService)

	// Run

	go func() { s.server = httpserver.New(handler, httpserver.Port(cfg.HTTP.Port)) }()
	// Wait for listener to start
	assert.Eventually(s.T(),
		func() bool {
			res, err := s.client.Get("http://localhost:8080/healthz")
			if res != nil {
				defer res.Body.Close()
			}
			return err == nil
		},
		50*time.Millisecond,
		10*time.Millisecond,
	)
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) TearDownSuite() {
	assert.NoError(s.T(), s.server.Shutdown())
}

func (s *APITestSuite) TestGetCountryNCityByIPHappy() {
	country, city, statusCode, err := s.getCountryNCityByIP(`2.22.233.255`)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, statusCode)
	assert.Equal(s.T(), "Sample Country", country)
	assert.Equal(s.T(), "Sample City", city)
}

func (s *APITestSuite) TestGetCountryNCityByIPNotFound() {
	country, city, statusCode, err := s.getCountryNCityByIP(`1.2.3.4`)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusNotFound, statusCode)
	assert.Empty(s.T(), country)
	assert.Empty(s.T(), city)
}

func (s *APITestSuite) TestGetCountryNCityByIPInvalidIP() {
	country, city, statusCode, err := s.getCountryNCityByIP(`1.2.3.4.5`)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusBadRequest, statusCode)
	assert.Empty(s.T(), country)
	assert.Empty(s.T(), city)
}
