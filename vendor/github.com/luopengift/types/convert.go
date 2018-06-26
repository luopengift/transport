package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

// byte => string
// 直接转换底层指针，两者指向的相同的内存，改一个另外一个也会变。
// 效率是string(Bytes{})的百倍以上，且转换量越大效率优势越明显。
func BytesToString(b Bytes) string {
	return *(*string)(unsafe.Pointer(&b))
}

// string => Bytes
// 直接转换底层指针，两者指向的相同的内存，改一个另外一个也会变。
// 效率是string(Bytes{})的百倍以上，且转换量越大效率优势越明显。
// 转换之后若没做其他操作直接改变里面的字符，则程序会崩溃。
// 如 b:=String2bytes("xxx"); b[1]='d'; 程序将panic。
func StringToBytes(s string) Bytes {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*Bytes)(unsafe.Pointer(&h))
}

// file => []byte
func FileToBytes(s string) (Bytes, error) {
	return ioutil.ReadFile(s)
}

// file => map[string]interface{}
func FileToMap(s string) (Map, error) {
	b, err := FileToBytes(s)
	if err != nil {
		return nil, err
	}
	return BytesToMap(b)
}

// string => bool, if fail return false
func StringToBoolean(s string) (bool, error) {
	return strconv.ParseBool(s)
}

//string => int, if fail return 0
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

//string => int64, if fail return 0
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// string => float64, if fail return 0
func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// string => map[string]interface{}
func StringToMap(s string) (m Map, err error) {
	err = json.Unmarshal(StringToBytes(s), &m)
	return
}

//[]byte => map[string]interface{}
func BytesToMap(b Bytes) (m Map, err error) {
	err = json.Unmarshal(b, &m)
	return
}

// map[string]interface{} => []byte
func MapToBytes(m Map) (Bytes, error) {
	return json.Marshal(m)
}

// interface => []byte
func ToBytes(v interface{}) (Bytes, error) {
	switch value := reflect.ValueOf(v); v.(type) {
	case string:
		return StringToBytes(value.String()), nil
	case Bytes: //[]byte
		return value.Bytes(), nil
	default:
		return json.Marshal(v)
	}
}

// interface => string
func ToString(v interface{}) (string, error) {
	switch value := reflect.ValueOf(v); v.(type) {
	case string:
		return value.String(), nil
	case Bytes: //[]byte
		return BytesToString(value.Bytes()), nil
	default:
		b, err := json.Marshal(v)
		return string(b), err
	}
}

// interface => int
func ToInt(v interface{}) (int, error) {
	switch value := reflect.ValueOf(v); v.(type) {
	case string:
		return StringToInt(value.String())
	case int, int8, int16, int32, int64:
		return int(value.Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return int(value.Uint()), nil
	case float32, float64:
		return int(value.Float()), nil
	default:
		return 0, fmt.Errorf("ToInt Unknow type:%#v", v)
	}
}

// interface => int64
func ToInt64(v interface{}) (int64, error) {
	switch value := reflect.ValueOf(v); v.(type) {
	case string:
		return StringToInt64(value.String())
	case int, int8, int16, int32, int64:
		return value.Int(), nil
	case uint, uint8, uint16, uint32, uint64:
		return int64(value.Uint()), nil
	case float32, float64:
		return int64(value.Float()), nil
	default:
		return 0, fmt.Errorf("ToInt64 unknow type:%#v", v)
	}
}

// interface => float64
func ToFloat64(v interface{}) (float64, error) {
	switch value := reflect.ValueOf(v); v.(type) {
	case string:
		return StringToFloat64(value.String())
	case int, int8, int16, int32, int64:
		return float64(value.Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return float64(value.Uint()), nil
	case float32, float64:
		return value.Float(), nil
	case bool:
		if value.Bool() {
			return float64(1), nil
		} else {
			return float64(0), nil
		}
	default:
		return 0, fmt.Errorf("ToFloat64 unknow type:%#v", v)
	}
}

// interface => bool
func ToBool(v interface{}) (bool, error) {
	switch value := reflect.ValueOf(v); v.(type) {
	case string:
		if value.String() == "" {
			return false, nil
		}
		return true, nil
	case int, int8, int16, int32, int64:
		if value.Int() == 0 {
			return false, nil
		}
		return true, nil
	case uint, uint8, uint16, uint32, uint64:
		if value.Uint() == 0 {
			return false, nil
		}
		return true, nil
	case float32, float64:
		if value.Float() == 0 {
			return false, nil
		}
		return true, nil
	case bool:
		return value.Bool(), nil
	default:
		return false, fmt.Errorf("ToBool unknow type:%#v", v)
	}
}

// interface => map[string]interface{}
func ToMap(v interface{}) (m Map, err error) {
	err = FormatJSON(v, &m)
	return
}

func ToList(v interface{}) (l List, err error) {
	err = FormatJSON(v, &l)
	return
}

type Formatter interface {
	Format(in, out interface{}) error
}

// map[string]interface{} => struct{}
// eg: Format(map[string]interface{...}, &Struct{})
func Format(in, out interface{}) error {
	var err error
	if b, err := json.Marshal(in); err == nil {
		return json.Unmarshal(b, out)
	}
	return err
}

// float64四舍五入，并取前几位
func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}
