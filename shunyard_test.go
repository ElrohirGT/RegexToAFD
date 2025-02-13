package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ElrohirGT/RegexToAFD/lib"
)

type testInfo struct {
	input    string
	expected []lib.RX_Token
	alphabet Alphabet
}

func toString(stream *[]lib.RX_Token) string {
	sb := strings.Builder{}
	sb.WriteString("{ ")
	for _, token := range *stream {
		if token.IsOperator() {
			op := *token.GetOperator()
			displayOp := "invalid"

			switch op {
			case lib.OR:
				displayOp = "OR"
			case lib.AND:
				displayOp = "AND"
			case lib.ZERO_OR_MANY:
				displayOp = "*"
			}

			sb.WriteString(
				fmt.Sprintf("[op: %s] ", displayOp),
			)
		} else {
			val := "epsilon"
			if token.GetValue().HasValue() {
				val = string(token.GetValue().GetValue())
			}

			sb.WriteString(
				fmt.Sprintf("[val: %s]", val),
			)
		}
	}
	sb.WriteString(" }")

	return sb.String()
}

func assertEquals(t *testing.T, expected []lib.RX_Token, actual []lib.RX_Token, originalInput string) {
	resultLength := len(actual)
	expectedLength := len(expected)

	if expectedLength != resultLength {
		t.Errorf("The lengths don't match! %d != %d\nActual  : %+v\nExpected: %+v\nFailed on: %s",
			expectedLength, resultLength, toString(&actual), toString(&expected), originalInput)
	}

	for i, expected := range expected {
		value := actual[i]

		if expected.Equals(value) {
			t.Fatalf("The result doesn't match expected!\n%v != %v\nFailed on: %s", expected, value, originalInput)
		}
	}
}

func test(t *testing.T, info testInfo) {
	result := info.alphabet.ToPostfix(info.input)

	assertEquals(t, info.expected, result, info.input)
}

func TestSimpleOr(t *testing.T) {
	regexp := "a|b"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('b'),
			lib.CreateOperatorToken(lib.OR),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestSimpleAnd(t *testing.T) {
	regexp := "abc"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('b'),
			lib.CreateValueToken('c'),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateOperatorToken(lib.AND),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestSimpleCombination(t *testing.T) {
	regexp := "a|bc"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('b'),
			lib.CreateValueToken('c'),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateOperatorToken(lib.OR),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestOptional(t *testing.T) {
	regexp := "abc?"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('b'),
			lib.CreateValueToken('c'),
			lib.CreateEpsilonValue(),
			lib.CreateOperatorToken(lib.OR),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateOperatorToken(lib.AND),
		},
		alphabet: DEFAULT_ALPHABET,
	})

	regexp = "ab?c"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('b'),
			lib.CreateEpsilonValue(),
			lib.CreateOperatorToken(lib.OR),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateValueToken('c'),
			lib.CreateOperatorToken(lib.AND),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestZeroOrMore(t *testing.T) {
	regexp := "abc*"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('b'),
			lib.CreateValueToken('c'),
			lib.CreateOperatorToken(lib.ZERO_OR_MANY),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateOperatorToken(lib.AND),
		},
		alphabet: DEFAULT_ALPHABET,
	})

	regexp = "ab*c"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('b'),
			lib.CreateOperatorToken(lib.ZERO_OR_MANY),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateValueToken('c'),
			lib.CreateOperatorToken(lib.AND),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestParenthesis(t *testing.T) {
	regexp := "b|(a*bc)"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('b'),
			lib.CreateValueToken('a'),
			lib.CreateOperatorToken(lib.ZERO_OR_MANY),
			lib.CreateValueToken('b'),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateValueToken('c'),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateOperatorToken(lib.OR),
		},
		alphabet: DEFAULT_ALPHABET,
	})

	regexp = "b|(ac)|o"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('b'),
			lib.CreateValueToken('a'),
			lib.CreateValueToken('c'),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateOperatorToken(lib.OR),
			lib.CreateValueToken('o'),
			lib.CreateOperatorToken(lib.OR),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestOrBrackets(t *testing.T) {
	regexp := "[ath]"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('t'),
			lib.CreateOperatorToken(lib.OR),
			lib.CreateValueToken('h'),
			lib.CreateOperatorToken(lib.OR),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestRanges(t *testing.T) {
	regexp := "[a-b]"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('b'),
			lib.CreateOperatorToken(lib.OR),
		},
		alphabet: DEFAULT_ALPHABET,
	})

	regexp = "[a-a]"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestLargeRange(t *testing.T) {
	regexp := "[a-d]"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('b'),
			lib.CreateOperatorToken(lib.OR),
			lib.CreateValueToken('c'),
			lib.CreateOperatorToken(lib.OR),
			lib.CreateValueToken('d'),
			lib.CreateOperatorToken(lib.OR),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestEscapeSequences(t *testing.T) {
	regexp := "\\a\\[\\*"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('['),
			lib.CreateValueToken('*'),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestOneOrMore(t *testing.T) {
	regexp := "a+bh"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('a'),
			lib.CreateOperatorToken(lib.ZERO_OR_MANY),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateValueToken('b'),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateValueToken('h'),
			lib.CreateOperatorToken(lib.AND),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestOneOrMoreComplicated(t *testing.T) {
	regexp := "(bh)+"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('b'),
			lib.CreateValueToken('h'),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateValueToken('b'),
			lib.CreateValueToken('h'),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateOperatorToken(lib.ZERO_OR_MANY),
			lib.CreateOperatorToken(lib.AND),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestOneOrMoreRecursive(t *testing.T) {
	regexp := "(b+a)+"

	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('b'),
			lib.CreateValueToken('b'),
			lib.CreateOperatorToken(lib.ZERO_OR_MANY),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateValueToken('a'),
			lib.CreateOperatorToken(lib.AND),

			lib.CreateValueToken('b'),
			lib.CreateValueToken('b'),
			lib.CreateOperatorToken(lib.ZERO_OR_MANY),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateValueToken('a'),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateOperatorToken(lib.ZERO_OR_MANY),
			lib.CreateOperatorToken(lib.AND),
		},
		alphabet: DEFAULT_ALPHABET,
	})
}

func TestNotRanges(t *testing.T) {
	regexp := "[^c-z]"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('b'),
			lib.CreateOperatorToken(lib.OR),
			lib.CreateValueToken('ñ'),
			lib.CreateOperatorToken(lib.OR),
		},
		alphabet: NewAlphabetFromString("abcdefghijklmnñopqrstuvwxyz"),
	})

}

func TestBigNotRange(t *testing.T) {
	regexp := "[^e-z]"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('b'),
			lib.CreateOperatorToken(lib.OR),
			lib.CreateValueToken('c'),
			lib.CreateOperatorToken(lib.OR),
			lib.CreateValueToken('d'),
			lib.CreateOperatorToken(lib.OR),
		},
		alphabet: NewAlphabetFromString("abcdefghijklmnopqrstuvwxyz"),
	})
}

func TestLogicEquivalence(t *testing.T) {
	reducedResult := DEFAULT_ALPHABET.ToPostfix("(b+a)+")
	expandedResult := DEFAULT_ALPHABET.ToPostfix("(bb*a)(bb*a)*")

	assertEquals(t, reducedResult, expandedResult, "Postfix inconsistency: `(b+a)+` != `(bb*a)(bb*a)*`")
}

// func TestPasswordPolicy(t *testing.T) {
// 	regexp := "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{12,}$"
//
// 	test(t, testInfo{
// 		input:    regexp,
// 		expected: "",
// 	})
// }
//
// func TestExtractURLFromText(t *testing.T) {
// 	regexp := "\\b((?:https?|ftp):\\/\\/[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|])"
//
// 	test(t, testInfo{
// 		input:    regexp,
// 		expected: "",
// 	})
// }
//
// func TestBalancedParenthesis(t *testing.T) {
// 	regexp := "\\((?:[^()]+|(?R))*\\)"
//
// 	test(t, testInfo{
// 		input:    regexp,
// 		expected: "",
// 	})
// }
//
// func TestMatchAmericanNumber(t *testing.T) {
// 	regexp := "^(\\+1\\s?)?(\\(?\\d{3}\\)?[\\s.-]?)?\\d{3}[\\s.-]?\\d{4}$"
//
// 	test(t, testInfo{
// 		input:    regexp,
// 		expected: "",
// 	})
// }
//
// func TestPrinceRegexp(t *testing.T) {
// 	regexp := "^[a-zA-Z0-9\\._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
//
// 	test(t, testInfo{
// 		input:    regexp,
// 		expected: "",
// 	})
// }
