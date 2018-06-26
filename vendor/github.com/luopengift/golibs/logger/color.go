/**
  打印字符颜色
*/
package logger

import (
	"fmt"
)

const (
	_RED     = uint8(iota + 91) //红色
	_GREEN                      //绿色
	_YELLOW                     //黄色
	_BLUE                       //蓝色
	_MAGENTA                    //洋红
	_BLUE2                      //湖蓝
)

func Color(col uint8, s interface{}) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", col, s)
}

func None(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func Red(v interface{}) string {
	return Color(_RED, v)
}

func Green(v interface{}) string {
	return Color(_GREEN, v)
}

func Yellow(v interface{}) string {
	return Color(_YELLOW, v)
}

func Blue(v interface{}) string {
	return Color(_BLUE, v)
}

func Magenta(v interface{}) string {
	return Color(_MAGENTA, v)
}

func Blue2(v interface{}) string {
	return Color(_BLUE2, v)
}

func setColor(lv uint8, v interface{}) string {
	switch lv {
	case TRACE:
		return Green(v)
	case DEBUG:
		return Blue(v)
	case INFO:
		return None(v)
	case WARNING:
		return Yellow(v)
	case ERROR:
		return Magenta(v)
	case PANIC:
		return Red(v)
	default:
		return Blue2(v)
	}
}
