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
	cities map[CityName]City
}

// Create a new, empty map.
func NewMap() *Map {
	return &Map{make(map[CityName]City)}
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
		roads := make([]string, len(c.roads))
		i := 0
		for dir, destCity := range c.roads {
			roads[i] = fmt.Sprintf("%s=%s", dir, destCity)
			i++
		}
		fmt.Fprintf(w, "%s %s\n", cityName, strings.Join(roads, " "))
	}
}
