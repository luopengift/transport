# influxdb Plugin for Transport


### influxdb output config
| config        | example       | required|
| ------------- |:-------------:| -------:|
| type          | influxdb          | YES |
| addr          | http://localhost:8086         | YES |
| database      | influxdatabase| YES |
| username      | root          | YES |
| password      | root          | YES |
### example
```
"output":{
    "type":"influxdb",
    "addr":"http://localhost:8086",
    "database":"test",
    "username":"root",
    "password":"root"
}
```

