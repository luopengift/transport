# transport
data transportation tool, from one to another.such as,file, kafka, hdfs etc.


```
// 将配置文件解析成制定格式v的结构
type Configer interface {
    Parse(v interface{}) error
}

任何实现了Inputer接口，即可做为input组件
type Inputer interface {
    Init(config Configer) error
    Start() error
    Read([]byte) (int, error)
    Close() error
}

任何实现了Handler接口，即可做为数据处理组件
type Handler interface {
    Init(config Configer) error
    Handle(in, out byte) error
}


任何实现了Outputer接口，即可作为output组件
type Outputer interface {
    Init(config Configer) error
    Start() error
    Write([]byte) (int, error)
    Close() error
}
```
### Input组件:
- [x] [exec](https://github.com/luopengift/transport/blob/master/plugins/input/exec/README.md): 执行程序/脚本
- [x] [file(s)](https://github.com/luopengift/transport/blob/master/plugins/input/file/README.md): 文件
- [x] [http](https://github.com/luopengift/transport/blob/master/plugins/input/http/README.md): HTTP POST方法
- [x] [kafka](https://github.com/luopengift/transport/blob/master/plugins/input/kafka/README.md): kafka
- [x] [std](https://github.com/luopengift/transport/blob/master/plugins/input/std/README.md): stdin,标准输入
- [x] [random](https://github.com/luopengift/transport/blob/master/plugins/input/random/README.md): 随机生成UUID,用于测试

- [ ] elasticsearch

### Output组件:
- [x] [file](https://github.com/luopengift/transport/blob/master/plugins/output/file/README.md): 文件
- [x] [kafka](https://github.com/luopengift/transport/blob/master/plugins/output/kafka/README.md): kafka
- [x] [null](https://github.com/luopengift/transport/blob/master/plugins/output/null/README.md): 类似于/dev/null,输出到空
- [x] [std](https://github.com/luopengift/transport/blob/master/plugins/output/std/README.md): stdout,标准输出
- [x] [elasticsearch](https://github.com/luopengift/transport/blob/master/plugins/output/elasticsearch/README.md): es
- [x] [tcp](https://github.com/luopengift/transport/blob/master/plugins/output/tcp/README.md): tcp

- [ ] hdfs

### Handler组件:
- [x] null,直接连接input,output
- [x] addenter,在行尾加入换行符,例子:写文件
- [x] grok,正则格式化成json格式,说明: ^(?P<命名>子表达式)$  被捕获的组，该组被编号且被命名 (子匹配)"
- [x] [kv](https://github.com/luopengift/transport/blob/master/plugins/codec/README.md),string split 成json格式

#### Handler可以组合Inject struct,以实现向Input/Output中注入数据,[示例](https://github.com/luopengift/transport/blob/master/plugins/codec/inject.go)
```
package codec

import (
    "github.com/luopengift/transport"
    "time"
)

type DebugInjectHandler struct {
    *transport.Inject
}

func (h *DebugInjectHandler) Init(config transport.Configer) error {
    return nil
}

func (h *DebugInjectHandler) Handle(in, out []byte) (int, error) {
    time.Sleep(1 * time.Second) // make program run slow down
    h.InjectInput(in)   //将输入数据，再次inject回recv_chan，实现数据循环处理
    n := copy(out, in)
    return n, nil
}

func init() {
    transport.RegistHandler("DEBUG_InjectInput", new(DebugInjectHandler))
}
```

## TODO
1. 优化性能
2. 加入更多组件


