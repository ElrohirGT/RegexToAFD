package main

import "./lib"

var precedence = map[string]int{
	"+": 1,
	".": 2,
	"*": 3,
}

type Operator int

const (
	OR Operator = iota
)

func ToPostfix(infixExpression string) string {
	return "not implemented"
}
