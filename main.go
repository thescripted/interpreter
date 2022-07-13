package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	repl()
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		run(scanner.Text())
	}
}

func exec(st Statement) {
	st.evaluate()
}

func run(code string) {
	tokenScanner := NewTokenScanner(code)
	tokens := tokenScanner.Scan()
	parser := NewParser(tokens)
	statements := parser.Parse()
	for _, st := range statements {
		val := st.evaluate()
		fmt.Println("return:", val) // not good in a non-REPL setting
	}
}
