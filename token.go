package main

import "fmt"

type TokenType int

const (
	// single-character tokens
	LEFT_PAREN TokenType = iota + 1 //start at one to prevent conflict with null (0) case.
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// one or two character tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// literals
	IDENTIFIER
	STRING
	NUMBER

	// keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

type Token struct {
	t      TokenType
	lexeme string
	value  interface{}
}

func (t Token) String() string {
	return fmt.Sprintf("TokenType: %v, Lexeme: %v, Value: %v", t.t, t.lexeme, t.value)
}
