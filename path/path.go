package path

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type PathPart interface {
	regex() string
}

type PathSlash struct {
	Value string `"/"`
}

func (p PathSlash) regex() string {
	return "/"
}

type PathShortPattern struct {
	Name string `"{" @Ident "}"`
}

func (p PathShortPattern) regex() string {
	return fmt.Sprintf("(?P<%s>[^/]+?)", p.Name)
}

type PathLongPattern struct {
	Name string `"{{" @Ident "}}"`
}

func (p PathLongPattern) regex() string {
	return fmt.Sprintf("(?P<%s>.+?)", p.Name)
}

type PathString struct {
	Value string `@Ident`
}

func (p PathString) regex() string {
	return regexp.QuoteMeta(p.Value)
}

type PathPattern struct {
	Parts []PathPart `@@+`
}

func (pp PathPattern) regex() string {
	sb := &strings.Builder{}

	for _, p := range pp.Parts {
		sb.WriteString(p.regex())
	}

	return sb.String()
}

func (pp PathPattern) Match(s string) (bool, map[string]string, error) {
	r, err := regexp.Compile("^" + pp.regex() + "$")
	if err != nil {
		return false, nil, err
	}

	ms := r.FindStringSubmatch(s)
	if ms == nil {
		return false, nil, nil
	}

	ctx := map[string]string{}
	for i, name := range r.SubexpNames() {
		ctx[name] = ms[i]
	}

	return true, ctx, nil
}

var parser = participle.MustBuild[PathPattern](
	participle.Union[PathPart](
		&PathSlash{},
		&PathLongPattern{},
		&PathShortPattern{},
		&PathString{},
	),
	participle.Lexer(
		lexer.MustSimple([]lexer.SimpleRule{
			{Name: "Slash", Pattern: `/`},
			{Name: "PatternLongOpen", Pattern: `{{`},
			{Name: "PatternLongClose", Pattern: `}}`},
			{Name: "PatternShortOpen", Pattern: `{`},
			{Name: "PatternShortClose", Pattern: `}`},
			{Name: "Ident", Pattern: `[^/{}]+`},
		}),
	),
)

func ParsePattern(path string) (*PathPattern, error) {
	return parser.ParseString("", path)
}

func MustParsePattern(path string) *PathPattern {
	p, err := ParsePattern(path)
	if err != nil {
		panic(err)
	}

	return p
}
