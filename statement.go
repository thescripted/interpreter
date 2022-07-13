// statements are not expressions. That's why they're separated. but they feed into the parser the same way-ish, right?
// might update.

package main

import "fmt"

type Statement struct { // I don't like this struct. Fix it.
	expression Expression
	print      Expression
}

func (s Statement) evaluate() interface{} {
	if s.expression != nil { // this might 1) be wrong, and 2) be inefficient
		return s.expression.evaluate()
	}
	fmt.Println(s.print.evaluate())
	return nil
}
