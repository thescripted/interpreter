package main

import (
	"bufio"
	"fmt"
	"os"
)

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		run(scanner.Text())
	}
}

func run(code string) {
	tokenScanner := NewTokenScanner(code)
	tokens := tokenScanner.Scan()
	parser := NewParser(tokens)
	expression := parser.Parse()
	fmt.Printf("expression: %v\n", expression)
	// for _, token := range tokens {
	// 	fmt.Printf("token: %v\n", token)
	// }
}

func main() {
	repl()
}
