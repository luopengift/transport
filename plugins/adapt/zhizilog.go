package codec

import (
	"encoding/json"
	"fmt"
	"github.com/luopengift/golibs/crypto"
	"github.com/luopengift/transport"
	"strconv"
	"strings"
)

type ZhiziLog struct {
	RequestId string `json:"request_id"`
	Uid       string `json:"uid"`
	ReqTime   int64  `json:"request_time"`
	Module    string `json:"module"`
	Cost      int    `json:"cost"`
	Timestamp int64  `json:"timestamp"`
	Data      string `json:"data"`
	Prefix    string `json:"Prefix"`
	MD5       string `json:"md5"`
}

//Format: xxxxxxx||K1=V1&K2=V2&K3=V3...Kn=Vn
// test||JSON:{K1={sk1=sV1}}
type ZhiziLogFormat struct{}

func (h *ZhiziLogFormat) Init(config transport.Configer) error {
	return nil
}

func (d *ZhiziLogFormat) Handle(in, out []byte) (int, error) {

	if len(in) == 0 {
		return 0, fmt.Errorf("input is null\n")
	}
	loglist := strings.Split(string(in), "||") // 切割无用日志
	if len(loglist) != 2 {
		return 0, fmt.Errorf("can't split input by ||,%s\n", string(in))
	}

	logformat := ZhiziLog{
		Prefix: loglist[0],
		MD5:    crypto.MD5(loglist[1]),
	}
	value := strings.Split(loglist[1], "&&")
	if len(value) <= 3 {
		return 0, fmt.Errorf("log format is error! log is %v", string(in))
	}

	timestamp, _ := strconv.ParseInt(value[0], 10, 64)
	logformat.Timestamp = timestamp

	logformat.RequestId = value[1]
	reqid := strings.Split(logformat.RequestId, "_")
	if len(reqid) != 2 {
		return 0, fmt.Errorf("%v<request_id> format is error! log is %v", logformat.RequestId, string(in))
	}
	logformat.Uid = reqid[0]
	req_time, _ := strconv.ParseInt(reqid[1], 10, 64)
	logformat.ReqTime = req_time

	logformat.Module = value[2]

	cost, _ := strconv.Atoi(value[3])
	logformat.Cost = cost

	output, err := json.Marshal(logformat)
	if err != nil {
		return 0, fmt.Errorf("JSON Marshal error:%v", err)
	}
	n := copy(out, output)
	return n, nil

}

func (d *ZhiziLogFormat) Version() string {
	return "0.0.1"
}

func init() {
	transport.RegistHandler("zhizilog", new(ZhiziLogFormat))
}
