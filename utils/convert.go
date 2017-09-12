package utils

import (
	"strconv"
    "encoding/json"
)

func Int(v interface{}) int {
	switch value := v.(type) {
	case string:
		ret, _ := strconv.Atoi(value)
		return ret
	case int64:
		return int(value)
	default:
		return 0
	}
}


// in format out
// eg: Format(map[string]interface{...}, &Struct{})
func Format(in, out interface{}) error {
    var err error
    if b, err := json.Marshal(in); err == nil {
        err = json.Unmarshal(b, out)
    }
    return err
}
