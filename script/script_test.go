package script_test

import (
	"testing"

	"github.com/aziis98/cabret/script"
	"github.com/bradleyjkemp/cupaloy"
	"gotest.tools/assert"
)

func TestParsing(t *testing.T) {
	t.Run("Number", func(t *testing.T) {
		program, err := script.Parse(`
			-3.14
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.BinaryOperation{Lhs: &script.LiteralNumber{-3.14}},
				},
			},
		)
	})

	t.Run("String", func(t *testing.T) {
		program, err := script.Parse(`
			"Hello, World!"
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.BinaryOperation{Lhs: &script.LiteralString{"Hello, World!"}},
				},
			},
		)
	})

	t.Run("Symbol", func(t *testing.T) {
		program, err := script.Parse(`
			#symbol
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.BinaryOperation{Lhs: &script.Quote{Value: &script.Identifier{"symbol"}}},
				},
			},
		)
	})

	t.Run("BinaryOperation", func(t *testing.T) {
		program, err := script.Parse(`
			1 < 2
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.BinaryOperation{
						Lhs:      &script.LiteralNumber{1},
						Operator: "<",
						Rhs:      &script.BinaryOperation{Lhs: &script.LiteralNumber{2}},
					},
				},
			},
		)
	})

	t.Run("Boolean/True", func(t *testing.T) {
		program, err := script.Parse(`
			true
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.BinaryOperation{Lhs: &script.LiteralBoolean{true}},
				},
			},
		)
	})

	t.Run("Boolean/False", func(t *testing.T) {
		program, err := script.Parse(`
			false
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.BinaryOperation{Lhs: &script.LiteralBoolean{false}},
				},
			},
		)
	})

	t.Run("List", func(t *testing.T) {
		program, err := script.Parse(`
			[1 2 3]
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.BinaryOperation{Lhs: &script.ListExpression{
						Values: []script.Expression{
							&script.BinaryOperation{Lhs: &script.LiteralNumber{1}},
							&script.BinaryOperation{Lhs: &script.LiteralNumber{2}},
							&script.BinaryOperation{Lhs: &script.LiteralNumber{3}},
						},
					}},
				},
			},
		)
	})

	t.Run("List", func(t *testing.T) {
		program, err := script.Parse(`
			[1 2 3 
			
			"a" "b" "c"

			false]
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.BinaryOperation{Lhs: &script.ListExpression{
						Values: []script.Expression{
							&script.BinaryOperation{Lhs: &script.LiteralNumber{1}},
							&script.BinaryOperation{Lhs: &script.LiteralNumber{2}},
							&script.BinaryOperation{Lhs: &script.LiteralNumber{3}},
							&script.BinaryOperation{Lhs: &script.LiteralString{"a"}},
							&script.BinaryOperation{Lhs: &script.LiteralString{"b"}},
							&script.BinaryOperation{Lhs: &script.LiteralString{"c"}},
							&script.BinaryOperation{Lhs: &script.LiteralBoolean{false}},
						},
					}},
				},
			},
		)
	})

	t.Run("Dict", func(t *testing.T) {
		program, err := script.Parse(`
			{}
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.BinaryOperation{Lhs: &script.DictExpression{}},
				},
			},
		)
	})

	t.Run("Dict", func(t *testing.T) {
		program, err := script.Parse(`
			{ a = "1", b = "2" }
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.BinaryOperation{Lhs: &script.DictExpression{
						Entries: []*script.DictEntry{
							{
								Key:   &script.Identifier{"a"},
								Value: &script.BinaryOperation{Lhs: &script.LiteralString{"1"}},
							},
							{
								Key:   &script.Identifier{"b"},
								Value: &script.BinaryOperation{Lhs: &script.LiteralString{"2"}},
							},
						},
					}},
				},
			},
		)
	})

	t.Run("Dict", func(t *testing.T) {
		program, err := script.Parse(`
			{ 
				a = "1"
				b = "2", c = "3"
			}
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.BinaryOperation{Lhs: &script.DictExpression{
						Entries: []*script.DictEntry{
							{
								Key:   &script.Identifier{"a"},
								Value: &script.BinaryOperation{Lhs: &script.LiteralString{"1"}},
							},
							{
								Key:   &script.Identifier{"b"},
								Value: &script.BinaryOperation{Lhs: &script.LiteralString{"2"}},
							},
							{
								Key:   &script.Identifier{"c"},
								Value: &script.BinaryOperation{Lhs: &script.LiteralString{"3"}},
							},
						},
					}},
				},
			},
		)
	})

	t.Run("FunctionCall", func(t *testing.T) {
		program, err := script.Parse(`
			fn 1 2 3
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.FunctionCall{
						Receiver: &script.Identifier{"fn"},
						Arguments: []script.ArgumentExpression{
							&script.BinaryOperation{Lhs: &script.LiteralNumber{1}},
							&script.BinaryOperation{Lhs: &script.LiteralNumber{2}},
							&script.BinaryOperation{Lhs: &script.LiteralNumber{3}},
						},
					},
				},
			},
		)
	})

	t.Run("FunctionCall", func(t *testing.T) {
		program, err := script.Parse(`
			foo { a = 1, b = [1 2 3] } [
				bar #test
				2
			]
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.FunctionCall{
						Receiver: &script.Identifier{"foo"},
						Arguments: []script.ArgumentExpression{
							&script.BinaryOperation{Lhs: &script.DictExpression{
								Entries: []*script.DictEntry{
									{
										Key:   &script.Identifier{"a"},
										Value: &script.BinaryOperation{Lhs: &script.LiteralNumber{1}},
									},
									{
										Key: &script.Identifier{"b"},
										Value: &script.BinaryOperation{Lhs: &script.ListExpression{
											Values: []script.Expression{
												&script.BinaryOperation{Lhs: &script.LiteralNumber{1}},
												&script.BinaryOperation{Lhs: &script.LiteralNumber{2}},
												&script.BinaryOperation{Lhs: &script.LiteralNumber{3}},
											},
										}},
									},
								},
							}},
							&script.BinaryOperation{Lhs: &script.ListExpression{
								Values: []script.Expression{
									&script.FunctionCall{
										Receiver: &script.Identifier{"bar"},
										Arguments: []script.ArgumentExpression{
											&script.BinaryOperation{Lhs: &script.Quote{&script.Identifier{"test"}}},
										},
									},
									&script.BinaryOperation{Lhs: &script.LiteralNumber{2}},
								},
							}},
						},
					},
				},
			},
		)
	})

	t.Run("Operators", func(t *testing.T) {
		program, err := script.Parse(`
			if 1 <= 2 [ println "Yes" ]
		`)

		assert.NilError(t, err)
		assert.DeepEqual(t, program,
			&script.Program{
				Statements: []script.Expression{
					&script.FunctionCall{
						Receiver: &script.Identifier{"if"},
						Arguments: []script.ArgumentExpression{
							&script.BinaryOperation{
								Lhs:      &script.LiteralNumber{1},
								Operator: "<=",
								Rhs: &script.BinaryOperation{
									Lhs: &script.LiteralNumber{2},
								},
							},
							&script.BinaryOperation{Lhs: &script.ListExpression{
								Values: []script.Expression{
									&script.FunctionCall{
										Receiver: &script.Identifier{"println"},
										Arguments: []script.ArgumentExpression{
											&script.BinaryOperation{Lhs: &script.LiteralString{"Yes"}},
										},
									},
								},
							}},
						},
					},
				},
			},
		)
	})

	t.Run("ComplexProgram", func(t *testing.T) {
		program, err := script.Parse(`
			match x [
				#None -> 0
				#(Some value) -> value
			]
		`)

		assert.NilError(t, err)
		cupaloy.SnapshotT(t, program)
	})

	t.Run("ComplexProgram", func(t *testing.T) {
		program, err := script.Parse(`
			for n (range 1 10) [
				m := n * n
				printfln "n^2 = %v" m
			]
		`)

		assert.NilError(t, err)
		cupaloy.SnapshotT(t, program)
	})

	t.Run("ComplexProgram", func(t *testing.T) {
		program, err := script.Parse(`
			foo [ 1 2 3 ] { 
				a = 1
				b = foo #other { bar = 456 }
			} [ 
				4 
				5 
				6
			]
		`)

		assert.NilError(t, err)
		cupaloy.SnapshotT(t, program)
	})
}
