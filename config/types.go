package config

import (
	"fmt"
	"io"
	"strings"

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
	destroyed map[CityName]bool
	cities    map[CityName]City
}

// Create a new, empty map.
func NewMap() *Map {
	return &Map{make(map[CityName]bool), make(map[CityName]City)}
}

// A single city within the map.
type City struct {
	// Roads out of the city.
	roads map[Direction]CityName
}

// Create a new city, with no roads.
func NewCity() *City {
	return &City{make(map[Direction]CityName)}
}

func (c *City) AddRoad(direction Direction, name CityName) error {
	_, ok := c.roads[direction]
	if ok {
		// Road direction already defined.
		return errors.Errorf("City already has a road in direction %s", direction)
	}
	c.roads[direction] = name
	return nil
}

// Add a city to a map, will return an error if the city already exists.
func (m *Map) AddCity(cityName CityName, city City) error {
	_, ok := m.cities[cityName]
	if ok {
		// City already exists, cannot have the same name.
		return errors.Errorf("Duplicate city name: %s", cityName)
	}
	m.cities[cityName] = city
	m.destroyed[cityName] = false
	for _, destCity := range city.roads {
		if _, ok := m.destroyed[destCity]; !ok {
			m.destroyed[destCity] = false
		}
	}
	return nil
}

// Get a list of known cities, cities will be returned in unknown order.
func (m *Map) KnownCities() []CityName {
	var ret []CityName

	for n, c := range m.cities {
		ret = append(ret, n)
		for _, n := range c.roads {
			ret = append(ret, n)
		}
	}

	return ret
}

func (m *Map) PrettyPrint(w io.Writer) {
	for cityName, c := range m.cities {
		if m.IsCityDestroyed(cityName) {
			// City is destroyed, skip it.
			continue
		}
		roads := make([]string, len(c.roads))
		i := 0
		for dir, destCity := range c.roads {
			if m.IsCityDestroyed(destCity) {
				// Destination city is destroyed, skip it.
				continue
			}
			roads[i] = fmt.Sprintf("%s=%s", dir, destCity)
			i++
		}
		fmt.Fprintf(w, "%s %s\n", cityName, strings.Join(roads, " "))
	}
}

// Is this city destroyed? Will return false for unknown cities.
func (m *Map) IsCityDestroyed(cityName CityName) bool {
	d, ok := m.destroyed[cityName]
	if !ok {
		return false
	}
	return d
}

// Destroy the given city.
func (m *Map) DestroyCity(cityName CityName) error {
	destroyed, ok := m.destroyed[cityName]
	if !ok {
		return errors.Errorf("City not known: %s", cityName)
	}
	if destroyed {
		return errors.Errorf("City already destroyed: %s", cityName)
	}
	m.destroyed[cityName] = true
	return nil
}

// Find non-destroyed cities connected to the given city.
func (m *Map) ConnectedCities(cityName CityName) []CityName {
	var cities []CityName
	city, ok := m.cities[cityName]
	if !ok {
		return cities
	}
	for _, destName := range city.roads {
		if !m.IsCityDestroyed(destName) {
			cities = append(cities, destName)
		}
	}
	return cities
}
