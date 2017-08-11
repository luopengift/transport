package plugins

import (
	_ "github.com/luopengift/transport/plugins/codec"

	_ "github.com/luopengift/transport/plugins/input/std"
	_ "github.com/luopengift/transport/plugins/input/file"
	_ "github.com/luopengift/transport/plugins/input/kafka"

	_ "github.com/luopengift/transport/plugins/output/std"
	_ "github.com/luopengift/transport/plugins/output/file"
	_ "github.com/luopengift/transport/plugins/output/kafka"
)

func init() {}
