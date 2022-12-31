package main

import (
	"log"

	"github.com/aziis98/cabret/config"
	"github.com/aziis98/cabret/runner"
)

func main() {
	log.SetFlags(0)

	cabretfile, err := config.ReadCabretfile("./Cabretfile.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if err := runner.RunConfig(cabretfile); err != nil {
		log.Fatal(err)
	}
}
