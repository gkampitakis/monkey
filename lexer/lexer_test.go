package lexer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/gkampitakis/monkey/lexer"
	"github.com/gkampitakis/monkey/token"
)

func TokensSnapshot(t *testing.T, input string) {
	t.Helper()

	l := lexer.New([]byte(input))
	b := strings.Builder{}
	for {
		tok := l.NextToken()
		b.WriteString(fmt.Sprintf("[%s]=>%q", tok.Type, tok.Literal))
		if tok.Type == token.EOF {
			break
		}
		b.WriteByte('\n')
	}

	snaps.MatchSnapshot(t, b.String())
}

func TestNextToken(t *testing.T) {
	t.Run("simple input", func(t *testing.T) {
		input := `=+(){},;`

		TokensSnapshot(t, input)
	})

	t.Run("monkey source code input", func(t *testing.T) {
		input := `let five = 5;

		let ten = 10;

		let add = fn(x, y) {

			x + y;

		};

		let result = add(five, ten);
		`

		TokensSnapshot(t, input)
	})

	t.Run("monkey source code more identifiers", func(t *testing.T) {
		input := `let fine = 5;

		let ten = 10;

		let add = fn(x, y) {

			x+y;

		};

		let result = add(five, ten);

		!-/*5;

		5 < 10 > 5;
		
		if (5 < 10) {
			return true;
		} else {
			return false;
		}

		10 == 10;
		10 != 9;
		let value = 10+10;
		"foobar"
		"foo bar"
		"hello \n world \t"		
		[1, 2];
		{"foo":"bar"}
		`

		TokensSnapshot(t, input)
	})
}
