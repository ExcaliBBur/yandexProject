package utility

type Stack struct {
	v []float64
}

func (s *Stack) Push(d float64) {
	s.v = append(s.v, d)
}

func (s *Stack) Pop() (float64, bool) {

	if len(s.v) == 0 {
		return 0, false
	}

	v := s.v[len(s.v)-1]

	s.v = s.v[:len(s.v)-1]
	return v, true
}

func (s *Stack) Size() int {
	return len(s.v)
}
