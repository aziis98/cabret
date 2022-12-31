package operation_test

import (
	"testing"

	"github.com/aziis98/cabret/operation"
	"github.com/aziis98/cabret/util"
	"gotest.tools/assert"
)

func TestTargetShortForm(t *testing.T) {
	t.Run("correct usage", func(t *testing.T) {
		op := &operation.Target{}
		err := op.Configure(util.ParseYAML(`
			target: foo/bar/baz.txt
		`))

		assert.NilError(t, err)
		assert.DeepEqual(t, op, &operation.Target{
			PathTemplate: "foo/bar/baz.txt",
		})
	})
	t.Run("wrong usage", func(t *testing.T) {
		op := &operation.Target{}
		err := op.Configure(util.ParseYAML(`
			target: 123
		`))

		assert.Error(t, err, `expected a path template but got "123" of type int`)
	})
}
