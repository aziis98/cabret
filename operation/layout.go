package operation

import (
	"fmt"
	goHtmlTemplate "html/template"
	"log"
	"mime"
	"path/filepath"
	"strings"

	"github.com/aziis98/cabret"
	"github.com/aziis98/cabret/operation/template"
	"github.com/aziis98/cabret/util"
)

var HtmlMimeType = mime.TypeByExtension(".html")

func init() {
	registerType("layout", &Layout{})
}

type Layout struct {
	// TemplatePatterns is a list of glob patterns of templates that will be loaded
	TemplatePatterns []string
}

func (op *Layout) Configure(config map[string]any) error {
	if v, ok := config[ShortFormValueKey]; ok {
		globPatternsStr, ok := v.(string)
		if !ok {
			return fmt.Errorf(`expected a comma separated list of glob patterns but got "%v" of type %T`, v, v)
		}

		globPatterns := strings.Split(globPatternsStr, ",")
		for _, pat := range globPatterns {
			op.TemplatePatterns = append(op.TemplatePatterns, strings.TrimSpace(pat))
		}

		return nil
	}
	if v, ok := config["paths"]; ok {
		globPatterns, ok := v.([]string)
		if !ok {
			return fmt.Errorf(`expected a list of glob patterns but got "%v" of type %T`, v, v)
		}

		for _, pat := range globPatterns {
			op.TemplatePatterns = append(op.TemplatePatterns, strings.TrimSpace(pat))
		}

		return nil
	}
	if v, ok := config["path"]; ok {
		globPatternStr, ok := v.(string)
		if !ok {
			return fmt.Errorf(`expected a glob pattern but got "%v" of type %T`, v, v)
		}

		op.TemplatePatterns = []string{strings.TrimSpace(globPatternStr)}

		return nil
	}

	return fmt.Errorf(`invalid config`)
}

func (op Layout) ProcessItem(content cabret.Content) (*cabret.Content, error) {
	// expand glob patterns
	tmplFiles := []string{}
	for _, pat := range op.TemplatePatterns {
		files, err := filepath.Glob(strings.TrimSpace(pat))
		if err != nil {
			return nil, err
		}

		tmplFiles = append(tmplFiles, files...)
	}

	// create template
	tmpl, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		return nil, err
	}

	ctx := util.CloneMap(content.Metadata)

	if content.Type == HtmlMimeType {
		ctx["Content"] = goHtmlTemplate.HTML(content.Data)
	} else {
		ctx["Content"] = content.Data
	}

	log.Printf(`[operation.Layout] rendering into layout "%s"`, strings.Join(op.TemplatePatterns, ", "))

	data, err := tmpl.Render(ctx)
	if err != nil {
		return nil, err
	}

	content.Type = mime.TypeByExtension(filepath.Ext(tmplFiles[0]))
	content.Data = data
	return &content, nil
}
