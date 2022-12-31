package exec

import (
	"github.com/aziis98/cabret"
	"github.com/aziis98/cabret/config"
	"github.com/aziis98/cabret/pipeline"
)

func Execute(cfg *config.Cabretfile) error {
	for _, p := range cfg.Build {
		ops, err := pipeline.Parse(p)
		if err != nil {
			return err
		}

		if _, err := pipeline.Process([]cabret.Content{}, ops); err != nil {
			return err
		}
	}

	return nil
}
