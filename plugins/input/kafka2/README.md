### [inputs] plugin kafka
```
{
    "kafka": {
        "addrs": [
            "192.168.0.1:9092",
            "192.168.0.2:9092",
            "192.168.0.3:9092"
        ],
        "topics": [
            "test"
        ],
        "group":"transport",
        "offset": -1    # -1:newest,-2:oldest
    }
}

| 参数 | 类型 | 要求 | 说明 |
| -----| ---- | ---- |------|
|addrs |list  | Y    |若group为空，addrs为kafka brokers的地址；否则为zookeeper的地址 |
|topics|list  | Y    |kafka topic list| 
|group |string| N    |kafka group|
|offset|int   | Y    |-1: newest, -2: oldest|
