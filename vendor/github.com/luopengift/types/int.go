package types

import "strconv"

type Int int

func (i Int) Int() int {
	return int(i)
}

func (i Int) Int64() int64 {
	return int64(i)
}

func (i Int) Float64() float64 {
	return float64(i)
}

func (i Int) String() string {
	return strconv.Itoa(i.Int())
}
