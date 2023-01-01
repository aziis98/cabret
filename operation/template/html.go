package template

import (
	"bytes"
	goHtmlTemplate "html/template"
)

func safeHTML(str string) goHtmlTemplate.HTML {
	return goHtmlTemplate.HTML(str)
}

var fn = goHtmlTemplate.FuncMap{
	"html": safeHTML,
}

type HtmlTemplate struct {
	*goHtmlTemplate.Template
}

// func NewHtmlTemplateFromReader(r io.Reader) (*HtmlTemplate, error) {
// 	t := goHtmlTemplate.New("")

// 	data, err := io.ReadAll(r)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tmpl, err := t.Parse(string(data))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &HtmlTemplate{tmpl}, nil
// }

func NewHtmlTemplate(files ...string) (*HtmlTemplate, error) {
	t, err := goHtmlTemplate.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	t.Funcs(fn)

	return &HtmlTemplate{t}, nil
}

func (t HtmlTemplate) Render(ctx map[string]any) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.Template.Execute(&buf, ctx); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
