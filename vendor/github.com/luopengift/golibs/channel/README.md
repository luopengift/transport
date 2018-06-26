## golang 队列

### 使用例子:
```
package main

import (
    "github.com/luopengift/golibs/channel"
    "time"
    "fmt"
)

func main() {
    var max_works int64 = 100
    c := channel.NewChannel(max_works)
    go func() {
        for {
            fmt.Println(c)
            time.Sleep(500 * time.Millisecond)
        }
    }()
    for i := 0; i < 20; i++ {
            c.Run(func() error {
                fmt.Println(fmt.Sprintf("groutine no.%d start,time %v", i, time.Now().Format("15:04:05")))
                time.Sleep(2 * time.Second)
                fmt.Println(fmt.Sprintf("groutine no.%d end,time %v", i, time.Now().Format("15:04:05")))
                return nil
            })
    }
    select{}

}

```
