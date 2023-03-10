package operation

import (
	"fmt"
	"reflect"

	"github.com/aziis98/cabret"
)

// ShortFormValueKey is used by some operations that support an inline form
const ShortFormValueKey = "value"

var registry = map[string]reflect.Type{}

func registerType(name string, op cabret.Operation) {
	typ := reflect.TypeOf(op).Elem()
	registry[name] = typ
}

func NewWithName(name string) (cabret.Operation, error) {
	typ, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf(`no registered operation named %q`, name)
	}

	return reflect.New(typ).Interface().(cabret.Operation), nil
}
