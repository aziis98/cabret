package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Operation should have at least one key in "source", "use", "target". The remaining keys are options for that operation
type Operation map[string]any

type Pipeline struct {
	Pipeline []Operation `yaml:"pipeline"`
}

type BuildOptions struct {
	// Excludes lists files and folders to globally exclude from compilation
	Excludes []string `yaml:",omitempty"`
}

type Cabretfile struct {
	Options BuildOptions
	Build   []Pipeline
}

func ReadCabretfile(file string) (*Cabretfile, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	site := new(Cabretfile)
	if err := yaml.NewDecoder(f).Decode(site); err != nil {
		return nil, err
	}

	return site, nil
}
