package path_test

import (
	"testing"

	"github.com/aziis98/cabret/path"
	"gotest.tools/assert"
)

func Test1(t *testing.T) {
	p, err := path.ParsePattern("/foo/bar")

	assert.NilError(t, err)
	assert.DeepEqual(t, p, &path.PathPattern{
		Parts: []path.PathPart{
			&path.PathSlash{}, &path.PathString{"foo"}, &path.PathSlash{}, &path.PathString{"bar"},
		},
	})
}

func Test2(t *testing.T) {
	p, err := path.ParsePattern("posts/{name}.md")

	assert.NilError(t, err)
	assert.DeepEqual(t, p, &path.PathPattern{
		Parts: []path.PathPart{
			&path.PathString{"posts"}, &path.PathSlash{}, &path.PathShortPattern{"name"}, &path.PathString{".md"},
		},
	})
}

func Test3(t *testing.T) {
	p, err := path.ParsePattern("posts/{{path}}/{name}.md")

	assert.NilError(t, err)
	assert.DeepEqual(t, p, &path.PathPattern{
		Parts: []path.PathPart{
			&path.PathString{"posts"},
			&path.PathSlash{},
			&path.PathLongPattern{"path"},
			&path.PathSlash{},
			&path.PathShortPattern{"name"},
			&path.PathString{".md"},
		},
	})
}

func Test4(t *testing.T) {
	p, err := path.ParsePattern("foo{{path}}/{name}.md")

	assert.NilError(t, err)
	assert.DeepEqual(t, p, &path.PathPattern{
		Parts: []path.PathPart{
			&path.PathString{"foo"},
			&path.PathLongPattern{"path"},
			&path.PathSlash{},
			&path.PathShortPattern{"name"},
			&path.PathString{".md"},
		},
	})

	t.Run("match", func(t *testing.T) {
		m, ctx, err := p.Match("foo-1/bar/baz/post-1.md")

		assert.NilError(t, err, nil)
		assert.Equal(t, m, true)
		assert.DeepEqual(t, ctx, map[string]string{
			"":     "foo-1/bar/baz/post-1.md",
			"name": "post-1",
			"path": "-1/bar/baz",
		})
	})
}
