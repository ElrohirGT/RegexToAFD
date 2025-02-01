package main

// You can rename package names when importing them!
// Here the "l" alias is being used!
import l "github.com/ElrohirGT/RegexToAFD/lib"

var precedence = map[string]int{
	"+": 1,
	".": 2,
	"*": 3,
}

func ToPostfix(infixExpression string) []l.RX_Token {
	return []l.CreateValueToken("")
}
