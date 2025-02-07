package main

import "testing"
import "github.com/ElrohirGT/RegexToAFD/lib"

type testInfo struct {
	input    string
	expected []lib.RX_Token
}

func test(t *testing.T, info testInfo) {
	result := ToPostfix(info.input)

	resultLength := len(result)
	expectedLength := len(info.expected)
	if expectedLength != resultLength {
		t.Errorf("The lengths don't match! %d != %d\nResult: %+v\nExpected: %+v\nFailed on: %s",
			expectedLength, resultLength, result, info.expected, info.input)
	}

	for i, expected := range info.expected {
		value := result[i]

		if expected.Equals(value) {
			t.Fatalf("The result doesn't match expected!\n%v != %v\nFailed on: %s", expected, value, info.input)
		}
	}

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
			lib.CreateOperatorToken(lib.ONE_OR_MANY),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateOperatorToken(lib.AND),
		},
	})

	regexp = "ab?c"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
			lib.CreateValueToken('b'),
			lib.CreateOperatorToken(lib.ONE_OR_MANY),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateValueToken('c'),
			lib.CreateOperatorToken(lib.AND),
		},
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
	})

	regexp = "[a-a]"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken('a'),
		},
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
	})
}

// func TestNotRanges(t *testing.T) {
// 	regexp := "[^a-c]"
// 	test(t, testInfo{
// 		input: regexp,
// 		expected: []lib.RX_Token{
// 			lib.CreateValueToken('a'),
// 			lib.CreateValueToken('b'),
// 			lib.CreateOperatorToken(lib.OR),
// 			lib.CreateValueToken('c'),
// 			lib.CreateOperatorToken(lib.OR),
// 		},
// 	})
// }

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
