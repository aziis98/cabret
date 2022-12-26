package layout

import (
	"bytes"
	"html/template"
)

var _ Template = HtmlTemplate{}

type HtmlTemplate struct {
	*template.Template
}

func NewHtmlTemplate(files ...string) (*HtmlTemplate, error) {
	t, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return &HtmlTemplate{t}, nil
}

func (t HtmlTemplate) Render(ctx map[string]any) ([]byte, error) {
	var b bytes.Buffer

	if err := t.Template.Execute(&b, ctx); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
