package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/gkampitakis/monkey/object"
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
