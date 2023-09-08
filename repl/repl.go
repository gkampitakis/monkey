package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/gkampitakis/monkey/lexer"
	"github.com/gkampitakis/monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	scanner := bufio.NewScanner(in)
	go func() {
		<-c
		fmt.Printf("\n(To exit, press Ctrl+C again or type .exit)\n%s", PROMPT)
		<-c
		fmt.Fprint(out, "\nSee you next time!\n")
		os.Exit(0)
	}()
	for {
		fmt.Fprint(out, PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Bytes()
		if bytes.Equal(line, []byte(".exit")) {
			fmt.Fprint(out, "See you next time!\n")
			os.Exit(0)
		}

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
