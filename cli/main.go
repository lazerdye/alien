package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/lazerdye/alien/alien"
	"github.com/lazerdye/alien/config"
)

var (
	mapFile     = kingpin.Flag("map", "Location of map config file").Required().String()
	infoLogging = kingpin.Flag("info-logging", "Turn info logging on").Bool()

	run          = kingpin.Command("run", "Run the simulator")
	runNumAliens = run.Arg("num-aliens", "Number of aliens").Required().Int()
	runMaxTime   = run.Arg("max-time", "Maximum time").Default("10000").Int()
	verify       = kingpin.Command("verify", "Verify the config file")
)

func loadConfig(mapFile string) (*config.Map, error) {
	file, err := os.Open(mapFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var parser config.Parser
	m, err := parser.Parse(file)
	if err != nil {
		return nil, err
	}

	log.Infof("loadConfig successful: %+v", *m)

	return m, nil
}

func doRun(m *config.Map, numAliens int, maxTime int) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generate alien list.
	aliens := make([]alien.Alien, numAliens)
	for i := 0; i < numAliens; i++ {
		aliens[i].InitInRandomLocation(i+1, m, r)
	}
	for t := 0; t < maxTime; t++ {
		if err := runSingleLoop(t, m, aliens, r); err != nil {
			return err
		}
	}
	fmt.Println("==== Final map")
	m.PrettyPrint(os.Stdout)
	return nil
}

func runSingleLoop(t int, m *config.Map, aliens []alien.Alien, r *rand.Rand) error {
	if log.GetLevel() >= log.InfoLevel {
		fmt.Fprintf(os.Stderr, "=== %d MAP\n", t)
		m.PrettyPrint(os.Stderr)
		fmt.Fprintf(os.Stderr, "=== %d ALIENS\n", t)
		for _, a := range aliens {
			if a.IsDestroyed() {
				continue
			}
			a.PrettyPrint(os.Stderr)
		}
	}
	// Generate a map city->[]aliens
	cityToAlien := make(map[config.CityName][]*alien.Alien)
	for i := 0; i < len(aliens); i++ {
		a := &aliens[i]
		if a.IsDestroyed() {
			// Alien is destroyed, skip it.
			continue
		}
		cityToAlien[a.City()] = append(cityToAlien[a.City()], a)
	}
	// Find each city with more than one alien
	// Destroy said city, and aliens
	for city, aliensInCity := range cityToAlien {
		if len(aliensInCity) > 1 {
			log.Infof("City %s has too many aliens: %v", city, aliensInCity)
			if err := m.DestroyCity(city); err != nil {
				return err
			}
			for _, a := range aliensInCity {
				if err := a.Destroy(); err != nil {
					return err
				}
			}
			alienDescription := make([]string, len(aliensInCity))
			for i, a := range aliensInCity {
				alienDescription[i] = fmt.Sprintf("alien %d", a.ID())
			}
			fmt.Printf("%s has been destroyed by %s!\n", city, strings.Join(alienDescription, " and "))
		}
	}
	// Move the remaining aliens
	for i := 0; i < len(aliens); i++ {
		if aliens[i].IsDestroyed() {
			continue
		}
		directions := m.AvailableDirections(aliens[i].City())
		if len(directions) > 0 {
			newDir := directions[r.Intn(len(directions))]
			aliens[i].MoveTo(newDir)
		}
	}
	return nil
}

func initLogging(infoLogging bool) {
	if infoLogging {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Terence Haddock")
	kingpin.CommandLine.Help = "Run the alien simulator"

	switch kingpin.Parse() {
	case "run":
		initLogging(*infoLogging)
		m, err := loadConfig(*mapFile)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if err := doRun(m, *runNumAliens, *runMaxTime); err != nil {
			log.Fatalf("%v", err)
		}
	case "verify":
		initLogging(*infoLogging)
		if _, err := loadConfig(*mapFile); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
