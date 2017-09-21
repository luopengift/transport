# transport

[![BuildStatus](https://travis-ci.org/luopengift/transport.svg?branch=master)](https://travis-ci.org/luopengift/transport)
[![GoDoc](https://godoc.org/github.com/luopengift/transport?status.svg)](https://godoc.org/github.com/luopengift/transport)
[![GoWalker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/luopengift/transport)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

data transportation tool, from one to another.such as,file, kafka, hdfs etc.

### Framwwork
![Framework](https://github.com/luopengift/transport/blob/master/Image/png/TransportFramework.png)

### Go version >= 1.9.0

* Demo: [Example](https://github.com/luopengift/transport/blob/master/doc/EXAMPLE.md)

## 接口
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
    Version() string
}

任何实现了Adapter接口，即可做为数据处理组件
type Adapter interface {
    Init(config Configer) error
    Handle(in, out byte) error
    Version() string
}


任何实现了Outputer接口，即可作为output组件
type Outputer interface {
    Init(config Configer) error
    Start() error
    Write([]byte) (int, error)
    Close() error
    Version() string
}
```
#### Handler可以组合Inject struct,以实现向Input/Output中注入数据,[示例](https://github.com/luopengift/transport/blob/master/plugins/codec/inject.go)

#### Transport作为library使用Input/Output,[示例](https://github.com/luopengift/transport/blob/master/plugins/output/elasticsearch/README.md)

### Input组件:
- [x] [exec](https://github.com/luopengift/transport/blob/master/plugins/input/exec/README.md): 执行程序/脚本
- [x] [file(s)](https://github.com/luopengift/transport/blob/master/plugins/input/file/README.md): 文件
- [x] [hdfs](https://github.com/luopengift/transport/blob/master/plugins/input/hdfs/README.md): hdfs
- [x] [http](https://github.com/luopengift/transport/blob/master/plugins/input/http/README.md): HTTP POST方法
- [x] [kafka](https://github.com/luopengift/transport/blob/master/plugins/input/kafka/README.md): kafka
- [x] [std](https://github.com/luopengift/transport/blob/master/plugins/input/std/README.md): stdin,标准输入
- [x] [random](https://github.com/luopengift/transport/blob/master/plugins/input/random/README.md): 随机生成UUID,用于测试

- [ ] [elasticsearch]():es API scroll query

### Output组件:
- [x] [file](https://github.com/luopengift/transport/blob/master/plugins/output/file/README.md): 文件
- [x] [kafka](https://github.com/luopengift/transport/blob/master/plugins/output/kafka/README.md): kafka
- [x] [null](https://github.com/luopengift/transport/blob/master/plugins/output/null/README.md): 类似于/dev/null,输出到空
- [x] [std](https://github.com/luopengift/transport/blob/master/plugins/output/std/README.md): stdout,标准输出
- [x] [elasticsearch](https://github.com/luopengift/transport/blob/master/plugins/output/elasticsearch/README.md): es API `/_bulk`
- [x] [tcp](https://github.com/luopengift/transport/blob/master/plugins/output/tcp/README.md): tcp
- [x] [hdfs](https://github.com/luopengift/transport/blob/master/plugins/output/hdfs/README.md): hdfs
- [x] [influxdb](https://github.com/luopengift/transport/blob/master/plugins/output/influxdb/README.md): influxdb

### Handler组件:
- [x] null,直接连接input,output
- [x] addenter,在行尾加入换行符,例子:写文件
- [x] grok,正则格式化成json格式,说明: ^(?P<命名>子表达式)$  被捕获的组，该组被编号且被命名 (子匹配)"
- [x] [kv](https://github.com/luopengift/transport/blob/master/plugins/codec/README.md),string split 成json格式

### [Docker Useage]
1. install docker

2. download src code
```
git clone https://github.com/luopengift/transport.git
cd transport
```

3. build Docker image
```
docker build -t transport:0.0.3 .
```

4. run with docker
```
docker run -p 12345:12345 transport:0.0.3 -f docker-test.json
```


### [使用](https://github.com/luopengift/transport/wiki/Useage)
1. 下载
```
git clone https://github.com/luopengift/transport.git
cd transport
```

2. 编译
```
[root@iZm5egf7xb48axmu4z1t3fZ transport]# ./init.sh build cmd/main.go 
2017-09-14.15:09:14
GOPATH init Finished. GOPATH=/data/golang:/data/golang/src
build transport success.
```

3. 查看插件列表
```
[root@iZm5egf7xb48axmu4z1t3fZ transport]# ./transport -h
Usage of ./transport:
  -f string
        (config)配置文件
  -l    (list)查看插件列表和插件版本
  -r    (read)读取当前配置文件
  -v    (version)版本号
[root@iZm5egf7xb48axmu4z1t3fZ transport]# ./transport -l
[Inputs]         version
  random          0.0.1
  std             0.0.1
  exec            0.0.1
  file            0.0.1
  files           0.0.1
  hdfs            0.0.1
  http            0.0.1
  kafka           0.0.1
[Adapters]       
  kv              0.0.3
  zhizilog        0.0.1
  null            0.0.1
  addenter        0.0.1
  grok            0.0.1
  inject          0.0.1_debug
[Outputers]      
  elasticsearch   0.0.1
  file            0.0.1
  hdfs            0.0.1
  kafka           0.0.1
  null            0.0.1
  std             0.0.1
  tcp             0.0.1

```

4. 查看当前配置文件是否可以加载成功
```
[root@iZm5egf7xb48axmu4z1t3fZ transport]# ./transport -f test/kafka-file.json -r
2017-09-14 15:11:59.955 [I] <file test/kafka-file.json is END:EOF> 
config info:
[Inputs]
  kafka:
    offset: -1
    addrs: ["10.10.20.14:9092","10.10.20.15:9092","10.10.20.16:9092"]
    topics: ["zhizi-log"]
[Adapts]
  addenter:
[Outputs]
  file:
    path: /tmp/tmp.log

```
5. 运行
```
[root@iZm5egf7xb48axmu4z1t3fZ transport]# ./transport -f test/kafka-file.json  
2017-09-14 15:13:22.107 [I] <file test/kafka-file.json is END:EOF> 
2017-09-14 15:13:22.107 [I] Transport starting... 
2017-09-14 15:13:22.107 [W] Starting loading performance data, please press CTRL+C exit... 
HttpsServer Start 0.0.0.0:12345

^C2017-09-14 15:13:31.526 [W] Get signal:interrupt, Profile File is cpu.prof/mem.prof
```

6. 启动服务[加载config.json配置文件]
```
[root@iZm5egf7xb48axmu4z1t3fZ transport]# ./init.sh start
2017-09-14.15:16:45
GOPATH init Finished. GOPATH=/data/golang:/data/golang/src
transport started..., PID=22778
[root@iZm5egf7xb48axmu4z1t3fZ transport]# ps -ef |grep -v grep |grep transport
root     22778     1  0 15:18 pts/0    00:00:00 ./transport -f config.json
```

7. 查看服务状态
```
[root@iZm5egf7xb48axmu4z1t3fZ transport]# ./init.sh status
2017-09-14.15:18:24
GOPATH init Finished. GOPATH=/data/golang:/data/golang/src
root     22778     1  0 15:18 pts/0    00:00:00 ./transport -f config.json
transport now is running already, PID=22778
```
8. 查看运行日志
```
[root@iZm5egf7xb48axmu4z1t3fZ transport]# ./init.sh tail 
2017-09-14.15:22:07
GOPATH init Finished. GOPATH=/data/golang:/data/golang/src
2017-09-14 15:18:16.399 [D] [transport] [files] recv 2017-01-02 15:58:43 DEBUG This is a debug Test 
2017-09-14 15:18:16.399 [D] [transport] [files] recv 44 
2017-09-14 15:18:16.399 [D] [transport] send 44
......
```

9. 停止服务
```
[root@iZm5egf7xb48axmu4z1t3fZ transport]# ./init.sh stop
2017-09-14.15:19:09
GOPATH init Finished. GOPATH=/data/golang:/data/golang/src
transport stoped...
```



## TODO
1. 优化性能
2. 加入更多组件


