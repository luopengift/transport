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
        "batch":100
    }
}
```
