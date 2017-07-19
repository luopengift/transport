package plugins

import (
    _ "github.com/luopengift/transport/plugins/handler"
    _ "github.com/luopengift/transport/plugins/file"
    _ "github.com/luopengift/transport/plugins/http"
	_ "github.com/luopengift/transport/plugins/hdfs"
	_ "github.com/luopengift/transport/plugins/kafka"
	_ "github.com/luopengift/transport/plugins/elasticsearch"
	_ "github.com/luopengift/transport/plugins/std"
	//_ "github.com/luopengift/transport/plugins/tcp"
    
)

func init() {}
