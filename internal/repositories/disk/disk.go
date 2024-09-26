package disk

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
)

type Repository struct {
	data map[string]CountryCity
}

type CountryCity struct {
	IP      string `json:"ip"`
	City    string `json:"city"`
	Country string `json:"country"`
}

func New(path string) (Repository, error) {
	repo := Repository{
		data: make(map[string]CountryCity),
	}

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(filePath) == ".json" {
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			var countryCities []CountryCity
			if err := json.Unmarshal(fileData, &countryCities); err != nil {
				return err
			}

			for _, countryCity := range countryCities {
				repo.data[countryCity.IP] = countryCity
			}
		}
		return nil
	})

	if err != nil {
		return Repository{}, err
	}

	return repo, nil
}

func (r Repository) CountryNCityByIP(_ context.Context, ip string) (country, city string, err error) {
	if countryCity, exists := r.data[ip]; exists {
		return countryCity.Country, countryCity.City, nil
	}
	return "", "", nil
}
