package template

import (
	"bytes"
	goHtmlTemplate "html/template"
)

var _ Template = HtmlTemplate{}

type HtmlTemplate struct {
	*goHtmlTemplate.Template
}

func NewHtmlTemplate(files ...string) (*HtmlTemplate, error) {
	t, err := goHtmlTemplate.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return &HtmlTemplate{t}, nil
}

func (t HtmlTemplate) Render(ctx map[string]any) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.Template.Execute(&buf, ctx); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
