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

func toOperator(self rune) l.Optional[l.Operator] {
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
			if stack.Empty() {
				stack.Push(char)
				continue
			}

			top := stack.Peek()
			currentPrecedence := precedence[char]
			// Since the stack is not empty, top can't be null!
			stackPrecedence, found := precedence[top.GetValue()]

			if !found || stackPrecedence > currentPrecedence {
				stack.Push(char)
			} else {
				for stackPrecedence < currentPrecedence {
					poppedRune := stack.Pop()

					op := toOperator(poppedRune.GetValue())
					output = append(output, l.CreateOperatorToken(op.GetValue()))

					if stack.Empty() {
						break
					}

					top := stack.Peek()
					stackPrecedence, found = precedence[top.GetValue()]
					if !found {
						break
					}
				}
			}

		default:
			output = append(output, l.CreateValueToken(char))
		}

	}

	for !stack.Empty() {
		val := stack.Pop()
		op := toOperator(val.GetValue())

		if op.HasValue() {
			output = append(output, l.CreateOperatorToken(op.GetValue()))
		}
	}

	return output
}
