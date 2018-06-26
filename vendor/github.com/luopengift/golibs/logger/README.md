## 简单日志模块

#### 分级别(trace,debug,info,warn,error,fatal,panic),显示颜色。
#### TODO:增加其他输出方式

### 使用例子:
```
package main

import (
    "github.com/luopengift/golibs/logger"
)

func main() {
    logger.SetPrefix("module")
    logger.SetTimeFormat("2006/01/02 15:04:05.000")
    logger.Trace("%s,%s","hello","world")
    logger.Debug("%s,%s","hello","world")
    logger.Info("%s,%s","hello","world")
    logger.Warn("%s,%s","hello","world")
    logger.Error("%s,%s","hello","world")
    logger.Fatal("%s,%s","hello","world")
}
```
