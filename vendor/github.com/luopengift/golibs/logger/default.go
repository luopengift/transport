package logger

import (
	"io"
	"os"
)

var MyLogger *Logger

func SetLevel(lv uint8) {
	MyLogger.lv = lv
}

func SetTimeFormat(timeFormat string) {
	MyLogger.timeFormat = timeFormat
}

func SetPrefix(prefix string) {
	MyLogger.prefix = prefix
}

func SetColor(color bool) {
	MyLogger.color = color
}

func SetOutput(out ...io.Writer) {
	MyLogger.out = out
}
func Output(lv uint8, format string, msg ...interface{}) {
	MyLogger.Output(lv, format, msg...)
}

func Trace(format string, msg ...interface{}) {
	MyLogger.Trace(format, msg...)
}
func Debug(format string, msg ...interface{}) {
	MyLogger.Debug(format, msg...)
}
func Info(format string, msg ...interface{}) {
	MyLogger.Info(format, msg...)
}
func Warn(format string, msg ...interface{}) {
	MyLogger.Warn(format, msg...)
}
func Error(format string, msg ...interface{}) {
	MyLogger.Error(format, msg...)
}
func Fatal(format string, msg ...interface{}) {
	MyLogger.Fatal(format, msg...)
}
func Panic(format string, msg ...interface{}) {
	MyLogger.Panic(format, msg...)
}

func init() {
	MyLogger = NewLogger("2006/01/02 15:04:05.000", DEBUG, os.Stdout)
	MyLogger.SetStack(false)
}
