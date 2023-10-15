package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/gkampitakis/monkey/evaluator"
	"github.com/gkampitakis/monkey/lexer"
	"github.com/gkampitakis/monkey/object"
	"github.com/gkampitakis/monkey/parser"
	"github.com/gkampitakis/monkey/repl"
	cli "github.com/openengineer/go-repl"
)

const MONKEY_FACE = `
            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"'' ''"-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '~---~'

`

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		if args[0] != "run" {
			fmt.Println("only supports 'run'")
			os.Exit(1)
		}

		f, err := os.ReadFile(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		l := lexer.New(f)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(p.Errors())
			os.Exit(1)
		}

		evaluated := evaluator.Eval(program, object.NewEnvironment())
		if evaluated != nil && evaluated.Type() != object.NULL {
			fmt.Println(evaluated)
		}
		os.Exit(0)
	} else {
		fmt.Println("should run monkey run <file>")
		return
	}

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the monkey programming language!\n", user.Username)
	fmt.Print(MONKEY_FACE)
	fmt.Printf("Feel free to type in monkey repl or type \"help\" for more info\n")
	h := &repl.ReplHandler{
		Env: object.NewEnvironment(),
	}
	h.Repl = cli.NewRepl(h)

	// start the terminal loop
	if err := h.Repl.Loop(); err != nil {
		log.Fatal(err)
	}
}

func printParserErrors(errors []string) {
	fmt.Println("Woops! We ran into some monkey business here!")
	fmt.Println(" parser errors:")
	for _, msg := range errors {
		fmt.Println("\t" + msg)
	}
}
