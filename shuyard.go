package main

// You can rename package names when importing them!
// Here the "l" alias is being used!
import l "github.com/ElrohirGT/RegexToAFD/lib"

var precedence = map[rune]int{
	'|': 1, // OR Operator
	'.': 2, // AND Operator
	'*': 3, // ZERO_OR_MORE
	'?': 3, // ONE_OR_MORE
}

func ToOperator(self rune) l.Optional[l.Operator] {
	switch self {
	case '|':
		return l.CreateValue(l.OR)
	default:
		return l.CreateNull[l.Operator]()
	}
}

func ToPostfix(infixExpression string) []l.RX_Token {
	stack := l.Stack[rune]{}
	output := []l.RX_Token{}

	for _, char := range infixExpression {
		switch char {
		case '|':
			top := stack.Peek()
			currentPrecedence := precedence[char]
			stackPrecedence, found := precedence[top]

			if !found || stackPrecedence > currentPrecedence {
				stack.Push(char)
			} else {
				for stackPrecedence < currentPrecedence {
					poppedRune := stack.Pop()
					output = append(output, l.CreateOperatorToken())
				}
			}

		default:
			output = append(output, l.CreateValueToken(string(char)))
		}

	}

	return output
}
