// all statements can be evaluated. Statements are code followed by a semicolon.

package main

import "fmt"

type Statement interface {
	eval()
}

type VariableStatment struct {
	name       string
	expression Expression
}

type PrintStatement struct {
	expression Expression
}

type ExpressionStatement struct {
	expression Expression
}

func (v VariableStatment) eval() {
	if v.expression == nil {
		_globalEnvironment.put(v.name, nil)
	} else {
		value := v.expression.evaluate()
		_globalEnvironment.put(v.name, value)
	}
}

func (p PrintStatement) eval() {
	value := p.expression.evaluate()
	fmt.Println(value)
}

func (e ExpressionStatement) eval() {
	e.expression.evaluate()
}
