### [codec] plugin json
```
2017-01-02 13:58:43 INFO This is a debug Test
=>
{"level":"DEBUG","message":"This is a debug Test","time":"2017-01-02 15:58:43"}

{
    "json": {
        "regex": "^(?P<time>[0-9]{4}-[0-9]{2}-[0-9]{2}\\s+[0-9]{2}:[0-9]{2}:[0-9]{2})\\s+(?P<level>[A-Z]{3,4})\\s+(?P<message>.*)$"
    }
}
```
