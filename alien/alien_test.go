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
	var alien Alien
	alien.InitInRandomLocation(5, m, fakeRandom{0})
	assert.Equal(t, 5, alien.id)
	assert.Equal(t, quackCity, alien.city)
}
