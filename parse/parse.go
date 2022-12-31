package parse

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

func ParsePipeline(p config.Pipeline) ([]cabret.Operation, error) {
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

			op, err := ParseOperation("source", opConfig)
			if err != nil {
				return nil, err
			}

			ops = append(ops, op)

		case has("target"):
			value, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf(`expected string but got "%v" of type %T`, v, v)
			}

			opConfig[operation.ShortFormValueKey] = value

			op, err := ParseOperation("target", opConfig)
			if err != nil {
				return nil, err
			}

			ops = append(ops, op)

		case has("use"):
			name, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf(`expected string but got "%v" of type %T`, v, v)
			}

			op, err := ParseOperation(name, opConfig)
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

func ParseOperation(name string, options map[string]any) (cabret.Operation, error) {
	op, err := operation.NewWithName(name)
	if err != nil {
		return nil, err
	}

	if err := op.Configure(options); err != nil {
		return nil, err
	}

	return op, nil
}
