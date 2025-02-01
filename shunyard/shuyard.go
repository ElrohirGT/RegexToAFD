package shunyard

var precedence = map[string]int{
	"+": 1,
	".": 2,
	"*": 3,
}

func ToPostfix(infixExpression string) string {
	operators := Stack[string]{}
}
