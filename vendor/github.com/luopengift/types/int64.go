package types

type Int64 struct {
	i int64
}

func (i *Int64) Int() int {
	return int(i.i)
}

func (i *Int64) Int64() int64 {
	return i.i
}

func (i *Int64) Float64() float64 {
	return float64(i.i)
}

func (i *Int64) String() string {
	return string(i.i)
}
