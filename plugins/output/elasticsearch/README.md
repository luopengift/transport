### [outputs] plugin elasticsearch
```
{
    "elasticsearch": {
        "addrs": [
            "10.10.10.100:9200"
        ],
        "index":"lp_test.%Y%M%D",
        "type":"logs",
        "timeout":5,
        "batch":500,
        "max_procs": 100
    }
}
```

#### EsOutput library
```
package main

import (
    "fmt"
    "encoding/json"
    "github.com/luopengift/transport"
    "github.com/luopengift/transport/plugins/output/elasticsearch"
)

func main() {
    cfg := `{
            "addrs": [
                "10.10.10.100:9200"
            ],
            "index":"testing.%Y%M%D",
            "type":"logs",
            "timeout":5,
            "batch":1,
            "max_procs": 1
    }`
    config := transport.PluginConfig{}
    json.Unmarshal([]byte(cfg), &config)
    es := elasticsearch.NewEsOutput()
    if err := es.Init(config); err != nil {
        fmt.Println(err)
    }
    go es.Start()
    n, err := es.Write([]byte(`{"a":"b"}`))
    fmt.Println(n, err)
    select {}
}

```
