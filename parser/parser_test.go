package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gkampitakis/monkey/ast"
	"github.com/gkampitakis/monkey/lexer"
	"github.com/gkampitakis/monkey/parser"
	"github.com/stretchr/testify/require"
)

/*Start helper methods*/

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	t.Helper()

	require.Equal(t, "let", s.TokenLiteral())
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("expected type let statement but got %T\n", s)
	}
	require.Equal(t, name, string(letStmt.Name.Value))
	require.Equal(t, name, letStmt.Name.TokenLiteral())
}

func testReturnStatement(t *testing.T, s ast.Statement) {
	t.Helper()

	returnStmt, ok := s.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("expected type let statement but got %T\n", s)
	}
	require.Equal(t, "return", returnStmt.TokenLiteral())
}

func assertParseErrors(t *testing.T, p *parser.Parser, expectedLen int) {
	t.Helper()
	if len(p.Errors()) != expectedLen {
		t.Logf(
			"expected len %d but got %d\n%s",
			expectedLen,
			len(p.Errors()),
			strings.Join(p.Errors(), "\n"),
		)
		t.FailNow()
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int) {
	t.Helper()

	require.IsType(t, &ast.IntegerLiteral{}, il)

	integer := il.(*ast.IntegerLiteral)

	require.Equal(t, value, integer.Value)
	require.Equal(t, fmt.Sprintf("%d", integer.Value), integer.TokenLiteral())
}

func testIdentifier(t *testing.T, exp ast.Expression, value []byte) {
	t.Helper()

	require.IsType(t, &ast.Identifier{}, exp)

	ident := exp.(*ast.Identifier)

	require.Equal(t, value, ident.Value)
	require.Equal(t, string(value), ident.TokenLiteral())
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) {
	t.Helper()

	require.IsType(t, &ast.Boolean{}, exp)

	boolean := exp.(*ast.Boolean)

	require.Equal(t, value, boolean.Value)
	require.Equal(t, fmt.Sprintf("%t", value), boolean.TokenLiteral())
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) {
	t.Helper()

	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, exp, v)
	case int64:
		testIntegerLiteral(t, exp, int(v))
	case []byte:
		testIdentifier(t, exp, v)
	case bool:
		testBooleanLiteral(t, exp, v)
	default:
		t.Fatalf("type of exp not handled. got=%T", expected)
	}
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) {
	t.Helper()

	require.IsType(t, &ast.InfixExpression{}, exp)

	opExp := exp.(*ast.InfixExpression)

	testLiteralExpression(t, opExp.Left, left)
	require.Equal(t, operator, opExp.Operator)
	testLiteralExpression(t, opExp.Right, right)
}

/*End helper methods*/

func TestLetStatements(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		tests := []struct {
			input              string
			expectedIdentifier string
			expectedValue      interface{}
		}{
			{"let x = 5;", "x", 5},
			{"let y = true;", "y", true},
			{"let foobar = y;", "foobar", []byte("y")},
		}

		for _, tc := range tests {
			l := lexer.New([]byte(tc.input))
			p := parser.New(l)
			program := p.ParseProgram()
			assertParseErrors(t, p, 0)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain 1 statements. got=%d",
					len(program.Statements))
			}

			stmt := program.Statements[0]
			testLetStatement(t, stmt, tc.expectedIdentifier)

			testLiteralExpression(t, stmt.(*ast.LetStatement).Value, tc.expectedValue)
		}
	})

	t.Run("with errors", func(t *testing.T) {
		input := `
		let x 5;
		let = 10;
		let 838383;
		`
		l := lexer.New([]byte(input))
		p := parser.New(l)
		p.ParseProgram()

		assertParseErrors(t, p, 4)
	})
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", []byte("foobar")},
	}

	for _, tc := range tests {
		l := lexer.New([]byte(tc.input))
		p := parser.New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p, 0)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		testReturnStatement(t, stmt)

		testLiteralExpression(t, stmt.(*ast.ReturnStatement).ReturnValue, tc.expectedValue)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New([]byte(input))
	p := parser.New(l)
	program := p.ParseProgram()

	assertParseErrors(t, p, 0)
	require.Len(t, program.Statements, 1)
	require.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	require.IsType(t, &ast.Identifier{}, stmt.Expression)

	i := stmt.Expression.(*ast.Identifier)
	require.Equal(t, string(i.Value), "foobar")
	require.Equal(t, i.TokenLiteral(), "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New([]byte(input))
	p := parser.New(l)
	program := p.ParseProgram()

	assertParseErrors(t, p, 0)
	require.Len(t, program.Statements, 1)
	require.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	require.IsType(t, &ast.IntegerLiteral{}, stmt.Expression)

	integer := stmt.Expression.(*ast.IntegerLiteral)

	require.Equal(t, integer.Value, 5)
	require.Equal(t, integer.TokenLiteral(), "5")
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", []byte("foobar")},
		{"-foobar;", "-", []byte("foobar")},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tc := range prefixTests {
		l := lexer.New([]byte(tc.input))
		p := parser.New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p, 0)
		require.Len(t, program.Statements, 1)
		require.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		require.IsType(t, &ast.PrefixExpression{}, stmt.Expression)

		exp := stmt.Expression.(*ast.PrefixExpression)

		require.Equal(t, tc.operator, exp.Operator)
		testLiteralExpression(t, exp.Right, tc.value)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 10;", 5, "+", 10},
		{"5 - 10;", 5, "-", 10},
		{"5 * 10;", 5, "*", 10},
		{"5 / 10;", 5, "/", 10},
		{"5 > 10;", 5, ">", 10},
		{"5 < 10;", 5, "<", 10},
		{"5 == 10;", 5, "==", 10},
		{"5 != 10;", 5, "!=", 10},
		{"foobar + barfoo;", []byte("foobar"), "+", []byte("barfoo")},
		{"foobar - barfoo;", []byte("foobar"), "-", []byte("barfoo")},
		{"foobar * barfoo;", []byte("foobar"), "*", []byte("barfoo")},
		{"foobar / barfoo;", []byte("foobar"), "/", []byte("barfoo")},
		{"foobar > barfoo;", []byte("foobar"), ">", []byte("barfoo")},
		{"foobar < barfoo;", []byte("foobar"), "<", []byte("barfoo")},
		{"foobar == barfoo;", []byte("foobar"), "==", []byte("barfoo")},
		{"foobar != barfoo;", []byte("foobar"), "!=", []byte("barfoo")},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tc := range infixTests {
		l := lexer.New([]byte(tc.input))
		p := parser.New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p, 0)

		require.Len(t, program.Statements, 1)
		require.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		testInfixExpression(t, stmt.Expression, tc.leftValue, tc.operator, tc.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tc := range tests {
		l := lexer.New([]byte(tc.input))
		p := parser.New(l)
		program := p.ParseProgram()

		assertParseErrors(t, p, 0)

		require.Equal(t, tc.expected, program.String())
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tc := range tests {
		l := lexer.New([]byte(tc.input))
		p := parser.New(l)
		program := p.ParseProgram()
		assertParseErrors(t, p, 0)

		require.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(
			t,
			ok,
			fmt.Sprintf(
				"program.Statements[0] is not *ast.ExpressionStatement. got=%T",
				program.Statements[0],
			),
		)

		boolean, ok := stmt.Expression.(*ast.Boolean)
		require.True(
			t,
			ok,
			fmt.Sprintf("stmt.Expression is not *ast.Boolean. got=%T", program.Statements[0]),
		)
		testBooleanLiteral(t, boolean, tc.expectedBoolean)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if(x < y) { x }`

	l := lexer.New([]byte(input))
	p := parser.New(l)
	program := p.ParseProgram()
	assertParseErrors(t, p, 0)

	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(
		t,
		ok,
		fmt.Sprintf(
			"expected program.Statements[0] to be type of *ast.ExpressionStatement but got %T",
			program.Statements[0],
		),
	)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t,
		ok,
		fmt.Sprintf("expected stmt.Expression to be type of *ast.IfExpression but got %T", stmt),
	)

	testInfixExpression(t, exp.Condition, []byte("x"), "<", []byte("y"))
	require.Len(t, exp.Consequence.Statements, 1)

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(
		t,
		ok,
		fmt.Sprintf(
			"expected exp.Consequence.Statements[0] to be type of *ast.ExpressionStatement but got %T",
			stmt,
		),
	)
	testIdentifier(t, consequence.Expression, []byte("x"))
	require.Nil(t, exp.Alternative)
}

func TestIfElseExpression(t *testing.T) {
	input := `if(x < y) { x } else { y }`

	l := lexer.New([]byte(input))
	p := parser.New(l)
	program := p.ParseProgram()
	assertParseErrors(t, p, 0)

	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(
		t,
		ok,
		fmt.Sprintf(
			"expected program.Statements[0] to be type of *ast.ExpressionStatement but got %T",
			program.Statements[0],
		),
	)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t,
		ok,
		fmt.Sprintf("expected stmt.Expression to be type of *ast.IfExpression but got %T", stmt),
	)

	testInfixExpression(t, exp.Condition, []byte("x"), "<", []byte("y"))
	require.Len(t, exp.Consequence.Statements, 1)

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(
		t,
		ok,
		fmt.Sprintf(
			"expected exp.Consequence.Statements[0] to be type of *ast.ExpressionStatement but got %T",
			stmt,
		),
	)
	testIdentifier(t, consequence.Expression, []byte("x"))

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	require.True(
		t,
		ok,
		fmt.Sprintf(
			"expected exp.Alternative.Statements[0] to be type of *ast.ExpressionStatement but got %T",
			stmt,
		),
	)
	testIdentifier(t, alternative.Expression, []byte("y"))
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x,y) { x + y; }`
	l := lexer.New([]byte(input))
	p := parser.New(l)
	program := p.ParseProgram()
	assertParseErrors(t, p, 0)

	require.Len(t, program.Statements, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(
		t,
		ok,
		fmt.Sprintf(
			"expected program.Statements[0] to be type of *ast.ExpressionStatement but got %T",
			program.Statements[0],
		),
	)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	require.True(t,
		ok,
		fmt.Sprintf("expected stmt.Expression to be type of *ast.FunctionLiteral but got %T", stmt),
	)
	require.Len(t, function.Parameters, 2)

	testLiteralExpression(t, function.Parameters[0], []byte("x"))
	testLiteralExpression(t, function.Parameters[1], []byte("y"))
	require.Len(t, function.Body.Statements, 1)

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	require.True(
		t,
		ok,
		fmt.Sprintf(
			"expected function.Body.Statements[0] to be type of *ast.BlockStatement but got %T",
			stmt,
		),
	)

	testInfixExpression(t, bodyStmt.Expression, []byte("x"), "+", []byte("y"))
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tc := range tests {
		l := lexer.New([]byte(tc.input))
		p := parser.New(l)
		program := p.ParseProgram()
		assertParseErrors(t, p, 0)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		require.Len(t, function.Parameters, len(tc.expectedParams))

		for i, ident := range tc.expectedParams {
			testLiteralExpression(t, function.Parameters[i], []byte(ident))
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1,2*3,4+5);`
	l := lexer.New([]byte(input))
	p := parser.New(l)
	program := p.ParseProgram()
	assertParseErrors(t, p, 0)

	require.Len(t, program.Statements, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(
		t,
		ok,
		fmt.Sprintf(
			"expected program.Statements[0] to be type of *ast.ExpressionStatement but got %T",
			program.Statements[0],
		),
	)
	exp, ok := stmt.Expression.(*ast.CallExpression)
	require.True(t,
		ok,
		fmt.Sprintf("expected stmt.Expression to be type of *ast.CallExpression but got %T", stmt),
	)

	testIdentifier(t, exp.Function, []byte("add"))
	require.Len(t, exp.Arguments, 3)
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tc := range tests {
		l := lexer.New([]byte(tc.input))
		p := parser.New(l)
		program := p.ParseProgram()
		assertParseErrors(t, p, 0)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp := stmt.Expression.(*ast.CallExpression)

		testIdentifier(t, exp.Function, []byte(tc.expectedIdent))
		require.Len(t, exp.Arguments, len(tc.expectedArgs))

		for i, arg := range tc.expectedArgs {
			require.Equal(t, arg, exp.Arguments[i].String())
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`
	l := lexer.New([]byte(input))
	p := parser.New(l)
	program := p.ParseProgram()
	assertParseErrors(t, p, 0)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal := stmt.Expression.(*ast.StringLiteral)

	require.Equal(t, "hello world", string(literal.Value))
}
