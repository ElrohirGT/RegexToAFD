package lib

type Operator int

const (
	OR Operator = iota
)

// Represents a token.
// It can either be a value or an operator between two values.
// If value is null then it should have an operator value, otherwise a value should be provided!
type RX_Token struct {
	operator *Operator
	value    *string
}

func CreateOperatorToken(t *Operator) RX_Token {
	return RX_Token{
		operator: t,
	}
}

func CreateValueToken(value *string) RX_Token {
	return RX_Token{
		value: value,
	}
}
