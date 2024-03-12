package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"
)

// City represents the information about a city
type City struct {
	Code         string
	ProvinceCode string
	CountryCode  string
	Name         string
	ProvinceName string
	CountryName  string
}

// PermissionStrategy defines the strategy for checking permissions
type PermissionStrategy interface {
	CheckPermission(cityCode string) bool
}

// IncludeExcludePermission implements PermissionStrategy
type IncludeExcludePermission struct {
	Include []string
	Exclude []string
}

// Distributor represents a distributor with permissions
type Distributor struct {
	Name        string
	Permissions PermissionStrategy
}

func main() {
	cities, err := LoadCities("cities.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	distributors := []Distributor{
		{
			Name: "DISTRIBUTOR1",
			Permissions: &IncludeExcludePermission{
				Include: []string{"IN", "US"},
				Exclude: []string{"KARNATAKA-IN", "CHENNAI-TAMILNADU-IN"},
			},
		},
		{
			Name: "DISTRIBUTOR2",
			Permissions: &IncludeExcludePermission{
				Include: []string{"IN"},
				Exclude: []string{"TAMILNADU-IN"},
			},
		},
		{
			Name: "DISTRIBUTOR3",
			Permissions: &IncludeExcludePermission{
				Include: []string{"HUBLI-KARNATAKA-IN"},
			},
		},
	}

	PrintPermissions(cities, distributors)
}

// CheckPermission checks if a distributor has permission for a given city
func (p *IncludeExcludePermission) CheckPermission(cityCode string) bool {
	for _, excluded := range p.Exclude {
		if strings.HasPrefix(cityCode, excluded) {
			return false
		}
	}
	for _, included := range p.Include {
		if strings.HasPrefix(cityCode, included) {
			return true
		}
	}
	return false
}

// LoadCities loads city data from the provided CSV file concurrently
func LoadCities(csvFile string) ([]City, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	var cities []City
	var wg sync.WaitGroup
	var mutex sync.Mutex
	errCh := make(chan error)

	for {
		line, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			errCh <- fmt.Errorf("failed to read CSV line: %w", err)
			return nil, <-errCh
		}
		wg.Add(1)
		go func(line []string) {
			defer wg.Done()
			city := City{
				Code:         line[0],
				ProvinceCode: line[1],
				CountryCode:  line[2],
				Name:         line[3],
				ProvinceName: line[4],
				CountryName:  line[5],
			}
			mutex.Lock()
			cities = append(cities, city)
			mutex.Unlock()
		}(line)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	return cities, nil
}

// PrintPermissions prints permissions for all distributors and cities concurrently
func PrintPermissions(cities []City, distributors []Distributor) {
	var wg sync.WaitGroup
	results := make(chan string)

	for _, distributor := range distributors {
		wg.Add(1)
		go func(distributor Distributor) {
			defer wg.Done()
			var result strings.Builder
			result.WriteString(fmt.Sprintf("%s Permissions:\n", distributor.Name))
			for _, city := range cities {
				permission := distributor.Permissions.CheckPermission(city.Code)
				result.WriteString(fmt.Sprintf("%s, %s, %s: %t\n", city.Name, city.ProvinceName, city.CountryName, permission))
			}
			results <- result.String()
		}(distributor)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println(res)
	}
}
