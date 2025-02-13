package lib

type Operator int

const (
	OR           Operator = iota // This OR that operator.
	AND                          // Concatenation operator.
	ZERO_OR_MANY                 // * Operator
)

// Represents a token.
// It can either be a value or an operator between two values.
// If value is null then it should have an operator value, otherwise a value should be provided!
type RX_Token struct {
	operator *Operator
	// If the value is null then this token is an operator.
	// If the optional doesn't have a value then the value is epsilon.
	// If the optional has a value then this token has the value of the rune.
	value *Optional[rune]
}

func (self *RX_Token) GetValue() *Optional[rune] {
	return self.value
}

func (self *RX_Token) GetOperator() *Operator {
	return self.operator
}

func (self *RX_Token) IsValue() bool {
	return self.value != nil
}

func (self *RX_Token) IsOperator() bool {
	return self.operator != nil
}

func CreateOperatorToken(t Operator) RX_Token {
	return RX_Token{
		operator: &t,
	}
}

func CreateValueToken(value rune) RX_Token {
	val := CreateValue(value)
	return RX_Token{
		value: &val,
	}
}

func CreateEpsilonValue() RX_Token {
	val := CreateNull[rune]()
	return RX_Token{
		value: &val,
	}
}

func (self RX_Token) Equals(other RX_Token) bool {
	bothOpsNil := self.operator == nil && other.operator == nil
	if !bothOpsNil {
		if self.operator == nil || other.operator == nil {
			return false
		}

		return *self.operator == *other.operator
	}

	return self.value.Equals(other.value)
}
