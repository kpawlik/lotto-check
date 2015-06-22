package lotto

type intSet struct {
	set map[int]bool
}

func NewIntSet() *intSet {
	return &intSet{make(map[int]bool)}
}

func NewIntSetFrom(ints []int) *intSet {
	s := NewIntSet()
	for i := 0; i < len(ints); i++ {
		s.set[ints[i]] = true
	}
	return s
}

func (s *intSet) Size() int {
	return len(s.set)
}

func (s *intSet) Add(i int) {
	s.set[i] = true
}
func (s *intSet) Remove(i int) {
	if s.set[i] {
		delete(s.set, i)
	}
}

func (s intSet) Has(i int) bool {
	return s.set[i]
}

func (s intSet) Items() []int {
	result := make([]int, 0, len(s.set))
	for k, v := range s.set {
		if v {
			result = append(result, k)
		}
	}
	return result

}
func (s *intSet) Intersect(other *intSet) *intSet {
	var (
		one, two *intSet
	)

	if s.Size() > other.Size() {
		one, two = s, other
	} else {
		one, two = other, s
	}
	result := NewIntSet()
	for k := range one.set {
		if two.Has(k) {
			result.Add(k)
		}
	}
	return result
}
