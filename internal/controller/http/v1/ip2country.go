package v1

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ransoor2/ip2country/pkg/logger"
)

type IP2CountryService interface {
	IP2CountryNCity(context.Context, string) (string, string, error)
}

type findCountryResponse struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

type ip2CountryNCityRoutes struct {
	ip2Country IP2CountryService
	logger     logger.Interface
}

func newIPToCountryRoutes(routerGroup *gin.RouterGroup, t IP2CountryService, l logger.Interface) {
	ip := &ip2CountryNCityRoutes{t, l}

	routerGroup.GET("/find-country", ip.findCountry)
}

// @Summary     Find Country
// @Description Find country by IP
// @ID          find-country
// @Accept      json
// @Produce     json
// @Param       ip query string true "IP address"
// @Success     200 {object} findCountryResponse
// @Failure     400 {object} response
// @Failure     429 {object} response
// @Failure     500 {object} response
// @Router      /find-country [get]
func (r *ip2CountryNCityRoutes) findCountry(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		r.logger.Error("ip query parameter is required", "http - v1 - findCountry")
		errorResponse(c, http.StatusBadRequest, "ip query parameter is required")
		return
	}

	if net.ParseIP(ip) == nil {
		r.logger.Error("invalid IP address format", "http - v1 - findCountry")
		errorResponse(c, http.StatusBadRequest, "invalid IP address format")
		return
	}

	country, city, err := r.ip2Country.IP2CountryNCity(c.Request.Context(), ip)
	if err != nil {
		r.logger.Error("error finding country", "http - v1 - findCountry", "error", err)
		errorResponse(c, http.StatusInternalServerError, "error finding country")
		return
	}

	if country == "" && city == "" {
		r.logger.Error("country and city not found", "http - v1 - findCountry")
		errorResponse(c, http.StatusNotFound, "country and city not found")
		return
	}

	c.JSON(http.StatusOK, findCountryResponse{
		Country: country,
		City:    city,
	})
}
