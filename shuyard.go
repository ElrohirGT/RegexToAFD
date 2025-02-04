package main

// You can rename package names when importing them!
// Here the "l" alias is being used!
import (
	"log"

	l "github.com/ElrohirGT/RegexToAFD/lib"
)

// Maps an operator in the form of a rune into a precedence number.
// Smaller means it has more priority
var precedence = map[rune]int{
	'|': 2, // OR Operator
	'.': 2, // AND Operator
	'*': 1, // ZERO_OR_MORE
	'?': 1, // ONE_OR_MORE
}

func toOperator(self rune) l.Optional[l.Operator] {
	switch self {
	case '|':
		return l.CreateValue(l.OR)
	case '.':
		return l.CreateValue(l.AND)
	case '?':
		return l.CreateValue(l.ONE_OR_MANY)
	case '*':
		return l.CreateValue(l.ZERO_OR_MANY)
	default:
		return l.CreateNull[l.Operator]()
	}
}

func tryToAppendWithPrecedence(stack *l.Stack[rune], operator rune, output *[]l.RX_Token) {
	if stack.Empty() {
		log.Default().Printf("Adding %c to stack!", operator)
		stack.Push(operator)
		return
	}

	top := stack.Peek()
	currentPrecedence := precedence[operator]
	stackPrecedence, found := precedence[top.GetValue()]

	log.Default().Printf("Checking if it can add operator directly %d > %d...", stackPrecedence, currentPrecedence)
	if !found || stackPrecedence > currentPrecedence {
		log.Default().Printf("Adding %c to stack!", operator)
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

		log.Default().Printf("Adding %c to stack!", operator)
		stack.Push(operator)
	}
}

func ToPostfix(infixExpression string) []l.RX_Token {
	stack := l.Stack[rune]{}
	output := []l.RX_Token{}

	previousCanBeANDedTo := true
	for i, char := range infixExpression {
		switch char {
		case '|':
			if stack.Empty() {
				log.Default().Printf("Adding %c to stack!", char)
				stack.Push(char)
			} else {
				tryToAppendWithPrecedence(&stack, char, &output)
			}
			previousCanBeANDedTo = false

		case '?', '*':
			if stack.Empty() {
				log.Default().Printf("Adding %c to stack!", char)
				stack.Push(char)
			} else {
				tryToAppendWithPrecedence(&stack, char, &output)
			}
			previousCanBeANDedTo = true

		default:
			if previousCanBeANDedTo && i != 0 {
				log.Default().Printf("Trying to append '.' operator...")
				tryToAppendWithPrecedence(&stack, '.', &output)
			}
			log.Default().Printf("Adding %c to output...", char)
			output = append(output, l.CreateValueToken(char))
			previousCanBeANDedTo = true
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
