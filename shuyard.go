package main

// You can rename package names when importing them!
// Here the "l" alias is being used!
import (
	"log"

	l "github.com/ElrohirGT/RegexToAFD/lib"
)

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
	case '.':
		return l.CreateValue(l.AND)
	default:
		return l.CreateNull[l.Operator]()
	}
}

func tryToAppendWithPrecedence(stack *l.Stack[rune], operator rune, output *[]l.RX_Token) {
	if stack.Empty() {
		stack.Push(operator)
		return
	}

	top := stack.Peek()
	currentPrecedence := precedence[operator]
	stackPrecedence, found := precedence[top.GetValue()]

	if !found || stackPrecedence > currentPrecedence {
		stack.Push(operator)
	} else {
		for stackPrecedence < currentPrecedence {
			poppedRune := stack.Pop()

			op := toOperator(poppedRune.GetValue())
			log.Default().Printf("Adding %c to output...", poppedRune.GetValue())
			*output = append(*output, l.CreateOperatorToken(op.GetValue()))

			if stack.Empty() {
				break
			}

			top := stack.Peek()
			stackPrecedence, found = precedence[top.GetValue()]
			if !found {
				break
			}
		}

		stack.Push(operator)
	}
}

func ToPostfix(infixExpression string) []l.RX_Token {
	stack := l.Stack[rune]{}
	output := []l.RX_Token{}

	previousCharWasNotOperator := true
	for i, char := range infixExpression {
		switch char {
		case '|':
			if stack.Empty() {
				log.Default().Printf("Adding %c to stack!", char)
				stack.Push(char)
			} else {
				tryToAppendWithPrecedence(&stack, char, &output)
			}
			previousCharWasNotOperator = false

		default:
			if previousCharWasNotOperator && i != 0 {
				log.Default().Printf("Trying to append '.' operator...")
				tryToAppendWithPrecedence(&stack, '.', &output)
			}
			log.Default().Printf("Adding %c to output...", char)
			output = append(output, l.CreateValueToken(char))
			previousCharWasNotOperator = true
		}
	}

	for !stack.Empty() {
		val := stack.Pop()
		op := toOperator(val.GetValue())

		if op.HasValue() {
			log.Default().Printf("Adding %c to output...", val.GetValue())
			output = append(output, l.CreateOperatorToken(op.GetValue()))
		}
	}

	return output
}
