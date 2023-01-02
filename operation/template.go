package operation

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/aziis98/cabret"
)

func init() {
	registerType("template", &Template{})
}

type Template struct {
	Engine string
}

// Configure will configure this operation
//
//	use: template
//	engine: <which engine to use> # required, can be "html" or "text"
func (op *Template) Configure(options map[string]any) error {
	var err error
	op.Engine, err = getKey[string](options, "engine")
	if err != nil {
		return err
	}

	return nil
}

func (op *Template) ProcessList(items []cabret.Content) ([]cabret.Content, error) {
	var t bytes.Buffer

	// concatenate all templates
	for _, item := range items {
		t.Write(item.Data)
	}

	tmpl := t.String()

	var data bytes.Buffer
	switch op.Engine {
	case "html":
		if err := op.RenderHtml(tmpl, items, &data); err != nil {
			return nil, err
		}
	case "text":
		if err := op.RenderText(tmpl, items, &data); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(`unknown format "%s"`, op.Engine)
	}

	return []cabret.Content{
		{
			Type:     items[0].Type,
			Metadata: cabret.Map{},
			Data:     data.Bytes(),
		},
	}, nil
}

func (op *Template) RenderHtml(tmpl string, items []cabret.Content, w io.Writer) error {
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		return err
	}

	if err := t.ExecuteTemplate(w, "template", map[string]any{"Items": items}); err != nil {
		return err
	}

	return nil
}

func (op *Template) RenderText(tmpl string, items []cabret.Content, w io.Writer) error {
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		return err
	}

	if err := t.ExecuteTemplate(w, "template", map[string]any{"Items": items}); err != nil {
		return err
	}

	return nil
}
