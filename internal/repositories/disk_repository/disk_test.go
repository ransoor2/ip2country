package disk_repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a sample data.json file
	sampleData := `[
        {"ip": "2.22.233.255", "city": "Sample City", "country": "Sample Country"},
        {"ip": "8.8.8.8", "city": "Mountain View", "country": "United States"},
        {"ip": "1.1.1.1", "city": "Research", "country": "Australia"}
    ]`
	sampleFilePath := filepath.Join(tempDir, "data.json")
	if err := os.WriteFile(sampleFilePath, []byte(sampleData), 0644); err != nil {
		t.Fatalf("Failed to write sample JSON file: %v", err)
	}

	// Initialize the repository
	repo, err := New(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Check if the data was loaded correctly
	if len(repo.data) != 3 {
		t.Errorf("Expected 3 entries in the repository, got %d", len(repo.data))
	}

	testCases := []struct {
		ip      string
		city    string
		country string
		exists  bool
	}{
		{"2.22.233.255", "Sample City", "Sample Country", true},
		{"8.8.8.8", "Mountain View", "United States", true},
		{"1.1.1.1", "Research", "Australia", true},
		{"123.123.123.123", "", "", false}, // Non-existent IP
	}

	for _, tc := range testCases {
		countryCity, exists := repo.data[tc.ip]
		if tc.exists && !exists {
			t.Errorf("Expected entry for IP %s not found", tc.ip)
		}
		if !tc.exists && exists {
			t.Errorf("Unexpected entry found for IP %s", tc.ip)
		}
		if exists && (countryCity.City != tc.city || countryCity.Country != tc.country) {
			t.Errorf("Data mismatch for IP %s: got %+v", tc.ip, countryCity)
		}
	}
}
