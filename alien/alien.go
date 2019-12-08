package alien

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/lazerdye/alien/config"
)

// The state of an alien.
type Alien struct {
	id        int
	destroyed bool
	city      config.CityName
}

// Interface to allow easy mocking of generating a random inteter.
type RandomInt interface {
	Intn(max int) int
}

// Initialize the alien in a random location.
func NewRandomAlien(id int, m *config.Map, randInt RandomInt) *Alien {
	knownCities := m.KnownCities()
	index := randInt.Intn(len(knownCities))
	return &Alien{
		id:        id,
		destroyed: false,
		city:      knownCities[index],
	}
}

// Is this alien destroyed?
func (a *Alien) IsDestroyed() bool {
	return a.destroyed
}

// Print out the alien info.
func (a *Alien) PrettyPrint(w io.Writer) {
	if !a.IsDestroyed() {
		fmt.Fprintf(w, "%d %s\n", a.id, a.city)
	} else {
		fmt.Fprintf(w, "%d %s (destroyed)\n", a.id, a.city)
	}
}

// Destroy this alien, this can happen only once per alien.
func (a *Alien) Destroy() error {
	if a.destroyed {
		return errors.Errorf("Alien %d already destroyed", a.id)
	}
	a.destroyed = true
	return nil
}

// Return this alien's id.
func (a *Alien) ID() int {
	return a.id
}

// Return which city this alien is in.
func (a *Alien) City() config.CityName {
	return a.city
}

// Move this alien to a random allowed location.
func (a *Alien) Move(m *config.Map, r RandomInt) {
	roads := m.ConnectedCities(a.City())
	if len(roads) > 0 {
		newCity := roads[r.Intn(len(roads))]
		a.city = newCity
	}
}

// Perform a fight between a group of aliens.
func Fight(m *config.Map, aliens []*Alien) error {
	if len(aliens) <= 1 {
		return errors.New("A fight requires at least two aliens.")
	}
	// Assume all of these aliens are in the same city.
	city := aliens[0].City()
	// This city is destroyed.
	if err := m.DestroyCity(city); err != nil {
		return err
	}
	// Now all of the aliens in the city are destroyed.
	alienDescription := make([]string, len(aliens))
	for i, a := range aliens {
		if a.City() != city {
			return errors.New("Cannot fight between aliens in different cities")
		}
		if err := a.Destroy(); err != nil {
			return err
		}
		alienDescription[i] = fmt.Sprintf("alien %d", a.ID())
	}
	fmt.Printf("%s has been destroyed by %s!\n", city, strings.Join(alienDescription, " and "))
	return nil
}
