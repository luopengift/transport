# Elasticsearch Plugin for Transport


### es output config
| config        | example       | required|
| ------------- |:-------------:| -------:|
| type          | elasticsearch | YES |
| protocol      | http/https    | YES |
| addrs         | localhost:9200| YES |
| _index        | test          | YES |
| _type         | test          | YES |

### example
```
"input": {
    "type":"elasticsearch",
    "protocol":"http",
    "addrs":"172.31.16.120:9200",
    "_index":"test",
    "_type":"test"
}
```
