package main

import (
	"log"

	"github.com/aziis98/cabret/config"
	"github.com/aziis98/cabret/exec"
)

func main() {
	log.SetFlags(0)

	cabretfile, err := config.ReadCabretfile("./Cabretfile.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// repr.Println(cabretfile)

	if err := exec.Execute(cabretfile); err != nil {
		log.Fatal(err)
	}
}
