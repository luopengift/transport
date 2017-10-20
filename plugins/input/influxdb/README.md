### [inputs] plugin influxdb
```
{
    "influxdb": {
        "addr": "http://10.10.10.10:8086",
        "database": "telegraf",
        "username": "root",
        "password": "root",
        "query":"select host, usage_idle from cpu where cpu='cpu-total' and host = 'ip-10-10-20-114' and time > '2017-10-20T05:05:00Z'"
        #query := `select host, usage_idle from "cpu" where "cpu"='cpu-total' and "host" = 'ip-10-10-20-114' and "time" > '2017-10-20T05:05:00Z'`
        #query := `select used_percent from "mem" where "host" = 'ip-10-10-20-14' and "time" > '2017-10-20T05:05:00Z'`
        #query:= `select mean(column_one) from haproxy_33 group by time(10m) where time > now()-1d limit 500`
    }
}
```
