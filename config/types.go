package config

import (
	"github.com/pkg/errors"
)

type CityName string

type Direction string

const (
	// "Nil" direction to use in case of errors.
	NilDirection Direction = ""
	North                  = "north"
	West                   = "west"
	East                   = "east"
	South                  = "south"
)

const (
	// "Nil" city name to return in case of errors.
	NilCityName CityName = ""
)

// Create a Direction object from a string.
func DirectionFromString(s string) Direction {
	dir := Direction(s)
	switch dir {
	case North, East, South, West:
		return dir
	}
	return NilDirection
}

// Create a City from a string.
// Currently, there are no restrictions on city names.
func CityNameFromString(s string) CityName {
	return CityName(s)
}

// Internal representation of the map.
type Map struct {
	Cities map[CityName]City
}

// A single city within the map.
type City struct {
	// Roads out of the city.
	Roads map[Direction]CityName
}

// Create a new, empty map.
func NewMap() *Map {
	return &Map{make(map[CityName]City)}
}

// Add a city to a map, will return an error if the city already exists.
func (m *Map) AddCity(cityName CityName, city City) error {
	_, ok := m.Cities[cityName]
	if ok {
		// City already exists, cannot have the same name.
		return errors.Errorf("Duplicate city name: %s", cityName)
	}
	m.Cities[cityName] = city
	return nil
}
