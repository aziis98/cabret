package operation_test

import (
	"testing"

	"github.com/aziis98/cabret/operation"
	"github.com/aziis98/cabret/util"
	"gotest.tools/assert"
)

func TestSourceShortForm(t *testing.T) {
	t.Run("correct usage", func(t *testing.T) {
		op := &operation.Source{}
		err := op.Configure(util.ParseYAML(`
			source: foo/bar/baz.txt
		`))

		assert.NilError(t, err)
		assert.DeepEqual(t, op, &operation.Source{
			Patterns: []string{
				"foo/bar/baz.txt",
			},
		})
	})
	t.Run("wrong usage", func(t *testing.T) {
		op := &operation.Source{}
		err := op.Configure(util.ParseYAML(`
			source: 123
		`))

		assert.Error(t, err, `expected a path pattern but got "123" of type int`)
	})
}

func TestSourceWithPaths(t *testing.T) {
	t.Run("correct usage", func(t *testing.T) {
		op := &operation.Source{}
		err := op.Configure(util.ParseYAML(`
			use: source
			paths: 
			- foo/bar/baz-1.txt
			- foo/bar/baz-2.txt
			- foo/bar/baz-3.txt
		`))

		assert.NilError(t, err)
		assert.DeepEqual(t, op, &operation.Source{
			Patterns: []string{
				"foo/bar/baz-1.txt",
				"foo/bar/baz-2.txt",
				"foo/bar/baz-3.txt",
			},
		})
	})

	t.Run("wrong usage", func(t *testing.T) {
		op := &operation.Source{}
		err := op.Configure(util.ParseYAML(`
			use: source
			paths: 
			- foo/bar/baz-1.txt
			- foo/bar/baz-2.txt
			- 123
		`))

		assert.Error(t, err, `expected a string but got "123" of type int`)
	})
}
