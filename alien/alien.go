package alien

import (
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/lazerdye/alien/config"
)

type RandomInt interface {
	Intn(max int) int
}

type Alien struct {
	id        int
	destroyed bool
	city      config.CityName
}

func (a *Alien) InitInRandomLocation(id int, m *config.Map, randInt RandomInt) {
	knownCities := m.KnownCities()
	index := randInt.Intn(len(knownCities))
	a.id = id
	a.destroyed = false
	a.city = knownCities[index]
}

func (a *Alien) IsDestroyed() bool {
	return a.destroyed
}

func (a *Alien) PrettyPrint(w io.Writer) {
	fmt.Fprintf(w, "%d %s\n", a.id, a.city)
}

func (a *Alien) Destroy() error {
	if a.destroyed {
		return errors.Errorf("Alien %d already destroyed", a.id)
	}
	a.destroyed = true
	return nil
}

func (a *Alien) ID() int {
	return a.id
}

func (a *Alien) City() config.CityName {
	return a.city
}

func (a *Alien) MoveTo(newCity config.CityName) {
	a.city = newCity
}
