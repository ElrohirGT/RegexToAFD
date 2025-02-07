package main

// You can rename package names when importing them!
// Here the "l" alias is being used!
import (
	l "github.com/ElrohirGT/RegexToAFD/lib"
	"log"
)

// Maps an operator in the form of a rune into a precedence number.
// Smaller means it has more priority
var precedence = map[byte]int{
	'|': 2, // OR Operator
	'.': 2, // AND Operator
	'*': 1, // ZERO_OR_MORE
	'?': 1, // ZERO_OR_ONE
}

func toOperator(self byte) l.Optional[l.Operator] {
	log.Default().Printf("Trying to get operator from: %c", self)

	switch self {
	case '|':
		return l.CreateValue(l.OR)
	case '.':
		return l.CreateValue(l.AND)
	case '*':
		return l.CreateValue(l.ZERO_OR_MANY)
	default:
		return l.CreateNull[l.Operator]()
	}
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func tryToAppendWithPrecedence(stack *l.Stack[byte], operator byte, output *[]l.RX_Token) {
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
			poppedRune := stack.Pop().GetValue()

			if poppedRune == '?' {
				*output = append(*output, l.CreateEpsilonValue())
				*output = append(*output, l.CreateOperatorToken(l.OR))
				// } else if poppedRune == '+' {
				// 	*output = append(*output, l.CreateEpsilonValue())
			} else {
				op := toOperator(poppedRune)
				log.Default().Printf("Adding %c to output...", poppedRune)
				*output = append(*output, l.CreateOperatorToken(op.GetValue()))
			}

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

type RegexState int

const (
	NORMAL RegexState = iota
	IN_BRACKETS
	IN_PARENTHESIS
)

type ShunStack = l.Stack[byte]
type ShunOutput = []l.RX_Token

func toPostFix(infixExpression *string, stack *ShunStack, output *ShunOutput) {
	infixExpr := *infixExpression
	previousCanBeANDedTo := false
	state := NORMAL

	previousExprStack := l.ExprStack{}
	for i := 0; i < len(infixExpr); i++ {
		char := infixExpr[i]

		switch char {
		case '|':
			if stack.Empty() {
				log.Default().Printf("Adding %c to stack!", char)
				stack.Push(char)
			} else {
				tryToAppendWithPrecedence(stack, char, output)
			}
			previousCanBeANDedTo = false

		case '?', '*':
			if stack.Empty() {
				log.Default().Printf("Adding %c to stack!", char)
				stack.Push(char)
			} else {
				tryToAppendWithPrecedence(stack, char, output)
			}
			previousCanBeANDedTo = true

		case '(':
			stack.Push('(')
			previousCanBeANDedTo = false
			state = IN_PARENTHESIS

			expr := ""
			if !previousExprStack.IsEmpty() {
				expr = previousExprStack.Peek().GetValue()
			}

			log.Default().Printf("The previous expression before deleting is: %s", expr)
			previousExprStack.Pop()     // Deletes previous expression
			previousExprStack.Push("(") // Adds ( context
			previousExprStack.Push("")  // Adds inner ( ) context

		case ')':
			log.Default().Printf("Popping until it finds: '('")
			for peeked := stack.Peek(); peeked.GetValue() != '('; peeked = stack.Peek() {
				val := stack.Pop()
				op := toOperator(val.GetValue()).GetValue()

				*output = append(*output, l.CreateOperatorToken(op))
			}

			// Popping '('
			stack.Pop()
			state = NORMAL
			previousExprStack.AppendTop(")")
			previousExprStack.Pop() // Popping inner ( ) context

		case '[':
			stack.Push('[')
			previousCanBeANDedTo = false
			state = IN_BRACKETS

			expr := ""
			if !previousExprStack.IsEmpty() {
				expr = previousExprStack.Peek().GetValue()
			}
			log.Default().Printf("The previous expression before deleting is: %s", expr)
			previousExprStack.Pop()     // Deletes previous expression
			previousExprStack.Push("[") // Adds [ context
			previousExprStack.Push("")  // Adds inner [ ] context

		case ']':
			log.Default().Printf("Popping until it finds: '['")
			for peeked := stack.Peek(); peeked.GetValue() != '['; peeked = stack.Peek() {
				val := stack.Pop()
				op := toOperator(val.GetValue()).GetValue()

				*output = append(*output, l.CreateOperatorToken(op))
			}

			// Popping '['
			stack.Pop()
			state = NORMAL
			previousExprStack.AppendTop("]")
			previousExprStack.Pop() // Popping inner [ ] context

		case '+':
			previousExpr := previousExprStack.Pop().GetValue()
			log.Default().Printf("'+' found! Recursing with: %s ...", previousExpr)
			toPostFix(&previousExpr, stack, output)
			tryToAppendWithPrecedence(stack, '*', output)
			tryToAppendWithPrecedence(stack, '.', output)
			previousExprStack.Push("")

		case '\\':
			nextChar := infixExpr[i+1]
			log.Default().Printf("Escape sequence found! Adding %c as a char...", nextChar)
			*output = append(*output, l.CreateValueToken(rune(nextChar)))
			i += 1

		default:
			log.Default().Printf("Iteration: (%c) %d != 0 && previousCanBeANDed: %t", char, i, previousCanBeANDedTo)
			if i != 0 && previousCanBeANDedTo {
				if state == NORMAL {
					log.Default().Printf("Trying to append '.' operator...")
					tryToAppendWithPrecedence(stack, '.', output)
				} else {
					log.Default().Printf("Trying to append '|' operator...")
					tryToAppendWithPrecedence(stack, '|', output)
				}
			}

			rangeStart := byte(char)
			if state == IN_BRACKETS {
				log.Default().Printf("Checking if the char (%c) is a range start...", char)
				if isLetter(rangeStart) || isDigit(rangeStart) {
					nextChar := infixExpr[i+1]

					if nextChar == '-' {
						rangeEnd := infixExpr[i+2]
						isEndTheSameAsStart := (isLetter(rangeStart) && isLetter(rangeEnd)) || (isDigit(rangeStart) && isDigit(rangeStart))

						log.Default().Printf("The end char (%c) is the same type as start? %v", rangeEnd, isEndTheSameAsStart)
						if isEndTheSameAsStart {
							if rangeEnd < rangeStart {
								rangeEnd, rangeStart = rangeStart, rangeEnd
							}

							for j := byte(0); j <= (rangeEnd - rangeStart); j++ {
								val := rune(rangeStart + j)
								log.Default().Printf("Adding %c to output...", val)
								*output = append(*output, l.CreateValueToken(val))

								if j >= 1 {
									tryToAppendWithPrecedence(stack, '|', output)
								}
							}

							// We already parsed '-' and the other byte
							// So we need to ignore them
							i += 2
							continue
						}
					}
				}
			}

			if state == IN_BRACKETS || state == IN_PARENTHESIS {
				expr := ""
				if !previousExprStack.IsEmpty() {
					expr = previousExprStack.Peek().GetValue()
				}
				log.Default().Printf("Appending %s to expression: %s", string(char), expr)
				previousExprStack.AppendTop(string(char))
			} else {

				expr := ""
				if !previousExprStack.IsEmpty() {
					expr = previousExprStack.Peek().GetValue()
				}

				log.Default().Printf("Changing previous expr from %s to %s", expr, string(char))
				previousExprStack.Pop()
				previousExprStack.Push(string(char))
			}

			log.Default().Printf("Adding %c to output...", char)
			*output = append(*output, l.CreateValueToken(rune(char)))
			previousCanBeANDedTo = true
		}
	}

	for !stack.Empty() {
		val := stack.Pop().GetValue()
		op := toOperator(val)

		if val == '?' {
			*output = append(*output, l.CreateEpsilonValue())
			*output = append(*output, l.CreateOperatorToken(l.OR))
		} else if op.HasValue() {
			log.Default().Printf("Adding %c to output...", val)
			*output = append(*output, l.CreateOperatorToken(op.GetValue()))
		}
	}
}

func ToPostfix(infixExpression string) []l.RX_Token {
	stack := l.Stack[byte]{}
	output := []l.RX_Token{}

	toPostFix(&infixExpression, &stack, &output)

	return output
}
