package evaluator_test

import (
	"testing"

	"github.com/gkampitakis/monkey/evaluator"
	"github.com/gkampitakis/monkey/lexer"
	"github.com/gkampitakis/monkey/object"
	"github.com/gkampitakis/monkey/parser"
	"github.com/stretchr/testify/require"
)

/*Start helper methods*/

func testEval(input string) object.Object {
	l := lexer.New([]byte(input))
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return evaluator.Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int) {
	t.Helper()

	result := obj.(*object.Integer)
	require.Equal(t, expected, result.Value)
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) {
	t.Helper()

	result := obj.(*object.Boolean)
	require.Equal(t, expected, result.Value)
}

func testNullObject(t *testing.T, obj object.Object) {
	t.Helper()
	require.IsType(t, &object.Null{}, obj)
}

/*End helper methods*/

func TestEvalIntegerExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tc := range tests {
		testIntegerObject(t, testEval(tc.input), tc.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tc := range tests {
		testBooleanObject(t, testEval(tc.input), tc.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tc := range tests {
		testBooleanObject(t, testEval(tc.input), tc.expected)
	}
}

func TestIfExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		if integer, ok := tc.expected.(int); ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"return 10;", 10},
		{"return 10;9;", 10},
		{"return 2 * 5;9;", 10},
		{"9;return 2 * 5;9;", 10},
		{`if (10 >1) {
			if (10 > 1) {
				return 10;
			}
			return 1;
		}`, 10},
		{
			`
let f = fn(x) {
  return x;
  x + 10;
};
f(10);`,
			10,
		},
		{
			`
let f = fn(x) {
   let result = x + 10;
   return result;
   return 10;
};
f(10);`,
			20,
		},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input       string
		expectedMsg string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"true + false + true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World!`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)

		require.IsType(t, &object.ErrorValue{}, evaluated)
		errObj := evaluated.(*object.ErrorValue)

		require.Equal(t, tc.expectedMsg, errObj.Message)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let a = 5; a;", 5},
		{"let a = 5;let b = a;b;", 5},
		{"let a = 5;let b = a;let c = a + b + 5;c;", 15},
	}

	for _, tc := range tests {
		testIntegerObject(t, testEval(tc.input), tc.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluate := testEval(input)

	require.IsType(t, &object.Function{}, evaluate)
	fn := evaluate.(*object.Function)

	require.Len(t, fn.Parameters, 1)
	require.Equal(t, "x", fn.Parameters[0].String())
	require.Equal(t, "(x + 2)", fn.Body.String())
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		// {"let id = fn(x) {x;} id(5);", 5},
		// {"let id = fn(x) { return x;} id(5);", 5},
		// {"let double = fn(x) { x*2;} double(5);", 10},
		// {"let add = fn(x,y) { x+y;} add(5, 5);", 10},
		// {"let add = fn(x,y) { x+y;} add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tc := range tests {
		testIntegerObject(t, testEval(tc.input), tc.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
		let newAdder = fn(x) {
			fn(y) { x + y;};
		};

		let addTwo = newAdder(2);
		addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		input := `"Hello World!`

		evaluate := testEval(input)

		require.IsType(t, &object.String{}, evaluate)
		str := evaluate.(*object.String)

		require.Equal(t, "Hello World!", str.Value)
	})

	t.Run("concatenation", func(t *testing.T) {
		input := `"Hello" + " " + "World!"`

		evaluate := testEval(input)

		require.IsType(t, &object.String{}, evaluate)
		str := evaluate.(*object.String)

		require.Equal(t, "Hello World!", str.Value)
	})

	t.Run("equality", func(t *testing.T) {
		input := ` let tmp = "hello";
		let tmptwo = "hello";
		tmp == tmptwo;
		`

		evaluate := testEval(input)

		require.IsType(t, &object.Boolean{}, evaluate)
		boolean := evaluate.(*object.Boolean)

		require.True(t, boolean.Value)
	})

	t.Run("inequality", func(t *testing.T) {
		input := ` let tmp = "hello";
		let tmptwo = "world";
		tmp != tmptwo;
		`

		evaluate := testEval(input)

		require.IsType(t, &object.Boolean{}, evaluate)
		boolean := evaluate.(*object.Boolean)

		require.True(t, boolean.Value)
	})
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		// {`puts("hello", "world!")`, nil},
		{`first([1, 2, 3])`, 1},
		{`first([])`, nil},
		{`first(1)`, "argument to `first` must be ARRAY, got INTEGER"},
		{`last([1, 2, 3])`, 3},
		{`last([])`, nil},
		{`last(1)`, "argument to `last` must be ARRAY, got INTEGER"},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, nil},
		{`push(1, 1)`, "argument to `push` must be ARRAY, got INTEGER"},
		{`push([], 1)`, []int{1}},
		{`push([], 1,2,3,4,5)`, []int{1, 2, 3, 4, 5}},
	}

	for _, tc := range tests {
		evaluate := testEval(tc.input)
		switch expected := tc.expected.(type) {
		case nil:
			testNullObject(t, evaluate)
		case []int:
			array := evaluate.(*object.Array)

			require.Len(t, array.Elements, len(expected))
			for i, expectedElem := range expected {
				testIntegerObject(t, array.Elements[i], expectedElem)
			}
		case int:
			testIntegerObject(t, evaluate, expected)
		case string:
			errObject := evaluate.(*object.ErrorValue)

			require.Equal(t, tc.expected, errObject.Message)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := `[1, 2*2, 3+3]`
	evaluate := testEval(input)

	require.IsType(t, &object.Array{}, evaluate)
	arrayLiteral := evaluate.(*object.Array)

	require.Len(t, arrayLiteral.Elements, 3)
	testIntegerObject(t, arrayLiteral.Elements[0], 1)
	testIntegerObject(t, arrayLiteral.Elements[1], 4)
	testIntegerObject(t, arrayLiteral.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i]",
			1,
		},
		{
			"[1, 2, 3][1+1]",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0]+myArray[1]+myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tc := range tests {
		evaluate := testEval(tc.input)
		if integer, ok := tc.expected.(int); ok {
			testIntegerObject(t, evaluate, integer)
			return
		}

		testNullObject(t, evaluate)
	}
}
