package alien

import (
	"testing"

	"github.com/lazerdye/alien/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	chaCity   = config.CityName("cha")
	blahCity  = config.CityName("blah")
	quackCity = config.CityName("quack")
)

func testMap() (*config.Map, error) {
	m := config.NewMap()
	c1 := config.NewCity()
	if err := c1.AddRoad(config.North, chaCity); err != nil {
		return nil, err
	}
	if err := c1.AddRoad(config.South, blahCity); err != nil {
		return nil, err
	}
	if err := m.AddCity(quackCity, *c1); err != nil {
		return nil, err
	}
	return m, nil
}

type fakeRandom struct {
	returnValue int
}

func (r fakeRandom) Intn(n int) int {
	return r.returnValue % n
}

func TestInitInRandomLocation(t *testing.T) {
	m, err := testMap()
	require.NoError(t, err)
	alien1 := NewRandomAlien(1, m, fakeRandom{0})
	alien2 := NewRandomAlien(2, m, fakeRandom{1})
	assert.Equal(t, 1, alien1.id)
	assert.Equal(t, quackCity, alien1.city)
	assert.Equal(t, 2, alien2.id)
	assert.Equal(t, chaCity, alien2.city)
}

func TestFight(t *testing.T) {

	m, err := testMap()
	require.NoError(t, err)

	alien1 := NewRandomAlien(1, m, fakeRandom{0})
	alien2 := NewRandomAlien(2, m, fakeRandom{0})
	alien3 := NewRandomAlien(3, m, fakeRandom{1})

	err = Fight(m, []*Alien{alien1, alien2})
	require.NoError(t, err)

	assert.True(t, m.IsCityDestroyed(quackCity))
	assert.False(t, m.IsCityDestroyed(chaCity))
	assert.True(t, alien1.IsDestroyed())
	assert.True(t, alien2.IsDestroyed())
	assert.False(t, alien3.IsDestroyed())
	assert.Equal(t, chaCity, alien3.City())
}
