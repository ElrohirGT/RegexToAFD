package lib

type Set[T comparable] map[T]struct{}

func (self *Set[T]) Add(val T) bool {
	ref := *self
	_, alreadyAdded := ref[val]

	if !alreadyAdded {
		ref[val] = struct{}{}
	}

	return !alreadyAdded
}
