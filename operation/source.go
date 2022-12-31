package operation

import (
	"fmt"
	"log"
	"mime"
	"os"
	gopath "path"

	"github.com/aziis98/cabret"
	"github.com/aziis98/cabret/path"
)

func init() {
	registerType("source", &Source{})
}

// Source is a ListOperation that appends the matched files to the processing items
type Source struct {
	Patterns []string
}

func (op *Source) Configure(config map[string]any) error {
	if v, ok := config["source"]; ok {
		pattern, ok := v.(string)
		if !ok {
			return fmt.Errorf(`expected a path pattern but got "%v" of type %T`, v, v)
		}

		op.Patterns = []string{pattern}
		return nil
	}
	if v, ok := config["paths"]; ok {
		aPatterns, ok := v.([]any)
		if !ok {
			return fmt.Errorf(`expected a list of path patterns but got "%v" of type %T`, v, v)
		}

		patterns := []string{}
		for _, aPat := range aPatterns {
			p, ok := aPat.(string)
			if !ok {
				return fmt.Errorf(`expected a string but got "%v" of type %T`, aPat, aPat)
			}

			patterns = append(patterns, p)
		}

		op.Patterns = patterns
		return nil
	}

	return fmt.Errorf(`invalid config for "source": %#v`, config)
}

func (op Source) ProcessList(contents []cabret.Content) ([]cabret.Content, error) {
	files, err := cabret.FindFiles([]string{})
	if err != nil {
		return nil, err
	}

	matches := []cabret.MatchResult{}

	for _, patternStr := range op.Patterns {
		pat, err := path.ParsePattern(patternStr)
		if err != nil {
			return nil, err
		}

		for _, f := range files {
			if ok, captures, _ := pat.Match(f); ok {
				matches = append(matches, cabret.MatchResult{
					File:     f,
					Captures: captures,
				})
			}
		}
	}

	for _, m := range matches {
		log.Printf(`[operation.Source] reading "%s"`, m.File)

		data, err := os.ReadFile(m.File)
		if err != nil {
			return nil, err
		}

		contents = append(contents, cabret.Content{
			Type: mime.TypeByExtension(gopath.Ext(m.File)),
			Data: data,
			Metadata: cabret.Map{
				cabret.MatchResultKey: m.Captures,
			},
		})
	}

	return contents, nil
}
