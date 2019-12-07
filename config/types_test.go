package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	chaCity   = CityName("Cha")
	blahCity  = CityName("Blah")
	quackCity = CityName("quack")
)

// Check if we can get the list of known cities from a map.
func TestKnownCities(t *testing.T) {
	var err error

	m := NewMap()
	c1 := NewCity()
	err = c1.AddRoad(North, chaCity)
	require.NoError(t, err)
	err = c1.AddRoad(South, blahCity)
	require.NoError(t, err)

	err = m.AddCity(quackCity, *c1)
	require.NoError(t, err)

	assert.ElementsMatch(t, m.KnownCities(), []CityName{chaCity, blahCity, quackCity})
}
