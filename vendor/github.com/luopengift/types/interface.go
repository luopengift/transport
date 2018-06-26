package types

//import "reflect"

type Interface struct {
	v interface{}
}

func (v *Interface) Int() int {
	i, _ := ToInt(v.v)
	return i
}

func (v *Interface) Float64() float64 {
	i, _ := ToFloat64(v.v)
	return i
}

func (v *Interface) String() string {
	i, _ := ToString(v.v)
	return i
}

func (v *Interface) List() []interface{} {
	return v.v.([]interface{})
}

func (v *Interface) Map() map[string]interface{} {
	return v.v.(map[string]interface{})
}
