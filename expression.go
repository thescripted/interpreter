package main

import (
	"fmt"
	"log"
)

// expressions are evaluated with go's built-in primitives and stuff.

type Expression interface { // evaluate()
	evaluate() interface{} // this value is ambigious, right?
}

type BinaryExpression struct {
	left    Expression
	operand Token // Operator? ExpressionOperator?
	right   Expression
}

// TODO: comparisons and equalities are missing here.
func (e BinaryExpression) evaluate() interface{} {
	left := e.left.evaluate()
	right := e.right.evaluate()
	switch e.operand.t {
	case MINUS:
		return left.(float64) - right.(float64) // this ok? These will panic if not true. Should handle those errors before.
	case SLASH:
		return left.(float64) / right.(float64) // this ok?
	case STAR:
		return left.(float64) * right.(float64) // this ok?
	case PLUS: // we can string concat. Typeswitch on that.
		switch left.(type) {
		case string:
			if _, ok := right.(string); ok {
				return left.(string) + right.(string)
			}
		case float64: // maybe others?
			if _, ok := right.(float64); ok { // this is fucking ridicoulous.
				return left.(float64) + right.(float64)
			}
		}
	}
	fmt.Println("HELLO??")
	log.Fatalf("who designed this compiler that allowed this error to fall through???")
	return nil // if you reach thhis point, then you fucked up.
}

type LiteralExpression struct {
	value interface{}
}

func (e LiteralExpression) evaluate() interface{} {
	return e.value
}

type UnaryExpression struct {
	operand Token
	right   Expression
}

func (e UnaryExpression) evaluate() interface{} {
	right := e.right.evaluate()
	switch e.operand.t {
	case MINUS:
		return -right.(float64) // what if the user does -"I am a dumbass user" ?? not all casts are ok.
	case BANG:
		return !right.(bool) // same here.
	}
	log.Fatalf("who designed this compiler that allowed this error to fall through???")
	return nil // if you reach thhis point, then you fucked up.
}

type GroupingExpression struct {
	expression Expression
}

func (e GroupingExpression) evaluate() interface{} {
	return e.expression.evaluate()
}
