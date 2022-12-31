package operation

import (
	"fmt"
	"log"
	"reflect"

	"github.com/aziis98/cabret"
)

// ShortFormValueKey is used by some operations that support an inline form
const ShortFormValueKey = "value"

var registry = map[string]reflect.Type{}

func registerType(name string, op cabret.Operation) {
	typ := reflect.TypeOf(op).Elem()
	log.Printf(`[operation] registered type "%v"`, typ)
	registry[name] = typ
}

func Build(name string, options map[string]any) (cabret.Operation, error) {
	typ, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf(`no registered operation named %q`, name)
	}

	op := reflect.New(typ).Interface().(cabret.Operation)
	if err := op.Load(options); err != nil {
		return nil, err
	}

	return op, nil
}
