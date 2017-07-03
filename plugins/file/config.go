package file

import (
	"encoding/json"
	"strings"
)

type Config struct {
	file     string
	comment  string
	contents []string
}

func NewConfig(file string) *Config {
	return &Config{
		file:     file,
		comment:  "#",
		contents: []string{},
	}
}

// 处理行内注释
func (self *Config) cleanLine(line string) string {
	offset := strings.Index(line, self.comment)
	switch offset {
	case -1, 0:
	default:
		if line[offset-1] != '\\' {
			sline := strings.Split(line, self.comment)
			return sline[0]
		}
	}
	return line
}

// 去除json文件中的注释
func (self *Config) Clean() string {
	f := NewTail(self.file)
	f.EndStop(true)
	f.ReadLine()
	for v := range f.NextLine() {
		if !strings.HasPrefix(*v, self.comment) {
			self.contents = append(self.contents, self.cleanLine(*v))
		}
	}
	return self.String()
}

//格式化配置文件
func (self *Config) Parse(v interface{}) error {
	self.Clean()
	return json.Unmarshal([]byte(self.String()), v)
}

func (self *Config) String() string {
	return strings.Join(self.contents, "\n")
}
