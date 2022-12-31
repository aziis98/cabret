package operation

import (
	"testing"

	"gotest.tools/assert"
)

func TestBuild(t *testing.T) {
	op, err := Build("categorize", map[string]any{
		"key": "tags",
	})

	assert.NilError(t, err)
	assert.DeepEqual(t, op, &Categorize{
		Key: "tags",
	})
}
