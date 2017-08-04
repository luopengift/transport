package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/luopengift/transport/pipeline"
	"strconv"
	"strings"
	"sync"
	"time"
)

const T = 60

var timeout = time.Duration(T) * time.Second

var tick = time.NewTimer(timeout)

type StoreMap struct {
	*sync.Mutex
	store map[string]map[string][]int
}

func NewStoreMap() *StoreMap {
	return &StoreMap{
		Mutex: new(sync.Mutex),
		store: make(map[string]map[string][]int),
	}
}

func (s *StoreMap) Update(request_id, module string, cost int) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if _, ok := s.store[request_id]; ok {
		//存在request_id
		if _, ok := s.store[request_id][module]; ok {
			//存在module
			s.store[request_id][module] = append(s.store[request_id][module], cost)
		} else {
			//不存在module
			s.store[request_id][module] = []int{cost}
		}
	} else {
		//不存在request_id
		s.store[request_id] = map[string][]int{
			module: []int{cost},
		}
	}
}

func (s *StoreMap) Parse(request_id, uid string, request_time int64, value map[string][]int) map[string]interface{} {
	ret := map[string]interface{}{
		"request_id":   request_id,
		"uid":          uid,
		"request_time": request_time,
	}
	for key, value := range value {
		if len(value) == 1 {
			ret[key] = value[0]
		} else {
			for i, v := range value {
				ret[key+".thread_"+strconv.Itoa(i+1)] = v
			}
		}
	}
	delete(s.store, request_id)
	return ret

}

var Store = NewStoreMap()

type ZhiziLogFormat1 struct{}

func (d *ZhiziLogFormat1) Handle(in, out []byte) (int, error) {
	if len(in) == 0 {
		return 0, fmt.Errorf("input is null\n")
	}
	loglist := strings.Split(string(in), "||") // 切割无用日志
	if len(loglist) != 2 {
		return 0, fmt.Errorf("can't split input by ||,%s\n", string(in))
	}

	value := strings.Split(loglist[1], "&&")
	if len(value) <= 3 {
		return 0, fmt.Errorf("log format is error! log is %v", string(in))
	}
	request_id := value[1]
	module := value[2]
	cost, _ := strconv.Atoi(value[3])
	Store.Update(request_id, module, cost)
	select {
	case <-tick.C:
		Store.Mutex.Lock()
		defer Store.Mutex.Unlock()
		defer tick.Reset(timeout)
		var buf bytes.Buffer
		for k, v := range Store.store {
			uid := strings.Split(k, "_")[0]
			request_time, _ := strconv.ParseInt(strings.Split(k, "_")[1], 10, 64)
			if request_time+T > time.Now().Unix() {
				ret := Store.Parse(request_id, uid, request_time, v)
				output, err := json.Marshal(ret)
				if err != nil {
					return 0, fmt.Errorf("JSON Marshal error:%v", err)
				}
				buf.Write(output)
				buf.Write([]byte("\n"))
			}
		}
		if buf.Len() > cap(out) {
			logger.Error("can't write to WriteBuffer,data length is %v,but buffer cap is %v,\n%v", buf.Len(), cap(out), buf.String())
			return 0, pipeline.MaxBytesError
		}
		n := copy(out, buf.Bytes())
		return n, nil
	default:

	}

	return 0, nil
}

func init() {
	pipeline.RegistHandler("zhizilog1", new(ZhiziLogFormat1))
}
