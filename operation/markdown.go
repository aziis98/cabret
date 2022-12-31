package operation

import (
	"bytes"
	"log"

	"github.com/aziis98/cabret"
	"github.com/iancoleman/strcase"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func init() {
	registerType("markdown", &Markdown{})
}

type Markdown struct {
	Options map[string]any
}

func (op *Markdown) Load(config map[string]any) error {
	return nil
}

func (op Markdown) ProcessItem(content cabret.Content) (*cabret.Content, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	var buf bytes.Buffer

	log.Printf(`[operation.Markdown] rendering markdown`)

	context := parser.NewContext()
	if err := md.Convert(content.Data, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	frontmatter := meta.Get(context)
	for k, v := range frontmatter {
		content.Metadata[strcase.ToCamel(k)] = v
	}

	return &cabret.Content{
		Type:     HtmlMimeType,
		Data:     buf.Bytes(),
		Metadata: content.Metadata,
	}, nil
}
