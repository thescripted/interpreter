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
	result, err := expression.evaluate()
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		return
	}
	fmt.Printf("expression result: %#v\n", result)
	// for _, token := range tokens {
	// 	fmt.Printf("token: %v\n", token)
	// }
}

func main() {
	repl()
}
