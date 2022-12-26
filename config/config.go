package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Operation is an enum of various operations
type Operation struct {
	Layout string `yaml:",omitempty"`
	Target string `yaml:",omitempty"`
	Plugin string `yaml:",omitempty"`

	Options map[string]any `yaml:",omitempty"`
}

type EntryPoint struct {
	Source   string      `yaml:",omitempty"`
	Pipeline []Operation `yaml:",omitempty"`
}

type Options struct {
	Excludes []string `yaml:",omitempty"`
	// Include []string `yaml:",omitempty"`
	Output string `yaml:",omitempty"`
}

// Cabretfile has some configuration for the
type Cabretfile struct {
	Options     Options       `yaml:",omitempty"`
	EntryPoints []*EntryPoint `yaml:"entryPoints"`
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
