package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	input1 := `
# We have some comments, and after this blank line.

Foo north=Bar west=Baz south=Qu-ux
`
	var parser Parser
	m, err := parser.Parse(strings.NewReader(input1))
	require.NoError(t, err)

	assert.Equal(t, 1, len(m.Cities))
	city, ok := m.Cities[CityName("Foo")]
	require.True(t, ok)
	assert.Equal(t, map[Direction]CityName{
		North: CityName("Bar"),
		West:  CityName("Baz"),
		South: CityName("Qu-ux"),
	}, city.Roads)
}
