package cabret

import (
	"io/fs"
	"path/filepath"

	"github.com/aziis98/cabret/path"
	"golang.org/x/exp/slices"
)

func FindFiles(excludes []string) ([]string, error) {
	paths := []string{}

	excludeMatchers := []*path.Pattern{}
	for _, p := range excludes {
		m, err := path.ParsePattern(p)
		if err != nil {
			return nil, err
		}

		excludeMatchers = append(excludeMatchers, m)
	}

	if err := filepath.Walk(".", func(p string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		excluded := slices.ContainsFunc(excludeMatchers, func(excludePattern *path.Pattern) bool {
			ok, _, _ := excludePattern.Match(p)
			return ok
		})

		if excluded {
			return nil
		}

		paths = append(paths, p)
		return nil
	}); err != nil {
		return nil, err
	}

	return paths, nil
}
