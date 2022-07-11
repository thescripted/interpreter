package main

import (
	"errors"
)

// expressions are evaluated with go's built-in primitives and stuff.

// panicking makes more sense here. If there is an user error deep inside some call stack somewhere in our application
// only the top-level error catcher should report it right?
type Expression interface { // evaluate()
	evaluate() (interface{}, error) // this value is ambigious, right?
}

type BinaryExpression struct {
	left    Expression
	operand Token // Operator? ExpressionOperator?
	right   Expression
}

// TODO: comparisons and equalities are missing here.
func (e BinaryExpression) evaluate() (interface{}, error) { // to panic or to return an error....
	left, err := e.left.evaluate()
	if err != nil {
		return nil, err
	}
	right, err := e.right.evaluate()
	if err != nil {
		return nil, err
	}
	switch e.operand.t {
	case MINUS:
		return left.(float64) - right.(float64), nil // this ok? These will panic if not true. Should handle those errors before.
	case SLASH:
		return left.(float64) / right.(float64), nil // this ok?
	case STAR:
		return left.(float64) * right.(float64), nil // this ok?
	case PLUS: // we can string concat. Typeswitch on that.
		switch left.(type) {
		case string:
			if _, ok := right.(string); ok {
				return left.(string) + right.(string), nil
			}
		case float64: // maybe others?
			if _, ok := right.(float64); ok { // this is fucking ridicoulous.
				return left.(float64) + right.(float64), nil
			}
		}
	}
	return nil, errors.New("dumbass") // if you reach thhis point, then you fucked up.
}

type LiteralExpression struct {
	value interface{}
}

func (e LiteralExpression) evaluate() (interface{}, error) {
	return e.value, nil
}

type UnaryExpression struct {
	operand Token
	right   Expression
}

func (e UnaryExpression) evaluate() (interface{}, error) {
	right, err := e.right.evaluate()
	if err != nil {
		return nil, err
	}
	switch e.operand.t {
	case MINUS:
		return -right.(float64), nil // what if the user does -"I am a dumbass user" ?? not all casts are ok. Should throw error.
	case BANG:
		return !right.(bool), nil // same here.
	}
	return nil, errors.New("idiot") // if you reach thhis point, then you fucked up.
}

type GroupingExpression struct {
	expression Expression
}

func (e GroupingExpression) evaluate() (interface{}, error) {
	return e.expression.evaluate()
}
