package types

import "strconv"

type String struct {
	s string
}

func (s *String) Int() int {
	i, _ := strconv.Atoi(s.s)
	return i
}

func (s *String) Int64() int64 {
	i, _ := strconv.ParseInt(s.s, 10, 64)
	return i
}

func (s *String) Float64() float64 {
	num, _ := strconv.ParseFloat(s.s, 64)
	return num
}

func (s *String) String() string {
	return s.s
}
