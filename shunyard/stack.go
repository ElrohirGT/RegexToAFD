package shunyard

type Stack[T any] []T

func (self *Stack[T]) Push(val T) *Stack[T] {
	*self = append(*self, val)
	return self
}

func (self *Stack[T]) Pop() T {
	length := len(*self)

	ref := *self
	val := ref[length-1]
	*self = ref[:length-1]

	return val
}
