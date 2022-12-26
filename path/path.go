package path

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// FullMatch represents the complete match in a match result context map
const FullMatch = "*"

type PathPart interface {
	regex() string
}

type PathLiteral struct {
	LongAny  bool `parser:"  @'**'"`
	ShortAny bool `parser:"| @'*' "`
	Slash    bool `parser:"| @'/' "`
}

func (p PathLiteral) regex() string {
	switch {
	case p.Slash:
		return "/"
	case p.ShortAny:
		return "([^/]+?)"
	case p.LongAny:
		return "(.+?)"
	default:
		panic("illegal enum state")
	}
}

type PathShortPattern struct {
	Name string `parser:"'{' @Ident '}'"`
}

func (p PathShortPattern) regex() string {
	return fmt.Sprintf("(?P<%s>[^/]+?)", p.Name)
}

type PathLongPattern struct {
	Name string `parser:"'{{' @Ident '}}'"`
}

func (p PathLongPattern) regex() string {
	return fmt.Sprintf("(?P<%s>.+?)", p.Name)
}

type PathString struct {
	Value string `parser:"@Ident"`
}

func (p PathString) regex() string {
	return regexp.QuoteMeta(p.Value)
}

type Pattern struct {
	Parts []PathPart `parser:"@@+"`
}

func (pp Pattern) regex() string {
	sb := &strings.Builder{}

	for _, p := range pp.Parts {
		sb.WriteString(p.regex())
	}

	return sb.String()
}

func (pp Pattern) Match(s string) (bool, map[string]string, error) {
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
		if name != "" {
			ctx[name] = ms[i]
		}
	}

	ctx[FullMatch] = ms[0]

	return true, ctx, nil
}

var parser = participle.MustBuild[Pattern](
	participle.Union[PathPart](
		&PathLiteral{},
		&PathLongPattern{},
		&PathShortPattern{},
		&PathString{},
	),
	participle.Lexer(
		lexer.MustSimple([]lexer.SimpleRule{
			{Name: "Slash", Pattern: `/`},
			{Name: "LongAny", Pattern: `\*\*`},
			{Name: "ShortAny", Pattern: `\*`},
			{Name: "PatternLongOpen", Pattern: `{{`},
			{Name: "PatternLongClose", Pattern: `}}`},
			{Name: "PatternShortOpen", Pattern: `{`},
			{Name: "PatternShortClose", Pattern: `}`},
			{Name: "Ident", Pattern: `[^/{}]+`},
		}),
	),
)

func RenderTemplate(tmpl string, ctx map[string]string) string {
	s := tmpl
	for k, v := range ctx {
		s = strings.ReplaceAll(s, "{"+k+"}", v)
	}
	return s
}

func ParsePattern(path string) (*Pattern, error) {
	return parser.ParseString("", path)
}

func MustParsePattern(path string) *Pattern {
	p, err := ParsePattern(path)
	if err != nil {
		panic(err)
	}

	return p
}
