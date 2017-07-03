# transport
data transportation tool, from one to another.such as,file, kafka, hdfs etc.


```
任何实现了ReadCloser接口，即可做为input组件
type ReadCloser interface {
    Read([]byte) (int, error)
    Close() error
}

任何实现了WritCloser接口，即可作为output组件
type WriteCloser interface {
    Write([]byte) (int, error)
    Close() error
}
```

### TODO
1. 实现数据清洗接口
2. 添加更多通用组件
