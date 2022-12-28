package operation

import (
	"fmt"
	"os"

	gopath "path"

	"github.com/aziis98/cabret"
	"github.com/aziis98/cabret/path"
)

var _ cabret.FlatMapOperation = Target{}

type Target struct {
	PathTemplate string
}

func (op Target) FlatMap(c cabret.Content) (*cabret.Content, error) {
	mr, ok := c.Metadata[cabret.MatchResult].(map[string]string)
	if !ok {
		return nil, fmt.Errorf(`invalid match result type %T`, c.Metadata[cabret.MatchResult])
	}

	target := path.RenderTemplate(op.PathTemplate, mr)

	if err := os.MkdirAll(gopath.Dir(target), 0777); err != nil {
		return nil, err
	}
	if err := os.WriteFile(target, c.Data, 0666); err != nil {
		return nil, err
	}

	c.Metadata["Target"] = op.PathTemplate
	return &c, nil
}
