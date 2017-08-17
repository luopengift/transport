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
- [x] exec: 执行程序/脚本
- [x] file(s): 文件
- [x] http: HTTP POST方法
- [x] kafka: kafka
- [x] std: stdin,标准输入

- [ ] elasticsearch
- [ ] random

### Output组件:
- [x] file: 文件
- [x] kafka: kafka
- [x] null: 类似于/dev/null,输出到空
- [x] std: stdout,标准输出

- [ ] elasticsearch
- [ ] hdfs

### Handler组件:
- [x] null,直接连接input,output
- [x] addenter,在行尾加入换行符,例子:写文件

## TODO
1. 优化性能
2. 加入更多组件


