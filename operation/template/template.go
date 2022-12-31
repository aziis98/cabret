package template

import (
	"path/filepath"
)

type Template interface {
	Render(ctx map[string]any) ([]byte, error)
}

func ParseFiles(files ...string) (Template, error) {
	if filepath.Ext(files[0]) == ".html" {
		return NewHtmlTemplate(files...)
	}

	return NewTextTemplate(files...)
}
