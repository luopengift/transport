package plugins

import (
	_ "github.com/luopengift/transport/plugins/codec"

	_ "github.com/luopengift/transport/plugins/input/exec"
	_ "github.com/luopengift/transport/plugins/input/file"
	_ "github.com/luopengift/transport/plugins/input/hdfs"
	_ "github.com/luopengift/transport/plugins/input/http"
	_ "github.com/luopengift/transport/plugins/input/kafka"
	_ "github.com/luopengift/transport/plugins/input/random"
	_ "github.com/luopengift/transport/plugins/input/std"

	_ "github.com/luopengift/transport/plugins/output/elasticsearch"
	_ "github.com/luopengift/transport/plugins/output/file"
	_ "github.com/luopengift/transport/plugins/output/hdfs"
	_ "github.com/luopengift/transport/plugins/output/influxdb"
	_ "github.com/luopengift/transport/plugins/output/kafka"
	_ "github.com/luopengift/transport/plugins/output/null"
	_ "github.com/luopengift/transport/plugins/output/std"
	_ "github.com/luopengift/transport/plugins/output/tcp"
)

func init() {}
