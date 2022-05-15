package world

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/alejoacosta74/allien_invasion/cities"
)

type IWorld interface {
	GetRandomCity() *cities.City
	GetCityByName(name string) *cities.City
	DestroyCity(city *cities.City) error
	GetCities() map[string]*cities.City
}

type World struct {
	cityMap map[string]*cities.City
}

func NewWorldFromFile(filename string) (*World, error) {
	w, err := cities.LoadCityMapFromFile(filename)
	return &World{cityMap: w}, err
}

func (w *World) GetCities() map[string]*cities.City {
	return w.cityMap
}

func (w *World) GetRandomCity() *cities.City {
	l := len(w.cityMap)
	if l == 0 {
		return nil
	}
	cities := make([]*cities.City, 0)
	for _, c := range w.cityMap {
		cities = append(cities, c)
	}
	rand.Seed(time.Now().UnixNano())
	return cities[rand.Intn(l)]
}

func (w *World) GetCityByName(name string) *cities.City {
	for _, c := range w.cityMap {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (w *World) DestroyCity(city *cities.City) error {
	neighbors := w.cityMap[city.Name].GetNeighborsMap()
	delete(w.cityMap, city.Name)
	for _, n := range neighbors {
		c := w.cityMap[n]
		for d, n := range c.GetNeighborsMap() {
			if n == city.Name {
				delete(c.GetNeighborsMap(), d)
			}
		}

	}
	return nil
}

func (w *World) PrintWorld() {
	println("\nPrinting world:")
	println("----------------------------------------")
	println("Number of cities: ", len(w.cityMap))
	for _, c := range w.cityMap {
		println("\tCity: ", c.Name)
		for d, n := range c.GetNeighborsMap() {
			fmt.Printf("\t\t direction -> %s : neighbor: %s\n", d, n)
		}
	}
	println("----------------------------------------\n")
}
