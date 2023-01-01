package main

import (
	"log"

	"github.com/aziis98/cabret/config"
	"github.com/aziis98/cabret/runner"
	"github.com/spf13/pflag"
)

func main() {
	log.SetFlags(0)

	optConfig := pflag.StringP("config", "c", "Cabretfile.yaml", `which configuration file to use`)
	pflag.Parse()

	cabretfile, err := config.ReadCabretfile(*optConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err := runner.RunConfig(cabretfile); err != nil {
		log.Fatal(err)
	}
}
