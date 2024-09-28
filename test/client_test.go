package test

import (
	"encoding/json"
	"fmt"

	"github.com/stretchr/testify/assert"
)

const (
	baseURI = "http://localhost:8080/v1/find-country"
)

type findCountryResponse struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

func (s *APITestSuite) getCountryNCityByIP(ip string) (country, city string, statusCode int, err error) {
	uri := fmt.Sprintf("%s?ip=%s", baseURI, ip)
	response, err := s.client.Get(uri)
	assert.NoError(s.T(), err)
	defer func() {
		closeErr := response.Body.Close()
		assert.NoError(s.T(), closeErr)
	}()

	statusCode = response.StatusCode
	var result findCountryResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", "", statusCode, err
	}

	return result.Country, result.City, statusCode, nil
}
