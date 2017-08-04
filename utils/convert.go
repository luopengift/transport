package utils

import (
    "strconv"
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
