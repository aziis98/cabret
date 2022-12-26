package operation

import (
	"bytes"

	"github.com/aziis98/cabret"
	"github.com/iancoleman/strcase"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type Markdown struct {
	Options map[string]any
}

func (op Markdown) Process(content cabret.Content) (*cabret.Content, error) {
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
