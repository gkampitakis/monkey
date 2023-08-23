package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/gkampitakis/monkey/lexer"
	"github.com/gkampitakis/monkey/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Bytes()
		if bytes.Equal(line, []byte("exit")) {
			fmt.Fprint(out, "See you next time!\n")
			os.Exit(0)
		}
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "[%s]=>%q\n", tok.Type, tok.Literal)
		}
	}
}
