package alien

import (
	"fmt"
	"io"

	"github.com/lazerdye/alien/config"
)

type RandomInt interface {
	Intn(max int) int
}

type Alien struct {
	id   int
	city config.CityName
}

func (a *Alien) InitInRandomLocation(id int, m *config.Map, randInt RandomInt) {
	knownCities := m.KnownCities()
	index := randInt.Intn(len(knownCities))
	a.id = id
	a.city = knownCities[index]
}

func (a *Alien) PrettyPrint(w io.Writer) {
	fmt.Fprintf(w, "%d %s\n", a.id, a.city)
}

func (a *Alien) City() CityName {
    return a.city
}
