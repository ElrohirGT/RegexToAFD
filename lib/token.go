package lib

type Operator int

const (
	OR           Operator = iota // This OR that operator.
	AND                          // Concatenation operator.
	RANGE                        // When we write {1,5}
	ONE_OR_MANY                  // ? Operator
	ZERO_OR_MANY                 // * Operator
)

// Represents a token.
// It can either be a value or an operator between two values.
// If value is null then it should have an operator value, otherwise a value should be provided!
type RX_Token struct {
	operator *Operator
	value    *rune
}

func CreateOperatorToken(t Operator) RX_Token {
	return RX_Token{
		operator: &t,
	}
}

func CreateValueToken(value rune) RX_Token {
	return RX_Token{
		value: &value,
	}
}

func (self RX_Token) Equals(other RX_Token) bool {
	return self.operator == other.operator && self.value == other.value
}
