package repl

import (
	"io"
	"strings"

	"github.com/gkampitakis/monkey/evaluator"
	"github.com/gkampitakis/monkey/lexer"
	"github.com/gkampitakis/monkey/object"
	"github.com/gkampitakis/monkey/parser"
	cli "github.com/openengineer/go-repl"
)

const helpMessage = `help              display this message
.exit              quit this program`

type ReplHandler struct {
	Repl *cli.Repl
	Env  *object.Environment
}

func (r *ReplHandler) Prompt() string {
	return ">> "
}

func (r *ReplHandler) Tab(_ string) string {
	return ""
}

func (r *ReplHandler) Eval(line string) string {
	if line == ".exit" {
		r.Repl.Quit()
		return ""
	}
	l := lexer.New([]byte(line))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return printParserErrors(p.Errors())
	}

	evaluated := evaluator.Eval(program, r.Env)
	if evaluated != nil {
		return evaluated.Inspect()
	}

	return ""
}

func printParserErrors(errors []string) string {
	str := strings.Builder{}

	io.WriteString(&str, "Woops! We ran into some monkey business here!\n")
	io.WriteString(&str, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(&str, "\t"+msg+"\n")
	}

	return str.String()
}
