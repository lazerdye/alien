package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/lazerdye/alien/config"
)

var (
	mapFile = kingpin.Flag("mapFile", "Location of map config file").Required().String()

	run    = kingpin.Command("run", "Run the simulator")
	verify = kingpin.Command("verify", "Verify the config file")
)

func verifyConfig() error {
	file, err := os.Open(*mapFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var parser config.Parser
	m, err := parser.Parse(file)
	if err != nil {
		return err
	}

	log.Infof("Verificaton successful: %+v", *m)

	return nil
}

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Terence Haddock")
	kingpin.CommandLine.Help = "Run the alien simulator"

	switch kingpin.Parse() {
	case "run":
		log.Fatal("Not implemented")
	case "verify":
		if err := verifyConfig(); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
