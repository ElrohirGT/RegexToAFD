package shunyard

import "testing"

func TestStack(t *testing.T) {
	stack := Stack[int]{}

	stack.
		Push(50).
		Push(100).
		Push(150)

	if val := stack.Pop(); val != 150 {
		t.Fatalf("150 != %d", val)
	}
	if val := stack.Pop(); val != 100 {
		t.Fatalf("100 != %d", val)
	}
	if val := stack.Pop(); val != 50 {
		t.Fatalf("50 != %d", val)
	}
}
