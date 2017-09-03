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
{"Level":"DEBUG","msg":"This is a debug Test","time":"2017-01-02 13:58:44"}
```
{
    "kv": {
        "keys": ["time","Level","msg"],
        "split": "||"
    }
}
```
