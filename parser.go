package main

import (
	"log"
)

type Parser struct {
	tokens  []Token
	current int
}

// Parser will parse via recursive descent. Should expand on this (maybe)
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Note: there is no reset. Once you parse it once, you can't parse it again. Is that ok?
func (p *Parser) Parse() []Statement {
	var statements []Statement
	for p.current < len(p.tokens) { // Is this the best way of knmowing if we're at EOF? Should really think this through
		statements = append(statements, p.statement())
		p.advance(1)
	}
	return statements
}

func (p *Parser) statement() Statement {
	var st Statement
	switch p.currentToken().t {
	case VAR:
		p.advance(1)
		if p.currentToken().t != IDENTIFIER {
			panic("something is wrong")
		}
		name := p.currentToken().lexeme
		var expr Expression = nil
		if p.peek().t == EQUAL {
			p.advance(2)
			expr = p.expression()
		}
		st = VariableStatment{name: name, expression: expr}
	case PRINT:
		p.advance(1)
		expr := p.expression()
		st = PrintStatement{expression: expr}
	default:
		expr := p.expression()
		st = ExpressionStatement{expression: expr}
	}

	p.advance(1)                         // what if we can't advance?
	if p.currentToken().t != SEMICOLON { // all statements ends in semicolons.
		panic("something is wrong")
	}

	return st
}

func (p *Parser) expression() Expression {
	return p.equality()
}

func (p *Parser) equality() Expression {
	left := p.comparison()
	op := p.peek()
	for op.t == BANG_EQUAL || op.t == EQUAL_EQUAL {
		operand := op
		p.advance(2)
		op = p.peek()
		right := p.comparison()
		return BinaryExpression{
			left:    left,
			operand: operand,
			right:   right,
		}
	}
	return left
}

func (p *Parser) comparison() Expression {
	left := p.term()
	op := p.peek()
	for op.t == LESS || op.t == LESS_EQUAL || op.t == GREATER || op.t == GREATER_EQUAL {
		operand := op
		p.advance(2)
		op = p.peek()
		right := p.term()
		return BinaryExpression{
			left:    left,
			operand: operand,
			right:   right,
		}
	}
	return left
}

func (p *Parser) term() Expression {
	left := p.factor()
	op := p.peek()
	for op.t == PLUS || op.t == MINUS {
		operand := op
		p.advance(2)
		op = p.peek()
		right := p.factor()
		left = BinaryExpression{
			left:    left,
			operand: operand,
			right:   right,
		}
	}

	return left
}

func (p *Parser) factor() Expression {
	left := p.unary()
	op := p.peek()
	for op.t == STAR || op.t == SLASH {
		operand := op
		p.advance(2)
		op = p.peek()
		right := p.unary()
		return BinaryExpression{
			left:    left,
			operand: operand,
			right:   right,
		}
	}
	return left
}

func (p *Parser) unary() Expression {
	if op := p.currentToken(); op.t == BANG || op.t == MINUS {
		p.advance(1) // can syntax error. Should handle.
		right := p.unary()
		return UnaryExpression{
			operand: op,
			right:   right,
		}
	}
	left := p.primary()
	return left
}

func (p *Parser) primary() Expression {
	literal := p.currentToken()
	switch literal.t {
	case NUMBER, STRING, TRUE, FALSE, NIL:
		return LiteralExpression{value: literal.value}
	case IDENTIFIER:
		return VariableExpression{name: literal.lexeme}
	case LEFT_PAREN:
		p.advance(1)
		expr := p.expression()
		p.advance(1) // this MUST advance to a right paren. Otherwise the user fucked up.
		if p.currentToken().t != RIGHT_PAREN {
			log.Fatal("The user fucked up.")
		}
		return GroupingExpression{
			expression: expr,
		}
	}
	return nil
}

// peek looks at what the next token is. Returns a nil token if there is not a next one.
func (p *Parser) peek() Token {
	if p.current+1 >= len(p.tokens) {
		return Token{} // handle this correclty.
	}
	return p.tokens[p.current+1]
}

// advanceToken moves to the next token
func (p *Parser) advance(n int) {
	p.current = p.current + n // this can break.
}

// currentToken returns the current token
func (p *Parser) currentToken() Token {
	return p.tokens[p.current]
}
