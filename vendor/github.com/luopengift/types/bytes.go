package types

import "bytes"

type Bytes = []byte

// 清理#注释
func Clean(b Bytes) Bytes {
	ret := [][]byte{}
	for _, dat := range bytes.Split(b, []byte("\n")) {
		switch offset := bytes.Index(dat, []byte("#")); offset {
		case -1: //不存在
			ret = append(ret, dat)
		case 0: //开头是注释，全部丢弃
			continue
		default: //中间有
			if dat[offset-1] != '\\' {
				dat = dat[:offset]
			}
			ret = append(ret, dat)
		}
	}
	return bytes.Join(ret, []byte("\n"))
}
