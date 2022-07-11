package main

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Note: there is no reset. Once you parse it once, you can't parse it again. Is that ok?
func (p *Parser) Parse() Expression {
	return p.expression()
}

func (p *Parser) expression() Expression {
	return p.equality()
}

func (p *Parser) equality() Expression {
	left := p.comparison()
	if op := p.peek(); op.t == BANG_EQUAL || op.t == EQUAL_EQUAL {
		p.advance()
		p.advance() // can syntax error. Should handle.
		right := p.comparison()
		return BinaryExpression{
			left:    left,
			operand: op,
			right:   right,
		}
	}
	return left
}

func (p *Parser) comparison() Expression {
	left := p.term()
	if op := p.peek(); op.t == LESS || op.t == LESS_EQUAL || op.t == EQUAL || op.t == GREATER_EQUAL {
		p.advance()
		p.advance() // can syntax error. Should handle.
		right := p.term()
		return BinaryExpression{
			left:    left,
			operand: op,
			right:   right,
		}
	}
	return left
}

func (p *Parser) term() Expression {
	left := p.factor()
	if op := p.peek(); op.t == PLUS || op.t == MINUS {
		p.advance()
		p.advance() // can syntax error. Should handle.
		right := p.factor()
		return BinaryExpression{
			left:    left,
			operand: op,
			right:   right,
		}
	}
	return left
}

func (p *Parser) factor() Expression {
	left := p.unary()
	if op := p.peek(); op.t == STAR || op.t == SLASH {
		p.advance()
		p.advance() // can syntax error. Should handle.
		right := p.unary()
		return BinaryExpression{
			left:    left,
			operand: op,
			right:   right,
		}
	}
	return left
}

func (p *Parser) unary() Expression {
	if op := p.currentToken(); op.t == BANG || op.t == MINUS {
		p.advance() // can syntax error. Should handle.
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
	case NUMBER:
		return LiteralExpression{value: literal.lexeme} // should be value
	case STRING:
		return LiteralExpression{value: literal.lexeme} // should be value
	case TRUE:
		return LiteralExpression{value: literal.lexeme} // should be value
	case FALSE:
		return LiteralExpression{value: literal.lexeme} // should be value
	case NIL:
		return LiteralExpression{value: literal.lexeme} // should be value

	}
	return LiteralExpression{value: nil} // should probably error
}

// peek looks at what the next token is. Returns a nil token if there is not a next one.
func (p *Parser) peek() Token {
	if p.current+1 >= len(p.tokens) {
		return Token{} // handle this correclty.
	}
	return p.tokens[p.current+1]
}

// advanceToken moves to the next token
func (p *Parser) advance() {
	p.current++
}

// currentToken returns the current token
func (p *Parser) currentToken() Token {
	return p.tokens[p.current]
}
