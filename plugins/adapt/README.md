### [codec] plugin grok
```
2017-01-02 13:58:43 INFO This is a debug Test
=>
{"level":"DEBUG","message":"This is a debug Test","time":"2017-01-02 15:58:43"}

{
    "grok": {
        "regex": "^(?P<time>[0-9]{4}-[0-9]{2}-[0-9]{2}\\s+[0-9]{2}:[0-9]{2}:[0-9]{2})\\s+(?P<level>[A-Z]{3,4})\\s+(?P<message>.*)$"
    }
}
```

### [codec] plugin kv
2017-01-02 13:58:43||INFO||This is a debug Test
=>
{"Level":"DEBUG","msg":"This is a debug Test","time":"2017-01-02 13:58:44", "T1": "G!" }
```
{
    "kv": {
        "keys": [
            [ "time", "string" ],
            [ "Level","string" ],
            [ "msg", "string" ]
        ],
        "split": "||",
        "ignore": "-",
        "tags": {
            "T1": "G!"
        },
        "geoip":"ip => geoip",
        "ipdb": "utils/GeoLite2-City.mmdb"

    }
}
```
说明:
keys 格式为[][]string <key, type> 
type: string, int, int64, float64, json
split: 分隔符
ignore: keys中为ignore的字段，为忽略字段
tags: 额外新增的字段
geoip: 将指定字段转换成位置格式
ipdb: GeoIP数据库地址
