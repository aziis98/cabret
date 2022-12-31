package operation

import (
	"fmt"
	"log"

	"github.com/aziis98/cabret"
	"github.com/iancoleman/strcase"
)

func init() {
	registerType("categorize", &Categorize{})
}

type Categorize struct {
	Key              string
	CategoryVariable string
}

func getKey[T any](m map[string]any, key string, defaultValue ...T) (T, error) {
	v, ok := m[key]
	if !ok {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}

		var zero T
		return zero, fmt.Errorf(`missing "%s" field`, key)
	}

	value, ok := v.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf(`expected %T but got "%v" of type %T`, zero, v, v)
	}

	return value, nil
}

func (op *Categorize) Configure(config map[string]any) error {
	var err error

	op.Key, err = getKey[string](config, "key")
	if err != nil {
		return err
	}

	op.CategoryVariable, err = getKey(config, "bind", "Category")
	if err != nil {
		return err
	}

	return nil
}

func (op *Categorize) ProcessList(contents []cabret.Content) ([]cabret.Content, error) {
	key := strcase.ToCamel(op.Key)

	categories := map[string][]cabret.Content{}

	for _, content := range contents {
		v, ok := content.Metadata[key]
		if !ok {
			log.Printf(`[operation.Categorize] item has no categories`)
			continue
		}

		cats, ok := v.([]any)
		if !ok {
			return nil, fmt.Errorf(`expected a list but got "%v" of type %T`, v, v)
		}

		for _, v := range cats {
			cat, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf(`expected a string but got "%v" of type %T`, v, v)
			}

			log.Printf(`[operation.Categorize] found item with category "%s"`, cat)

			categories[cat] = append(categories[cat], content)
		}
	}

	result := []cabret.Content{}
	for name, contents := range categories {
		log.Printf(`[operation.Categorize] generating category with %v item(s)`, len(contents))

		result = append(result, cabret.Content{
			Type: "metadata-only",
			Metadata: cabret.Map{
				op.CategoryVariable: name,
				"Items":             contents,
			},
		})
	}

	return result, nil
}
