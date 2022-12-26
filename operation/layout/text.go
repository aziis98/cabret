package layout

import (
	"bytes"
	"text/template"
)

var _ Template = TextTemplate{}

type TextTemplate struct {
	*template.Template
}

func NewTextTemplate(files ...string) (*TextTemplate, error) {
	t, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return &TextTemplate{t}, nil
}

func (t TextTemplate) Render(ctx map[string]any) ([]byte, error) {
	var b bytes.Buffer

	if err := t.Template.Execute(&b, ctx); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
