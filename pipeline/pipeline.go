package pipeline

import (
	"fmt"

	"github.com/aziis98/cabret"
	"github.com/aziis98/cabret/config"
	"github.com/aziis98/cabret/operation"
)

func switchMap(m map[string]any, v *any) func(k string) bool {
	return func(k string) bool {
		val, ok := m[k]
		if ok {
			*v = val
		}
		return ok
	}
}

func Parse(p config.Pipeline) ([]cabret.Operation, error) {
	ops := []cabret.Operation{}

	for _, opConfig := range p.Pipeline {
		var v any
		has := switchMap(opConfig, &v)

		switch {
		case has("source"):
			value, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf(`expected string but got "%v" of type %T`, v, v)
			}

			opConfig[operation.ShortFormValueKey] = value
			op := &operation.Source{}
			if err := op.Load(opConfig); err != nil {
				return nil, err
			}

			ops = append(ops, op)

		case has("target"):
			value, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf(`expected string but got "%v" of type %T`, v, v)
			}

			opConfig[operation.ShortFormValueKey] = value
			op := &operation.Target{}
			if err := op.Load(opConfig); err != nil {
				return nil, err
			}

			ops = append(ops, op)

		case has("use"):
			name, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf(`expected string but got "%v" of type %T`, v, v)
			}

			op, err := operation.Build(name, opConfig)
			if err != nil {
				return nil, err
			}

			ops = append(ops, op)

		default:
			return nil, fmt.Errorf(`pipeline entry is missing one of "use", "source" or "target", got %#v`, opConfig)
		}
	}

	return ops, nil
}

func Process(contents []cabret.Content, ops []cabret.Operation) ([]cabret.Content, error) {
	for _, op := range ops {
		var err error
		contents, err = cabret.ProcessOperation(op, contents)
		if err != nil {
			return nil, err
		}
	}

	return contents, nil
}
