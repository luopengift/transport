package hdfs

import (
    "github.com/colinmarc/hdfs"
    "github.com/luopengift/golibs/file"
    "github.com/luopengift/golibs/logger"
    "github.com/luopengift/transport"
    "path"
    "sync"
    "time"
)

type HDFSOutput struct {
    NameNode    string  `json:"nodename"`
    File        string  `json:"file"`

    client *hdfs.Client
    fd *hdfs.FileWriter
}



func (out *HDFSOutput) Init(config transport.Configer) error {
    err := config.Parse(out)
    if err != nil {
        return err
    }
    out.client, err = hdfs.New(h.NameNode)
    if err != nil {
        return nil
    }
    return nil
}

func (out *HDFSOutput) Start() error {
    return nil
}

func (out *HDFSOutput) Write(p []byte) (int, error) {
    return out.fd.Write(p)
}


func (out *HDFSOutput) Close() error {
    out.fd.Close()
    return out.Client.Close()
}
