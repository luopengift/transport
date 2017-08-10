package plugins

import (
	_ "github.com/luopengift/transport/plugins/codec"
	_ "github.com/luopengift/transport/plugins/input"
	_ "github.com/luopengift/transport/plugins/input/file"
	_ "github.com/luopengift/transport/plugins/output"
	_ "github.com/luopengift/transport/plugins/output/file"
)

func init() {}
