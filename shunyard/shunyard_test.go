package shunyard

import "testing"

type testInfo struct {
	input    string
	expected string
}

func test(t *testing.T, info testInfo) {
	result := ToPostfix(info.input)

	if result != info.expected {
		t.Fatalf("The result doesn't match expected!\n%s", result)
	}
}

func TestPasswordPolicy(t *testing.T) {
	regexp := "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{12,}$"

	test(t, testInfo{
		input:    regexp,
		expected: "",
	})
}

func TestExtractURLFromText(t *testing.T) {
	regexp := "\\b((?:https?|ftp):\\/\\/[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|])"

	test(t, testInfo{
		input:    regexp,
		expected: "",
	})
}

func TestBalancedParenthesis(t *testing.T) {
	regexp := "\\((?:[^()]+|(?R))*\\)"

	test(t, testInfo{
		input:    regexp,
		expected: "",
	})
}

func TestMatchAmericanNumber(t *testing.T) {
	regexp := "^(\\+1\\s?)?(\\(?\\d{3}\\)?[\\s.-]?)?\\d{3}[\\s.-]?\\d{4}$"

	test(t, testInfo{
		input:    regexp,
		expected: "",
	})
}
