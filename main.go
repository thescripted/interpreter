package main

import (
	"bufio"
	"fmt"
	"os"
)

var _globalEnvironment = Environment{mem: make(map[string]interface{})} // should place this somewhere else.
func main() {
	if len(os.Args) > 2 {
		sendHelp() // because they need it
	}
	if len(os.Args) == 2 {
		script(os.Args[1])
		return
	}
	repl()
}

// sendHelp prints the default message for usage
func sendHelp() {
	fmt.Print("glox: ./interpret [file]\n")
}

// script will interpret and evaluate the file
func script(file string) {
	dat, err := os.ReadFile(file)
	if err != nil {
		panic(err) // wtf??
	}
	data := string(dat) // we're just slurpring the whole file in memory. Probably not a good idea.
	eval(data)
}

// eval will evaluate source code, statement by statement.
func eval(code string) {
	tokenScanner := NewTokenScanner(code)
	tokens := tokenScanner.Scan()
	parser := NewParser(tokens)
	statements := parser.Parse()
	for _, st := range statements {
		st.eval()
	}
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		eval(scanner.Text())
	}
}
