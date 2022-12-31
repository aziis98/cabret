package runner

import (
	"fmt"

	"github.com/aziis98/cabret"
	"github.com/aziis98/cabret/config"
	"github.com/aziis98/cabret/parse"
)

func RunConfig(cfg *config.Cabretfile) error {
	for _, p := range cfg.Build {
		ops, err := parse.ParsePipeline(p)
		if err != nil {
			return err
		}

		if _, err := RunPipeline([]cabret.Content{}, ops); err != nil {
			return err
		}
	}

	return nil
}

func RunPipeline(contents []cabret.Content, ops []cabret.Operation) ([]cabret.Content, error) {
	for _, op := range ops {
		var err error
		contents, err = RunOperation(op, contents)
		if err != nil {
			return nil, err
		}
	}

	return contents, nil
}

func RunOperation(op cabret.Operation, inputs []cabret.Content) ([]cabret.Content, error) {
	switch op := op.(type) {
	case cabret.ListOperation:
		return op.ProcessList(inputs)

	case cabret.ItemOperation:
		outputs := []cabret.Content{}
		for _, item := range inputs {
			result, err := op.ProcessItem(item)
			if err != nil {
				return nil, err
			}

			// skip terminal operations
			if result == nil {
				continue
			}
			outputs = append(outputs, *result)
		}
		return outputs, nil
	default:
		return nil, fmt.Errorf(`invalid operation type %T`, op)
	}
}
