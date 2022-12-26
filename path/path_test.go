package path_test

import (
	"testing"

	"github.com/aziis98/cabret/path"
	"gotest.tools/assert"
)

func TestPath(t *testing.T) {
	p, err := path.ParsePattern("/foo/bar")

	assert.NilError(t, err)
	assert.DeepEqual(t, p, &path.Pattern{
		Parts: []path.PathPart{
			&path.PathLiteral{Slash: true},
			&path.PathString{"foo"},
			&path.PathLiteral{Slash: true},
			&path.PathString{"bar"},
		},
	})
}

func TestSimplePattern(t *testing.T) {
	p, err := path.ParsePattern("posts/{name}.md")

	assert.NilError(t, err)
	assert.DeepEqual(t, p, &path.Pattern{
		Parts: []path.PathPart{
			&path.PathString{"posts"},
			&path.PathLiteral{Slash: true},
			&path.PathShortPattern{"name"},
			&path.PathString{".md"},
		},
	})
}

func TestSimplePatternWithAny(t *testing.T) {
	p, err := path.ParsePattern("posts/*/{name}.md")

	assert.NilError(t, err)
	assert.DeepEqual(t, p, &path.Pattern{
		Parts: []path.PathPart{
			&path.PathString{"posts"},
			&path.PathLiteral{Slash: true},
			&path.PathLiteral{ShortAny: true},
			&path.PathLiteral{Slash: true},
			&path.PathShortPattern{"name"},
			&path.PathString{".md"},
		},
	})

	t.Run("MatchGood", func(t *testing.T) {
		m, ctx, err := p.Match("posts/2022/post-1.md")

		assert.NilError(t, err, nil)
		assert.Equal(t, m, true)
		assert.DeepEqual(t, ctx, map[string]string{
			path.FullMatch: "posts/2022/post-1.md",
			"name":         "post-1",
		})
	})

	t.Run("MatchBad", func(t *testing.T) {
		m, ctx, err := p.Match("posts/2022/12/post-1.md")

		assert.NilError(t, err)
		assert.Assert(t, m == false)
		assert.Assert(t, ctx == nil)
	})
}

func TestSimplePatternCombo(t *testing.T) {
	p, err := path.ParsePattern("posts/{{path}}/{name}.md")

	assert.NilError(t, err)
	assert.DeepEqual(t, p, &path.Pattern{
		Parts: []path.PathPart{
			&path.PathString{"posts"},
			&path.PathLiteral{Slash: true},
			&path.PathLongPattern{"path"},
			&path.PathLiteral{Slash: true},
			&path.PathShortPattern{"name"},
			&path.PathString{".md"},
		},
	})
}

func TestPatternCombo(t *testing.T) {
	p, err := path.ParsePattern("posts/**/{date}_{slug}.md")

	assert.NilError(t, err)
	assert.DeepEqual(t, p, &path.Pattern{
		Parts: []path.PathPart{
			&path.PathString{"posts"},
			&path.PathLiteral{Slash: true},
			&path.PathLiteral{LongAny: true},
			&path.PathLiteral{Slash: true},
			&path.PathShortPattern{"date"},
			&path.PathString{"_"},
			&path.PathShortPattern{"slug"},
			&path.PathString{".md"},
		},
	})

	t.Run("MatchGood", func(t *testing.T) {
		ok, ctx, err := p.Match("posts/a/b/c/2022-12-25_example.md")

		assert.Assert(t, ok)
		assert.NilError(t, err)
		assert.DeepEqual(t, ctx, map[string]string{
			path.FullMatch: "posts/a/b/c/2022-12-25_example.md",
			"date":         "2022-12-25",
			"slug":         "example",
		})
	})
}

func TestStrangePatternCombo(t *testing.T) {
	p, err := path.ParsePattern("foo{{path}}/{name}.md")

	assert.NilError(t, err)
	assert.DeepEqual(t, p, &path.Pattern{
		Parts: []path.PathPart{
			&path.PathString{"foo"},
			&path.PathLongPattern{"path"},
			&path.PathLiteral{Slash: true},
			&path.PathShortPattern{"name"},
			&path.PathString{".md"},
		},
	})

	t.Run("match", func(t *testing.T) {
		m, ctx, err := p.Match("foo-1/bar/baz/post-1.md")

		assert.NilError(t, err, nil)
		assert.Equal(t, m, true)
		assert.DeepEqual(t, ctx, map[string]string{
			path.FullMatch: "foo-1/bar/baz/post-1.md",
			"name":         "post-1",
			"path":         "-1/bar/baz",
		})
	})
}

func TestMultipleGroups(t *testing.T) {
	p, err := path.ParsePattern("{a}/{b}/{c}")

	assert.NilError(t, err)
	assert.DeepEqual(t, p, &path.Pattern{
		Parts: []path.PathPart{
			&path.PathShortPattern{"a"},
			&path.PathLiteral{Slash: true},
			&path.PathShortPattern{"b"},
			&path.PathLiteral{Slash: true},
			&path.PathShortPattern{"c"},
		},
	})

	t.Run("match", func(t *testing.T) {
		m, ctx, err := p.Match("foo/bar/baz")

		assert.NilError(t, err, nil)
		assert.Equal(t, m, true)
		assert.DeepEqual(t, ctx, map[string]string{
			path.FullMatch: "foo/bar/baz",
			"a":            "foo",
			"b":            "bar",
			"c":            "baz",
		})
	})
}
