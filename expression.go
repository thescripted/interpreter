package main

// expressions are evaluated with go's built-in primitives and stuff.

// panicking makes more sense here. If there is an user error deep inside some call stack somewhere in our application
// only the top-level error catcher should report it right?]

// No returned errors. We're panicking instead. All internal tools should return their error + details (Runtime error, Syntax, etc...)
// error.js will capture all panicked errors and hand them off to the enduser correctly.
type Expression interface { // evaluate()
	evaluate() interface{} // this value is ambigious, right?
}

type BinaryExpression struct {
	left    Expression
	operand Token // Operator? ExpressionOperator?
	right   Expression
}

// TODO: comparisons and equalities are missing here.
func (e BinaryExpression) evaluate() interface{} { // to panic or to return an error....
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
			if _, ok := right.(float64); ok { // this is ridicoulous.
				return left.(float64) + right.(float64)
			}
		}
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) <= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	}
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
		return -right.(float64) // what if the user does -"I am a dumbass user" ?? not all casts are ok. Should throw error.
	case BANG:
		return !right.(bool) // same here.
	}
	return nil // this is bad. If this happens, then you--the code writer--messed up.
}

type GroupingExpression struct {
	expression Expression
}

func (e GroupingExpression) evaluate() interface{} {
	return e.expression.evaluate()
}

func isEqual(l, r interface{}) bool {
	if l == nil && r == nil {
		return true
	}
	if l == nil {
		return false
	}
	return l == r
}
