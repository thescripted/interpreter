package main

type Expression interface { // evaluate()
}

type BinaryExpression struct {
	left    Expression
	operand Token // Operator? ExpressionOperator?
	right   Expression
}

type LiteralExpression struct {
	value interface{}
}

type UnaryExpression struct {
	operand Token
	right   Expression
}
