package cities

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
)

// Reads the cities and directions from a file and returns a "world" map
func LoadCityMapFromFile(filename string) (map[string]*City, error) {
	var cityMap = make(map[string]*City)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	buf := bytes.Buffer{}
	for scanner.Scan() {
		line := scanner.Text()
		err := parseCityToMap(cityMap, line)
		if err != nil {
			return nil, fmt.Errorf("error parsing city to map: %s", err)
		}
		fmt.Fprintln(&buf, line) // save for later use when parsin neighbors
	}
	scanner2 := bufio.NewScanner(&buf)
	for scanner2.Scan() {
		line := scanner2.Text()
		err := parseNeighborToMap(cityMap, line)
		if err != nil {
			return nil, fmt.Errorf("error parsing neighbors: %s", err)
		}
	}

	return cityMap, nil

}

// LoadCitiesFS loads cities from a filesystem (mainly used for testing)
func LoadCitiesFS(filesystem fs.FS) ([]*City, error) {
	var cityList []*City

	dir, err := fs.ReadDir(filesystem, ".")
	if err != nil {
		return nil, err
	}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		cities, err := loadcities(filesystem, f.Name())
		if err != nil {
			return nil, err
		}
		if len(cities) > 0 {
			cityList = append(cityList, cities...)
		}
	}
	return cityList, nil
}

func loadcities(filesystem fs.FS, filename string) ([]*City, error) {
	var cities []*City
	file, err := filesystem.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		newCity, err := parseCity(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("error parsing city: %s", err)
		}
		if newCity != nil {
			cities = append(cities, newCity)
		}
	}
	return cities, nil
}
