package types

import (
	"strings"
	"time"
)

var m map[string]string = map[string]string{
	"%Y": "2006",
	"%M": "01",
	"%D": "02",
	"%h": "15",
	"%m": "04",
	"%s": "05",
}

func Now() time.Time {
	return time.Now()
}

func TimeFormat(str string) string {
	for k, v := range m {
		str = strings.Replace(str, k, time.Now().Format(v), -1)
	}
	return str
}

type Time struct {
	t time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return nil, nil
}
