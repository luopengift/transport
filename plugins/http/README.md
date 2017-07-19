# http Plugin for Transport


### http input config
| config        | example       | required|
| ------------- |:-------------:| -------:|
| type          | http          | YES |
| addr          | :9090         | YES |
### example
```
"input":{
    "type":"http",
    "addr":":9090",
}
```

### test
```
curl -XPOST -H "Content-Type:application/json; charset=utf-8" "127.0.0.1:18081/post" -d '{"data":"test"}'
```
