package config

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"strings"
)

// Tool to parse config files.
type Parser struct{}

// Parse a single line, returns (NilCityName, nil, nil) if the line is empty.
func (p *Parser) parseLine(line string) (CityName, *City, error) {
	// Nice to have: remove # and anything after it.
	parts := strings.SplitN(line, "#", 2)
	// Cut out any extra whitespaces.
	trimmed := strings.TrimSpace(parts[0])
	// If we're left with an empty string, the line was empty or had only comments.
	if len(trimmed) == 0 {
		return NilCityName, nil, nil
	}
	mapEntryParts := strings.Split(trimmed, " ")
	thisCityName := CityNameFromString(mapEntryParts[0])
	if thisCityName == NilCityName {
		return NilCityName, nil, errors.New("Invalid CityName")
	}

	roads := make(map[Direction]CityName, len(mapEntryParts)-1)
	for i := 1; i < len(mapEntryParts); i++ {
		directionAndCity := strings.SplitN(mapEntryParts[i], "=", 2)
		direction := DirectionFromString(directionAndCity[0])
		if direction == NilDirection {
			return NilCityName, nil, errors.Errorf("Invalid Direction: %s", directionAndCity[0])
		}
		cityName := CityNameFromString(directionAndCity[1])
		if cityName == NilCityName {
			return NilCityName, nil, errors.Errorf("Invalid CityName: %s", directionAndCity[1])
		}
		roads[direction] = cityName
	}
	city := City{roads}

	return thisCityName, &city, nil
}

// Parse the given file.
func (p *Parser) Parse(r io.Reader) (*Map, error) {
	scanner := bufio.NewScanner(r)
	m := NewMap()
	for lineno := 1; scanner.Scan(); lineno++ {
		cityName, city, err := p.parseLine(scanner.Text())
		if err != nil {
			return nil, errors.Wrapf(err, "Scanner line %d", lineno)
		}
		if cityName == NilCityName {
			// Empty line.
			continue
		}
		if err := m.AddCity(cityName, *city); err != nil {
			return nil, errors.Wrapf(err, "AddCity line %d", lineno)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return m, nil
}
