package operation

import (
	"mime"
	"os"
	gopath "path"
	"path/filepath"

	"github.com/aziis98/cabret"
)

var _ cabret.ListOperation = Read{}

type Read struct {
	Patterns []string
}

func (op Read) MapAll(contents []cabret.Content) ([]cabret.Content, error) {
	for _, pattern := range op.Patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			data, err := os.ReadFile(file)
			if err != nil {
				return nil, err
			}

			contents = append(contents, cabret.Content{
				Type:     mime.TypeByExtension(gopath.Ext(file)),
				Data:     data,
				Metadata: cabret.Map{},
			})
		}
	}

	return contents, nil
}
