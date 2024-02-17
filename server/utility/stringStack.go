package utility

type StringStack struct {
	v []string
}

func (s *StringStack) Push(d string) {
	s.v = append(s.v, d)
}

func (s *StringStack) Pop() (string, bool) {

	if len(s.v) == 0 {
		return "", false
	}

	v := s.v[len(s.v)-1]

	s.v = s.v[:len(s.v)-1]
	return v, true
}

func (s *StringStack) Top() string {

	if len(s.v) == 0 {
		return ""
	}

	v := s.v[len(s.v)-1]

	return v
}

func (s *StringStack) Size() int {
	return len(s.v)
}
