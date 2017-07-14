# kakfa Plugin for Transport


### kafka input config
| config        | example       | required|
| ------------- |:-------------:| -------:|
| type          | kafka         | YES |
| addrs         | localhost:9092| YES |
| topic         | test          | YES |
| offset        | -1            | YES |
### example
```
"input":{
    "type":"kafka",
    "addrs":"172.31.4.53:9092",
    "topic":"lp_test",
    "offset":"-1"
}
```
