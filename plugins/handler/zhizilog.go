package  handler
import (
    "encoding/json"
    "github.com/luopengift/golibs/logger"
    "github.com/luopengift/transport"
)

var ModuleMap = map[string][]int{}

type LogDisPatch struct {
}

func (log *LogDisPatch) Handle(in,out []byte) (int,error) {
    format := ZhiziLog{}
    err := json.Unmarshal(in, &format) 
    if err != nil {
        return 0,err
    }
    if _,ok := ModuleMap[format.Module]; ok {
        ModuleMap[format.Module] = append(ModuleMap[format.Module],format.Cost)
    }else{
        ModuleMap[format.Module] = []int{format.Cost}
    }
    logger.Warn("%#v",ModuleMap)
    //logger.Error("%#v",format)
    return 0,nil
}

func init() {
    transport.RegistHandler("logdispatch", new(LogDisPatch))
}

