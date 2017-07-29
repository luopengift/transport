package handler

import (
	"encoding/json"
	"fmt"
	"github.com/luopengift/gohttp"
	. "github.com/luopengift/golibs/logger"
	"github.com/luopengift/pool"
	"github.com/luopengift/transport/pipeline"
	"github.com/luopengift/transport/plugins/elasticsearch"
	"os"
	"strconv"
	"strings"
)

var logger *Logger = NewLogger(DEBUG, "2006/01/02 15:04:05.000 [transport Handle]", true, os.Stdout)

var ModuleMap = map[string][]int{}

type LogDisPatch struct {
}

func (log *LogDisPatch) Handle(in, out []byte) (int, error) {
	format := ZhiziLog{}
	err := json.Unmarshal(in, &format)
	if err != nil {
		return 0, err
	}
	if _, ok := ModuleMap[format.Module]; ok {
		ModuleMap[format.Module] = append(ModuleMap[format.Module], format.Cost)
	} else {
		ModuleMap[format.Module] = []int{format.Cost}
	}
	logger.Warn("%#v", ModuleMap)
	//logger.Error("%#v",format)
	return 0, nil
}

type GetRequestId struct{}

func (req GetRequestId) Handle(in, out []byte) (int, error) {
	format := ZhiziLog{}
	err := json.Unmarshal(in, &format)
	if err != nil {
		return 0, err
	}
	n := copy(out, format.RequestId+"\n")
	return n, nil
}

type GetDataByRequestId struct{}

var factory = func() (interface{}, error) {
	client := gohttp.NewClient().Url("http://www.baidu.com")
	logger.Info("create conn:%p", client)
	return client, nil
}
var p = pool.NewPool(100, 1000, 60, factory)

func (self *GetDataByRequestId) Handle(in, out []byte) (int, error) {
	//p.LogLevel(INFO)
	request_id := string(in)
	body := `
    {
  "query": {
    "bool": {
      "must": [
        {
          "term": {
            "request_id": "%s"
          }
        }
      ]
    }
  },
    "from":0,
    "size":100
}
    `
	conn, err := p.Get()
	if err != nil {
		return 0, err

	}
	response, err := conn.(*gohttp.Client).Reset().Url("http://10.10.10.100:9200/zhizi-log*/_search?_source=module,cost").Body(fmt.Sprintf(body, request_id)).Post()
	if err != nil {
		return 0, err
	}
	err = p.Put(conn)
	if err != nil {
		logger.Error("PUT:%#v,%#V", conn, err)
	}
	println("++%s,%s", request_id, response.String())
	resp := elasticsearch.Scroll{}
	err = json.Unmarshal(response.Bytes(), &resp)
	if err != nil {
		return 0, err
	}
	req := strings.Split(request_id, "_")
	request_time, _ := strconv.ParseInt(req[1], 10, 64)
	data := map[string]interface{}{"request_id": request_id, "uid": req[0], "request_time": request_time}

	for _, doc := range resp.Hits.Hits {
		data[doc.Source["module"].(string)] = doc.Source["cost"]
	}
	bb, err := json.Marshal(data)
	if err != nil {
		logger.Warn("T3,%v", err)
		return 0, err
	}
	bb = append(bb, '\n')
	n := copy(out, bb)
	return n, nil

}

type CountGroupByModule struct{}

func (self *CountGroupByModule) Handle(in, out []byte) (int, error) {
	return 0, nil

}

func init() {
	pipeline.RegistHandler("logdispatch", new(LogDisPatch))
	pipeline.RegistHandler("getrequestid", new(GetRequestId))
	pipeline.RegistHandler("getdatabyrequestid", new(GetDataByRequestId))
	pipeline.RegistHandler("count_groupby_module", new(CountGroupByModule))
}
