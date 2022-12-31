package operation

import (
	"io"
	"log"

	"github.com/aziis98/cabret"
	"github.com/iancoleman/strcase"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func init() {
	registerType("frontmatter", &Frontmatter{})
}

type Frontmatter struct {
	Options map[string]any
}

func (op *Frontmatter) Configure(config map[string]any) error {
	return nil
}

func (op *Frontmatter) ProcessItem(content cabret.Content) (*cabret.Content, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	log.Printf(`[operation.Frontmatter] reading frontmatter`)

	context := parser.NewContext()
	if err := md.Convert(content.Data, io.Discard, parser.WithContext(context)); err != nil {
		panic(err)
	}

	frontmatter := meta.Get(context)
	for k, v := range frontmatter {
		content.Metadata[strcase.ToCamel(k)] = v
	}

	return &content, nil
}
