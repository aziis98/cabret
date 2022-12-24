package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Operation is an enum of various operations
type Operation struct {
	Layout string `yaml:"layout,omitempty"`
	Target string `yaml:"target,omitempty"`
	Plugin string `yaml:"plugin,omitempty"`
}

type EntryPoint struct {
	Source   string      `yaml:"source"`
	Pipeline []Operation `yaml:"pipeline"`
}

type Config struct {
	Output string `yaml:"output,omitempty"`
}

// Site has some configuration for the
type Site struct {
	Config      Config       `yaml:"config,omitempty"`
	EntryPoints []EntryPoint `yaml:"entry-points"`
}

func ParseCabretfile(file string) (*Site, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	site := new(Site)
	if err := yaml.NewDecoder(f).Decode(&site); err != nil {
		return nil, err
	}

	return site, nil
}
