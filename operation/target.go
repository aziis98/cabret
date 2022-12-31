package operation

import (
	"fmt"
	"log"
	"os"

	gopath "path"

	"github.com/aziis98/cabret"
	"github.com/aziis98/cabret/path"
)

func init() {
	registerType("target", &Target{})
}

type Target struct {
	PathTemplate string
}

func (op *Target) Load(config map[string]any) error {
	if v, ok := config[ShortFormValueKey]; ok {
		template, ok := v.(string)
		if !ok {
			return fmt.Errorf(`expected pattern but got "%v" of type %T`, v, v)
		}

		op.PathTemplate = template
		return nil
	}

	return fmt.Errorf(`invalid config for "target": %#v`, config)
}

func (op Target) ProcessItem(c cabret.Content) (*cabret.Content, error) {
	mr, ok := c.Metadata[cabret.MatchResultKey].(map[string]string)
	if !ok {
		return nil, fmt.Errorf(`invalid match result type %T`, c.Metadata[cabret.MatchResultKey])
	}

	target := path.RenderTemplate(op.PathTemplate, mr)

	log.Printf(`[operation.Target] writing "%s"`, target)

	if err := os.MkdirAll(gopath.Dir(target), 0777); err != nil {
		return nil, err
	}
	if err := os.WriteFile(target, c.Data, 0666); err != nil {
		return nil, err
	}

	c.Metadata["Target"] = op.PathTemplate
	return &c, nil
}
