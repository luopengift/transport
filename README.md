# transport
data transportation tool, from one to another.such as,file, kafka, hdfs etc.


```
任何实现了ReadCloser接口，即可做为input组件
type ReadCloser interface {
    Read([]byte) (int, error)
    Close() error
}

任何实现了Handler接口，即可做为数据处理组件
type Handler interface {
    Handle(in, out byte) error
}


任何实现了WritCloser接口，即可作为output组件
type WriteCloser interface {
    Write([]byte) (int, error)
    Close() error
}
```
### Input组件:
- [x] File
- [x] Kafka

### Output组件:
- [x] Kafka
- [x] HDFS

### Handler组件:
- [x] Default,直接连接input,output
- [x] 在行尾加入换行符,例子:写文件

##TODO
1. 优化性能
2. 加入更多组件


