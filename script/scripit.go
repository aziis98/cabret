package script

import (
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Program struct {
	Statements []Expression `parser:"Newline? ( @@ ( Newline @@ )* Newline? )?"`
}

type ArgumentExpression interface{}

type Expression interface{}

type FunctionReceiver interface{}

var _ FunctionReceiver = &Identifier{}
var _ FunctionReceiver = &ParenthesizedExpression{}

type FunctionCall struct {
	Receiver  FunctionReceiver     `parser:"@@"`
	Arguments []ArgumentExpression `parser:"@@+"`
}

type ParenthesizedExpression struct {
	Inner Expression `parser:"'(' @@ ')'"`
}

type ListExpression struct {
	Values []Expression `parser:"'[' ( Newline? @@ ( Newline? @@ )* Newline? )? ']'"`
}

type DictExpression struct {
	Entries []*DictEntry `parser:"'{' ( Newline? @@ ( ( ',' | Newline ) @@ )* Newline? )? '}'"`
}

type DictEntry struct {
	Key   DictEntryKey `parser:"@@"`
	Value Expression   `parser:"'=' @@"`
}

type DictEntryKey interface{}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

type LiteralBoolean struct {
	Value Boolean `parser:"@('true' | 'false')"`
}

type LiteralString struct {
	Value string `parser:"@String"`
}

type LiteralNumber struct {
	Value float64 `parser:"@Number"`
}

type Identifier struct {
	Value string `parser:"@Ident"`
}

type Quote struct {
	Value Quotation `parser:"'#' @@"`
}

type PropertyAccess struct {
	Receiver IdentifierOrParenthesis   `parser:"@@"`
	Accesses []IdentifierOrParenthesis `parser:"('.' @@)*"`
}

type Quotation interface{}

type ExpressionFactor interface{}

type BinaryOperation struct {
	Lhs      ExpressionFactor `parser:"@@"`
	Operator string           `parser:"( @Operator"`
	Rhs      Expression       `parser:"  @@ )?"`
}

type IdentifierOrParenthesis interface{}

var (
	scriptLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Comment", Pattern: `--[^\n]*\n?`},
		{Name: "String", Pattern: `"(\\"|[^"])*"`},
		{Name: "Number", Pattern: `[-+]?(?:\d*\.)?\d+`},
		{Name: "Ident", Pattern: `[a-zA-Z]\w*`},
		{Name: "Meta", Pattern: `\#`},
		// an operator is any combination of [+-*/%<>=:;,&|^?#@] expect for the string "="
		{Name: "Operator", Pattern: `(=[\+\-\*/%<>=:&|\^?#@]+)|([\+\-\*/%<>:&|\^?#@][\+\-\*/%<>=:&|\^?#@]*)`},
		{Name: "Punct", Pattern: `[\(\)\{\}\[\]=,]`},
		{Name: "Newline", Pattern: `[ \t]*\n\s*`},
		{Name: "Whitespace", Pattern: `[ \t]+`},
	})
	// Parser
	Parser = participle.MustBuild[Program](
		participle.Lexer(scriptLexer),
		participle.Elide("Comment", "Whitespace"),
		participle.Unquote("String"),
		participle.Union[Expression](
			&FunctionCall{},
			&BinaryOperation{},
		),
		participle.Union[ArgumentExpression](
			&BinaryOperation{},
		),
		participle.Union[ExpressionFactor](
			&ListExpression{},
			&DictExpression{},
			&LiteralNumber{},
			&LiteralString{},
			&LiteralBoolean{},
			&ParenthesizedExpression{},
			&Quote{},
			&Identifier{},
		),
		participle.Union[Quotation](
			&Identifier{},
			&ParenthesizedExpression{},
		),
		participle.Union[DictEntryKey](
			&Identifier{},
			&ParenthesizedExpression{},
		),
		participle.Union[FunctionReceiver](
			&Identifier{},
			&ParenthesizedExpression{},
		),
		participle.Union[IdentifierOrParenthesis](
			&Identifier{},
			&ParenthesizedExpression{},
		),
		participle.UseLookahead(2),
	)
)

func Parse(source string) (*Program, error) {
	return Parser.ParseString("", source, participle.Trace(os.Stdout))
	// return Parser.ParseString("", source)
}
