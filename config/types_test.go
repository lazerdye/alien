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

func testMap() (*Map, error) {
	m := NewMap()
	c1 := NewCity()
	if err := c1.AddRoad(North, chaCity); err != nil {
		return nil, err
	}
	if err := c1.AddRoad(South, blahCity); err != nil {
		return nil, err
	}
	if err := m.AddCity(quackCity, *c1); err != nil {
		return nil, err
	}
	return m, nil
}

// Check if we can get the list of known cities from a map.
func TestKnownCities(t *testing.T) {
	m, err := testMap()
	require.NoError(t, err)

	assert.ElementsMatch(t, m.KnownCities(), []CityName{chaCity, blahCity, quackCity})
}

// Check destroying a given city.
func TestDestroyCity(t *testing.T) {
	m, err := testMap()
	require.NoError(t, err)

	// Check initial status of blahCity.
	assert.False(t, m.IsCityDestroyed(blahCity))

	// Try destroying a nonexistant city.
	err = m.DestroyCity(CityName("NotHere"))
	assert.Error(t, err)

	// Destroy blah.
	err = m.DestroyCity(blahCity)
	assert.NoError(t, err)
	assert.True(t, m.IsCityDestroyed(blahCity))

	// Destroy blah again.
	err = m.DestroyCity(blahCity)
	assert.Error(t, err)
}
