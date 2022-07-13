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
		p.advance(1) // Someone needs to move to the next statement. Very very imperative, I dont' like it.
	}
	return statements
}

func (p *Parser) statement() Statement {
	if token := p.currentToken(); token.t == PRINT {
		p.advance(1)
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) expressionStatement() Statement {
	value := p.expression()
	if p.peek().t != SEMICOLON {
		panic("you messed up")
	}
	p.advance(1) // someone needs to consume the semicolon. Should it be here or higher up (or a different implementation?)
	return Statement{
		expression: value,
		print:      nil,
	}
}

func (p *Parser) printStatement() Statement {
	value := p.expression()
	if p.peek().t != SEMICOLON {
		return Statement{} // this is wrong. You messed up
	}
	p.advance(1)
	return Statement{
		expression: nil,
		print:      value,
	}
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
	for op.t == LESS || op.t == LESS_EQUAL || op.t == EQUAL || op.t == GREATER_EQUAL {
		operand := op
		p.advance(2) // don't know if this should be allowed hmmmmmmmmmammm
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
	return LiteralExpression{value: nil} // should probably error instead
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
