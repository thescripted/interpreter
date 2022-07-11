package main

import (
	"log"
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
	for !ts.finished() {
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
			if ts.peek() == '"' {
				ts.advance()
			} else { // we didn't close the quotation. Berate the user.
				log.Fatalf("stupid fucking user")
			}
			lexeme := ts.source[ts.start+1 : ts.current+1]
			tokens = append(tokens, Token{
				t:      STRING,
				lexeme: lexeme,
				value:  lexeme,
			})

		case ' ', '\r', '\t', '\n':
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
				value, err := strconv.ParseFloat(ts.sliceToCurrent(), 64)
				if err != nil {
					panic(err) // oh god
				}
				tokens = append(tokens, Token{
					t:      NUMBER,
					lexeme: ts.sliceToCurrent(), // still dont like this
					value:  value,
				})
			} else if unicode.IsLetter(rune(c)) {
				for unicode.IsLetter(rune(ts.peek())) || unicode.IsDigit(rune(ts.peek())) {
					ts.advance()
				}
				text := ts.sliceToCurrent()
				var tokenType TokenType
				if val, ok := tokenKeywords[text]; ok {
					tokenType = val
				}
				// can we do this elsewhere? I don't like this.
				// we're also using append and appendToken a lot. Can this be unified?
				switch tokenType {
				case TRUE:
					tokens = append(tokens, Token{
						t:      TRUE,
						lexeme: ts.sliceToCurrent(),
						value:  true,
					})
				case FALSE:
					tokens = append(tokens, Token{
						t:      TRUE,
						lexeme: ts.sliceToCurrent(),
						value:  false,
					})
				case NIL:
					tokens = append(tokens, Token{
						t:      TRUE,
						lexeme: ts.sliceToCurrent(),
						value:  nil,
					})
				}

			}
		}
		// advance to the next token on complete. I don't like this much either. It's my code. I hate it all.
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

// peekNext lookahead two characters. Stupid.
func (ts *TokenScanner) peekNext() byte {
	if ts.current+2 >= len(ts.source) {
		return '\x00'
	}

	return ts.source[ts.current+2]
}

// finished checks if we've completed scanning. Might not be useful. Only used in one spot.
func (ts *TokenScanner) finished() bool {
	if ts.current >= len(ts.source) {
		return true
	}
	return false
}

// appendToken appends a token to the tokens array. It will not assign a value and its lexeme will be from the start to the current character.
func (ts *TokenScanner) appendToken(tokens []Token, ttype TokenType) []Token { // really should be called "appendTokenWithAssumedLexeme"
	token := Token{
		t:      ttype,
		lexeme: ts.sliceToCurrent(), // I dont like this
	}
	return append(tokens, token)
}
