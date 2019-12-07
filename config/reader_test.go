package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestSuccessfulReader(t *testing.T) {
	var parser Parser
	var err error

	input1 := `
# We have some comments, and after this blank line.

Foo north=Bar west=Baz south=Qu-ux
`
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

func TestInvalidReader(t *testing.T) {
	var parser Parser
	var err error

	input1 := `
# This one is bad, as a city name has a space, making the format invalid.

Quack north=blah source=another one
`
	_, err = parser.Parse(strings.NewReader(input1))
	require.Error(t, err)

	input2 := `
# This one is bad, as it repeats a city name.

Which north=What south=Where
Which north=That south=Other
`
	_, err = parser.Parse(strings.NewReader(input2))
	require.Error(t, err)
}
