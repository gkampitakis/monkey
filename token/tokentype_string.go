// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ILLEGAL-0]
	_ = x[EOF-1]
	_ = x[IDENT-2]
	_ = x[INT-3]
	_ = x[ASSIGN-4]
	_ = x[PLUS-5]
	_ = x[MINUS-6]
	_ = x[BANG-7]
	_ = x[ASTERISK-8]
	_ = x[SLASH-9]
	_ = x[LT-10]
	_ = x[GT-11]
	_ = x[EQ-12]
	_ = x[NEQ-13]
	_ = x[COMMA-14]
	_ = x[SEMICOLON-15]
	_ = x[LPAREN-16]
	_ = x[RPAREN-17]
	_ = x[LBRACE-18]
	_ = x[RBRACE-19]
	_ = x[FUNCTION-20]
	_ = x[LET-21]
	_ = x[RETURN-22]
	_ = x[IF-23]
	_ = x[ELSE-24]
	_ = x[TRUE-25]
	_ = x[FALSE-26]
}

const _TokenType_name = "ILLEGALEOFIDENTINTASSIGNPLUSMINUSBANGASTERISKSLASHLTGTEQNEQCOMMASEMICOLONLPARENRPARENLBRACERBRACEFUNCTIONLETRETURNIFELSETRUEFALSE"

var _TokenType_index = [...]uint8{0, 7, 10, 15, 18, 24, 28, 33, 37, 45, 50, 52, 54, 56, 59, 64, 73, 79, 85, 91, 97, 105, 108, 114, 116, 120, 124, 129}

func (i TokenType) String() string {
	if i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}