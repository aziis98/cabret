package operation

import (
	"fmt"

	"github.com/aziis98/cabret"
)

func init() {
	registerType("categorize", &Categorize{})
}

type Categorize struct {
	Key string

	// Operation to be executed for each category
	Operation cabret.Operation
}

func (op *Categorize) Load(config map[string]any) error {
	{
		v, ok := config["key"]
		if !ok {
			return fmt.Errorf(`missing "key" field`)
		}
		key, ok := v.(string)
		if !ok {
			return fmt.Errorf(`expected string but got "%v" of type %T`, v, v)
		}

		op.Key = key
	}

	return nil
}

func (op *Categorize) Process(content cabret.Content) (*cabret.Content, error) {
	return nil, nil
}
