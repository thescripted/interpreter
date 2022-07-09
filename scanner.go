package main

import (
	"strconv"
	"unicode"
)

type TokenScanner struct {
	source  string
	start   int
	current int
}

var tokenKeywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func NewTokenScanner(source string) *TokenScanner {
	return &TokenScanner{
		source:  source,
		start:   0,
		current: 0,
	}
}

// Scan takes in source code and emits tokens
func (ts *TokenScanner) Scan() []Token {
	var tokens []Token
	for ts.current < len(ts.source) {
		ts.start = ts.current
		c := ts.source[ts.current]
		switch c {
		case '(':
			tokens = ts.appendToken(tokens, LEFT_PAREN)
		case ')':
			tokens = ts.appendToken(tokens, RIGHT_PAREN)
		case '{':
			tokens = ts.appendToken(tokens, LEFT_BRACE)
		case '}':
			tokens = ts.appendToken(tokens, RIGHT_BRACE)
		case ',':
			tokens = ts.appendToken(tokens, COMMA)
		case '.':
			tokens = ts.appendToken(tokens, DOT)
		case '-':
			tokens = ts.appendToken(tokens, MINUS)
		case '+':
			tokens = ts.appendToken(tokens, PLUS)
		case ';':
			tokens = ts.appendToken(tokens, SEMICOLON)
		case '*':
			tokens = ts.appendToken(tokens, STAR)
		case '!':
			if ts.peek() == '=' {
				ts.advance()
				tokens = ts.appendToken(tokens, BANG_EQUAL)
			} else {
				tokens = ts.appendToken(tokens, BANG)
			}
		case '=':
			if ts.peek() == '=' {
				ts.advance()
				tokens = ts.appendToken(tokens, EQUAL_EQUAL)
			} else {
				tokens = ts.appendToken(tokens, EQUAL)
			}
		case '<':
			if ts.peek() == '=' {
				ts.advance()
				tokens = ts.appendToken(tokens, LESS_EQUAL)
			} else {
				tokens = ts.appendToken(tokens, LESS)
			}
		case '>':
			if ts.peek() == '=' {
				ts.advance()
				tokens = ts.appendToken(tokens, GREATER_EQUAL)
			} else {
				tokens = ts.appendToken(tokens, GREATER)
			}
		case '/':
			if ts.peek() == '/' { // comment
				ts.advance()
				for ts.current < len(ts.source) && ts.peek() != '\n' {
					ts.advance()
				}
			} else {
				tokens = ts.appendToken(tokens, SLASH)
			}
		case '"':
			for ts.current < len(ts.source) && ts.peek() != '"' {
				ts.advance()
			}
			tokens = append(tokens, Token{
				t:      STRING,
				lexeme: ts.source[ts.start+1 : ts.current],
			})
		case ' ':
		case '\r':
		case '\t':
		case '\n':
			break
		default:
			if unicode.IsDigit(rune(c)) {
				for unicode.IsDigit(rune(ts.peek())) {
					ts.advance()
				}
				if ts.peek() == '.' && unicode.IsDigit(rune(ts.peekNext())) {
					ts.advance()
					for unicode.IsDigit(rune(ts.peek())) {
						ts.advance()
					}
				}
				_, err := strconv.ParseFloat(ts.sliceToCurrent(), 64)
				if err != nil {
					panic(err) // oh god
				}
				tokens = ts.appendToken(tokens, NUMBER)
			} else if unicode.IsLetter(rune(c)) {
				for unicode.IsLetter(rune(ts.peek())) || unicode.IsDigit(rune(ts.peek())) {
					ts.advance()
				}
				text := ts.sliceToCurrent()
				var tokenType TokenType
				if val, ok := tokenKeywords[text]; ok {
					tokenType = val
				}
				tokenType = IDENTIFIER
				tokens = ts.appendToken(tokens, tokenType)

			}
		}

		// advance to the next token on complete
		ts.advance()
	}
	return tokens
}

// advance increments the current pointer by 1
func (ts *TokenScanner) advance() {
	ts.current++
}

// sliceToCurrent slices the source code from start to current
func (ts *TokenScanner) sliceToCurrent() string { // TODO(ben): better word for this
	return ts.source[ts.start : ts.current+1]
}

// peek lookahead one character
func (ts *TokenScanner) peek() byte {
	if ts.current+1 >= len(ts.source) {
		return '\x00'
	}

	return ts.source[ts.current+1]
}

// peekNext lookahead two characters
func (ts *TokenScanner) peekNext() byte {
	if ts.current+2 >= len(ts.source) {
		return '\x00'
	}

	return ts.source[ts.current+2]
}

func (ts *TokenScanner) finished() bool {
	if ts.current >= len(ts.source) {
		return true
	}
	return false
}

func (ts *TokenScanner) appendToken(tokens []Token, ttype TokenType) []Token { // really should be called "appendTokenWithAssumedLexeme"
	token := Token{
		t:      ttype,
		lexeme: ts.sliceToCurrent(),
	}
	return append(tokens, token)
}
