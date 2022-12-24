package main

import (
	"log"

	"github.com/alecthomas/repr"
	"github.com/aziis98/cabret/config"
)

func main() {
	log.Printf("Rendering current project")

	site, err := config.ParseCabretfile("./Cabretfile.yaml")
	if err != nil {
		log.Fatal(err)
	}

	repr.Println(site)
}
