package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/luopengift/golibs/logger"
)

//Format: xxxxxxx||K1=V1&K2=V2&K3=V3...Kn=Vn
// test||JSON:{K1={sk1=sV1}}
type ZhiziLogFormat struct{}

func (d *ZhiziLogFormat) Handle(in, out []byte) (int, error) {
	//n := copy(out, in)
	//return n, nil

	if len(in) == 0 {
		return 0, fmt.Errorf("input is null\n")
	}
	loglist := bytes.Split(in, []byte("||")) // 切割无用日志
	if len(loglist) != 2 {
		return 0, fmt.Errorf("can't split input by ||,%s\n", string(in))
	}

	kvMap := map[string]interface{}{"others": string(loglist[0])}
	log := loglist[1]
	for _, v := range bytes.Split(log, []byte("&")) {
		n := bytes.IndexAny(v, ":")
		if n == -1 {
			return 0, fmt.Errorf("can't find : to split k and v ||,%s\n", string(v))
		}
		kvMap[string(v[:n])] = parse(v[n+1:])
	}
	output, err := json.Marshal(kvMap)
	if err != nil {
		return 0, fmt.Errorf("JSON Marshal error:%v", err)
	}
	n := copy(out, output)
	return n, nil

}

func parse(p []byte) interface{} {
	logger.Info("parse:%v", string(p))
	if p[0] != '[' && p[0] != '{' && p[1] != '"' {
		return string(p)
	}
	n := len(p)
	switch {
	case p[0] == '{' && p[1] == '"' && p[n-1] == '}':
		m := map[string]interface{}{}
		err := json.Unmarshal(p, &m)
		if err != nil {
			logger.Error("Map Json Unmarshal error:%v", err)
			return string(p)
		}
		return m
	case p[0] == '[' && p[1] == '"' && p[n-1] == ']':
		m := []interface{}{}
		err := json.Unmarshal(p, &m)
		if err != nil {
			logger.Error("List Json Unmarshal error:%v", err)
			return string(p)
		}
		return m
	case p[0] == '{' && p[1] != '"' && p[n-1] == '}':
		kvMap := map[string]interface{}{}
		for _, v := range bytes.Split(p[1:n-1], []byte(",")) {
			n := bytes.IndexAny(v, "=")
			kvMap[string(v[:n])] = parse(v[n+1:])
		}
		return kvMap
	case p[0] == '[' && p[1] != '"' && p[n-1] == ']':
		vList := []interface{}{}
		logger.Debug("D1,%#s", p[1:n-1])
		for _, v := range bytes.Split(p[1:n-1], []byte(",")) {
			vList = append(vList, parse(v))
		}
		return vList
	default:
		logger.Error("parse error")
		return string(p)

	}
}
