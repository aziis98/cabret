package main

import (
	"log"

	"github.com/aziis98/cabret/config"
	"github.com/aziis98/cabret/exec"
)

func main() {
	log.SetFlags(0)
	log.Printf("Rendering current project")

	site, err := config.ReadCabretfile("./Cabretfile.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// repr.Println(site)

	if err := exec.Execute(site); err != nil {
		log.Fatal(err)
	}

}
