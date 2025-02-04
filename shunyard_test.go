package main

import "testing"
import "github.com/ElrohirGT/RegexToAFD/lib"

type testInfo struct {
	input    string
	expected []lib.RX_Token
}

func test(t *testing.T, info testInfo) {
	result := ToPostfix(info.input)

	for i, expected := range info.expected {
		value := result[i]

		if expected.Equals(value) {
			t.Fatalf("The result doesn't match expected!\n%v != %v", expected, value)
		}
	}

	resultLength := len(result)
	expectedLength := len(info.expected)
	if expectedLength != resultLength {
		t.Fatalf("The lengths don't match! %d != %d", expectedLength, resultLength)
	}
}

func TestSimple(t *testing.T) {
	regexp := "a|b"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken("a"),
			lib.CreateValueToken("b"),
			lib.CreateOperatorToken(lib.OR),
		},
	})

	regexp = "abc"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken("a"),
			lib.CreateValueToken("b"),
			lib.CreateValueToken("c"),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateOperatorToken(lib.AND),
		},
	})

	regexp = "a|bc"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken("a"),
			lib.CreateValueToken("b"),
			lib.CreateValueToken("c"),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateOperatorToken(lib.OR),
		},
	})

	regexp = "abc?"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken("a"),
			lib.CreateValueToken("b"),
			lib.CreateValueToken("c"),
			lib.CreateOperatorToken(lib.ONE_OR_MANY),
			lib.CreateOperatorToken(lib.AND),
			lib.CreateOperatorToken(lib.AND),
		},
	})

}

func TestRanges(t *testing.T) {
	regexp := "[a-c]"
	test(t, testInfo{
		input: regexp,
		expected: []lib.RX_Token{
			lib.CreateValueToken("a"),
			lib.CreateValueToken("b"),
			lib.CreateOperatorToken(lib.OR),
			lib.CreateValueToken("c"),
			lib.CreateOperatorToken(lib.OR),
		},
	})
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
