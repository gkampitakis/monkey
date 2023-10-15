package token

//go:generate stringer -type=TokenType
type TokenType uint8

type Token struct {
	Type    TokenType
	Literal []byte
}

const (
	ILLEGAL TokenType = iota
	EOF
	// Identifiers + literals
	IDENT
	INT
	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH
	LT
	GT
	EQ
	NEQ
	// Delimiters
	COMMA
	SEMICOLON
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	COLON
	STRING_QUOTE
	// Keywords
	FUNCTION
	LET
	RETURN
	IF
	ELSE
	WHILE
	TRUE
	FALSE
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
	"if":     IF,
	"while":  WHILE,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident []byte) TokenType {
	if tok, ok := keywords[string(ident)]; ok {
		return tok
	}
	return IDENT
}
