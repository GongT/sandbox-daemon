package types

type Set[T comparable] []T

func (s Set[T]) Has(item T) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

func (s *Set[T]) Add(item T) {
	if !s.Has(item) {
		*s = append(*s, item)
	}
}

func (s *Set[T]) Delete(item T) {
	for i, v := range *s {
		if v == item {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return
		}
	}
}
