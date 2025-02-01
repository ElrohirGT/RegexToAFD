package main

// To import modules you can use the module name with the path to the module
// import "github.com/ElrohirGT/RegexToAFD/lib"

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
