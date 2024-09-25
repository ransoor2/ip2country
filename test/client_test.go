package test

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	baseURI = "http://localhost:8080/v1/find-country"
)

type findCountryResponse struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

func (s *APITestSuite) getCountryNCityByIP(ip string) (string, string, int, error) {
	uri := fmt.Sprintf("%s?ip=%s", baseURI, ip)
	response, err := s.client.Get(uri)
	statusCode := response.StatusCode
	if err != nil {
		return "", "", statusCode, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			s.T().Errorf("error closing response body: %v", err)
		}
	}(response.Body)

	var result findCountryResponse
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", "", statusCode, err
	}

	return result.Country, result.City, statusCode, nil
}
