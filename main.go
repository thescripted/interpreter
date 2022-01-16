package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

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
	ktype   TokenType
	lexeme  string
	literal interface{}
	line    int
}

func (t Token) String() string {
	return fmt.Sprintf("type: %d, lexeme: %v, literal: %v", t.ktype, t.lexeme, t.literal)
}

type Scanner struct {
	source   string
	tokens   []Token
	start    int
	current  int
	line     int
	keywords map[string]TokenType
}

func NewScanner(source string) *Scanner {
	keywords := make(map[string]TokenType)

	keywords["and"] = AND
	keywords["class"] = CLASS
	keywords["else"] = ELSE
	keywords["false"] = FALSE
	keywords["for"] = FOR
	keywords["fun"] = FUN
	keywords["if"] = IF
	keywords["nil"] = NIL
	keywords["or"] = OR
	keywords["print"] = PRINT
	keywords["return"] = RETURN
	keywords["super"] = SUPER
	keywords["this"] = THIS
	keywords["true"] = TRUE
	keywords["var"] = VAR
	keywords["while"] = WHILE

	return &Scanner{
		source:   source,
		tokens:   make([]Token, 0),
		start:    0,
		current:  0,
		line:     1,
		keywords: keywords,
	}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, Token{EOF, "EOF", nil, s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.peek() == '=' {
			s.advance()
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.peek() == '=' {
			s.advance()
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.peek() == '=' {
			s.advance()
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.peek() == '=' {
			s.advance()
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.peek() == '/' { // this is a comment. advance.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case '"':
		{
			for s.peek() != '"' && !s.isAtEnd() {
				if s.peek() == '\n' {
					s.line++
				}
				s.advance()
			}

			if s.isAtEnd() { // we did not terminate the quote
				loxError(s.line, "Unterminated string.")
			}
			s.advance() // consume the closing quote
			value := s.source[s.start+1 : s.current-1]
			s.addTokenWithValue(STRING, value)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
		break
	default:
		if unicode.IsDigit(rune(c)) {
			for unicode.IsDigit(rune(s.peek())) {
				s.advance()
			}
			if s.peek() == '.' && unicode.IsDigit(rune(s.peekNext())) {
				s.advance()
				for unicode.IsDigit(rune(s.peek())) {
					s.advance()
				}
			}
			value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
			if err != nil {
				panic(err) // oh god
			}
			s.addTokenWithValue(NUMBER, value)
		} else if unicode.IsLetter(rune(c)) {
			for unicode.IsLetter(rune(s.peek())) || unicode.IsDigit(rune(s.peek())) {
				s.advance()
			}
			text := s.source[s.start:s.current]
			textType := s.keywords[text]
			if textType == 0 {
				s.addToken(IDENTIFIER)
			} else {
				s.addToken(textType)
			}
		} else {
			loxError(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) addToken(token TokenType) {
	newToken := Token{
		token,
		s.source[s.start:s.current],
		nil,
		s.line,
	}
	s.tokens = append(s.tokens, newToken)
}

func (s *Scanner) addTokenWithValue(token TokenType, value interface{}) {
	newToken := Token{
		token,
		s.source[s.start:s.current],
		value,
		s.line,
	}
	s.tokens = append(s.tokens, newToken)
}

func (s *Scanner) advance() byte {
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
}

func (s *Scanner) isAtEnd() bool {
	if s.current >= len(s.source) {
		return true
	}
	return false
}

var hadError = false

func runFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Bad")
	}
	run(string(data))
	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		run(scanner.Text())
		hadError = false
	}
}

func run(code string) {
	scanner := NewScanner(code)
	tokens := scanner.ScanTokens()
	for _, token := range tokens {
		fmt.Printf("token: %v\n", token)
	}
}

// update this error handling.
func loxError(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
	hadError = true
}

func main() {
	var expression Binary
	expression = Binary{
		left: Unary{
			operator: Token{MINUS, "-", nil, 1},
			right:    Literal{123},
		},
		operator: Token{STAR, "*", nil, 1},
		right: Grouping{
			expression: Unary{
				operator: Token{BANG, "!", nil, 1},
				right:    Literal{"hello"},
			},
		},
	}

	var printerVisitor PrinterVisitor
	printerVisitor = PrinterVisitor{
		output: "",
	}
	expression.accept(&printerVisitor)
	fmt.Println("Output for printerVisitor:", printerVisitor.output)

	//	if len(os.Args) > 2 {
	//		fmt.Println("Usage: glox [script]")
	//		os.Exit(64)
	//	} else if len(os.Args) == 2 {
	//		runFile(os.Args[1])
	//	} else {
	//		runPrompt()
	//	}
}

// =============== Vistor Interface =============== //

type Visitor interface {
	binaryExpr(expr Binary)
	groupingExpr(expr Grouping)
	literalExpr(expr Literal)
	unaryExpr(expr Unary)
}

type PrinterVisitor struct {
	output string
}

func (p *PrinterVisitor) binaryExpr(expr Binary) {
	expr.left.accept(p)
	left := p.output
	expr.right.accept(p)
	right := p.output
	p.output = fmt.Sprintf("(%v %v %v)", expr.operator.lexeme, left, right)
}

func (p *PrinterVisitor) groupingExpr(expr Grouping) {
	expr.expression.accept(p)
	p.output = fmt.Sprintf("(group %v)", p.output)
}

func (p *PrinterVisitor) unaryExpr(expr Unary) {
	expr.right.accept(p)
	p.output = fmt.Sprintf("(%v %v)", expr.operator.lexeme, p.output)
}

func (p *PrinterVisitor) literalExpr(expr Literal) {
	p.output = fmt.Sprint(expr.value)
}

// =============== Expression Interface =============== //
type Expr interface {
	accept(v Visitor)
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

type Unary struct {
	operator Token
	right    Expr
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value interface{}
}

func (this Binary) accept(v Visitor) {
	v.binaryExpr(this)
}

func (this Unary) accept(v Visitor) {
	v.unaryExpr(this)
}

func (this Grouping) accept(v Visitor) {
	v.groupingExpr(this)
}

func (this Literal) accept(v Visitor) {
	v.literalExpr(this)
}
