package exec

import (
	"log"
	"mime"
	"os"
	gopath "path"

	"github.com/alecthomas/repr"
	"github.com/aziis98/cabret"
	"github.com/aziis98/cabret/config"
	"github.com/aziis98/cabret/operation"
	"github.com/aziis98/cabret/path"
)

type matchResult struct {
	file    string
	context map[string]string
}

func BuildOperation(op config.Operation) cabret.Operation {
	if path, ok := op["read"]; ok {
		return cabret.FlatMapToMapAll{&operation.Target{
			PathTemplate: path.(string),
		}}
	}
	if path, ok := op["write"]; ok {
		return cabret.FlatMapToMapAll{&operation.Target{
			PathTemplate: path.(string),
		}}
	}
	if name, ok := op["plugin"]; ok {
		switch name {
		case "layout":
			path := op["path"].(string)
			delete(op, "path")

			return cabret.FlatMapToMapAll{&operation.Layout{
				TemplateFilesPattern: path,
				Options:              op,
			}}
		case "markdown":
			return cabret.FlatMapToMapAll{&operation.Markdown{
				Options: op,
			}}
		default:
			log.Fatalf(`invalid operation: %s`, name)
		}
	}

	return nil
}

func Execute(cfg *config.Cabretfile) error {
	files, err := cabret.FindFiles([]string{})
	if err != nil {
		return err
	}

	// the first index is the entrypoint ID, the second is for the array of matched files for this entrypoint
	entryPointMatches := make([][]matchResult, len(cfg.EntryPoints))

	// load all files to process
	for id, ep := range cfg.EntryPoints {
		pat, err := path.ParsePattern(ep.Source)
		if err != nil {
			return err
		}

		matchedFiles := []matchResult{}
		for _, f := range files {
			if ok, ctx, _ := pat.Match(f); ok {
				log.Printf(`[Preload] [EntryPoint %d] Found "%s": %#v`, id, f, ctx)

				matchedFiles = append(matchedFiles, matchResult{f, ctx})
			}
		}

		entryPointMatches[id] = matchedFiles
	}

	// TODO: preload all metadata...

	// process all entrypoints
	for id, ep := range cfg.EntryPoints {
		log.Printf(`[EntryPoint %d] starting to process %d file(s)`, id, len(entryPointMatches[id]))

		for _, m := range entryPointMatches[id] {
			log.Printf(`[EntryPoint %d] ["%s"] reading file`, id, m.file)
			data, err := os.ReadFile(m.file)
			if err != nil {
				return err
			}

			content := cabret.Content{
				Type: mime.TypeByExtension(gopath.Ext(m.file)),
				Data: data,
				Metadata: cabret.Map{
					cabret.MatchResult: m.context,
				},
			}

			for i, opConfig := range ep.Pipeline {
				op := BuildOperation(opConfig)

				log.Printf(`[EntryPoint %d] ["%s"] [Operation(%d)] applying %s`, id, m.file, i, repr.String(op))

				newContent, err := op.Process(content)
				if err != nil {
					return err
				}
				if newContent == nil {
					break
				}

				// log.Printf(`[EntryPoint %d] ["%s"] [Operation(%d)] [Metadata] %s`, id, m.file, i, repr.String(newContent.Metadata, repr.Indent("  ")))
				content = *newContent
			}

			log.Printf(`[EntryPoint %d] ["%s"] done`, id, m.file)
		}
	}

	return nil
}
