package operation

// import (
// 	"fmt"
// 	"log"
// 	"mime"
// 	"path/filepath"
// 	"strings"

// 	"github.com/aziis98/cabret"
// 	"github.com/aziis98/cabret/operation/template"
// )

// func init() {
// 	registerType("template", &Template{})
// }

// type Template struct {
// 	// TemplatePatterns is a list of glob patterns of templates that will be loaded
// 	TemplatePatterns []string
// }

// func (op *Template) Load(config map[string]any) error {
// 	if v, ok := config[ShortFormValueKey]; ok {
// 		globPatternsStr, ok := v.(string)
// 		if !ok {
// 			return fmt.Errorf(`expected a comma separated list of glob patterns but got "%v" of type %T`, v, v)
// 		}

// 		globPatterns := strings.Split(globPatternsStr, ",")
// 		for _, pat := range globPatterns {
// 			op.TemplatePatterns = append(op.TemplatePatterns, strings.TrimSpace(pat))
// 		}

// 		return nil
// 	}
// 	if v, ok := config["paths"]; ok {
// 		globPatterns, ok := v.([]string)
// 		if !ok {
// 			return fmt.Errorf(`expected a list of glob patterns but got "%v" of type %T`, v, v)
// 		}

// 		for _, pat := range globPatterns {
// 			op.TemplatePatterns = append(op.TemplatePatterns, strings.TrimSpace(pat))
// 		}

// 		return nil
// 	}
// 	if v, ok := config["path"]; ok {
// 		globPatternStr, ok := v.(string)
// 		if !ok {
// 			return fmt.Errorf(`expected a glob pattern but got "%v" of type %T`, v, v)
// 		}

// 		op.TemplatePatterns = []string{strings.TrimSpace(globPatternStr)}

// 		return nil
// 	}

// 	return fmt.Errorf(`invalid config for "template": %#v`, config)
// }

// func (op *Template) ProcessList(contents []cabret.Content) ([]cabret.Content, error) {
// 	// expand glob patterns
// 	tmplFiles := []string{}
// 	for _, pat := range op.TemplatePatterns {
// 		files, err := filepath.Glob(strings.TrimSpace(pat))
// 		if err != nil {
// 			return nil, err
// 		}

// 		tmplFiles = append(tmplFiles, files...)
// 	}

// 	// create template
// 	tmpl, err := template.ParseFiles(tmplFiles...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Printf(`[operation.Layout] rendering into layout "%s"`, strings.Join(op.TemplatePatterns, ", "))

// 	ctx := map[string]any{}
// 	ctx["Items"] = contents
// 	data, err := tmpl.Render(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return []cabret.Content{
// 		{
// 			Type:     mime.TypeByExtension(filepath.Ext(tmplFiles[0])),
// 			Data:     data,
// 			Metadata: ctx,
// 		},
// 	}, nil
// }
