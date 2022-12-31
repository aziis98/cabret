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
	Key string
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
				"Category": name,
				"Items":    contents,
			},
		})
	}

	return result, nil
}
