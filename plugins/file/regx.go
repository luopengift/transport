package file

import (
	"regexp"
	"time"
)

var (
	//时间通配符，用于正则表达式替换
	Map map[string]string = map[string]string{
		"%Y": "2006",
		"%M": "01",
		"%D": "02",
		"%h": "15",
		"%m": "04",
		"%s": "05",
	}
)

//eg:"test-%Y%M%D.log" ->"test-20170203.log"
//eg:"test-%Y-%M-%D.log" ->"test-2017-02-03.log"
func HandlerRule(str string) string {
	for k, v := range Map {
		re, err := regexp.Compile(k)
		if err != nil {
			continue
		}
		str = re.ReplaceAllString(str, v)
	}
	return time.Now().Format(str)
}
