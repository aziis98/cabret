package operation

import (
	"html/template"
	"log"
	"mime"
	"path/filepath"
	"strings"

	gopath "path"

	"github.com/alecthomas/repr"
	"github.com/aziis98/cabret"
	"github.com/aziis98/cabret/operation/layout"
	"github.com/aziis98/cabret/util"
)

var HtmlMimeType = mime.TypeByExtension(".html")

var _ cabret.Operation = Layout{}

type Layout struct {
	// TemplateFilesPattern is a comma separated list of unix glob patterns
	TemplateFilesPattern string

	Options map[string]any
}

func (op Layout) Process(content cabret.Content) (*cabret.Content, error) {
	var tmpl layout.Template

	patterns := strings.Split(op.TemplateFilesPattern, ",")

	tmplFiles := []string{}
	for _, pat := range patterns {
		files, err := filepath.Glob(strings.TrimSpace(pat))
		if err != nil {
			return nil, err
		}

		tmplFiles = append(tmplFiles, files...)
	}

	log.Printf(`[Layout] template pattern "%s" expanded to %s`, op.TemplateFilesPattern, repr.String(tmplFiles))

	if gopath.Ext(tmplFiles[0]) == ".html" {
		var err error
		if tmpl, err = layout.NewHtmlTemplate(tmplFiles...); err != nil {
			return nil, err
		}
	} else {
		var err error
		if tmpl, err = layout.NewTextTemplate(tmplFiles...); err != nil {
			return nil, err
		}
	}

	ctx := util.CloneMap(content.Metadata)

	if content.Type == HtmlMimeType {
		ctx["Content"] = template.HTML(content.Data)
	} else {
		ctx["Content"] = content.Data
	}

	data, err := tmpl.Render(ctx)
	if err != nil {
		return nil, err
	}

	content.Data = data
	return &content, nil
}
