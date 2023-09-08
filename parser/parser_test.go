package parser_test

import (
	"testing"

	"github.com/gkampitakis/monkey/ast"
	"github.com/gkampitakis/monkey/lexer"
	"github.com/gkampitakis/monkey/parser"
	"github.com/stretchr/testify/require"
)

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	t.Helper()

	require.Equal(t, "let", s.TokenLiteral())
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("expected type let statement but got %T\n", s)
	}
	require.Equal(t, name, letStmt.Name.Value)
	require.Equal(t, name, letStmt.Name.TokenLiteral())
}

func testReturnStatement(t *testing.T, s ast.Statement, name string) {
	t.Helper()

	returnStmt, ok := s.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("expected type let statement but got %T\n", s)
	}
	require.Equal(t, "return", returnStmt.TokenLiteral())
}

func checkParseErrors(t *testing.T, p *parser.Parser, expectedLen int) {
	t.Helper()
	require.Len(t, p.Errors(), expectedLen)
}

func TestParser(t *testing.T) {
	t.Run("let statements", func(t *testing.T) {
		input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
		`

		l := lexer.New([]byte(input))
		p := parser.New(l)
		program := p.ParseProgram()

		require.Len(t, program.Statements, 3)
		identifiers := []string{"x", "y", "foobar"}
		for i, stmt := range program.Statements {
			testLetStatement(t, stmt, identifiers[i])
		}
	})

	t.Run("let statement with error", func(t *testing.T) {
		input := `
		let x 5;
		let = 10;
		let 838383;
		`
		l := lexer.New([]byte(input))
		p := parser.New(l)
		p.ParseProgram()

		checkParseErrors(t, p, 3)
	})

	t.Run("return statement", func(t *testing.T) {
		input := `
			return 5;
			return 10;
			return add(10);
		`

		l := lexer.New([]byte(input))
		p := parser.New(l)
		program := p.ParseProgram()

		checkParseErrors(t, p, 0)
		require.Len(t, program.Statements, 3)

		// for := range program.Statements {}
	})
}
