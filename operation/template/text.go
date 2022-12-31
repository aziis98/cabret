package template

import (
	"bytes"
	goTextTemplate "text/template"
)

var _ Template = TextTemplate{}

type TextTemplate struct {
	*goTextTemplate.Template
}

func NewTextTemplate(files ...string) (*TextTemplate, error) {
	t, err := goTextTemplate.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return &TextTemplate{t}, nil
}

func (t TextTemplate) Render(ctx map[string]any) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.Template.Execute(&buf, ctx); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
