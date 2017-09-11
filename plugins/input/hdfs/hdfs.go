package hdfs

import (
	"bufio"
	"bytes"
	"github.com/colinmarc/hdfs"
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	"io"
)

type HDFSInput struct {
	NameNode string `json:"namenode"`
	File     string `json:"file"`

	realFile string //record real file name
	client   *hdfs.Client
	buffer   chan []byte
}

func NewHDFSInput() *HDFSInput {
	return new(HDFSInput)
}

func (out *HDFSInput) Init(config transport.Configer) error {
	err := config.Parse(out)
	if err != nil {
		return err
	}

	out.buffer = make(chan []byte, 100)
	out.client, err = hdfs.New(out.NameNode)
	return err
}

func (out *HDFSInput) Start() error {
	for {
		out.realFile = file.TimeRule.Handle(out.File)
		r, err := out.client.Open(out.realFile)
		if err != nil {
			logger.Error("client open file error:%v", err)
			continue
		}
		reader := bufio.NewReader(r)
		for {
			line, err := reader.ReadBytes('\n')
			switch {
			case err != nil && err == io.EOF:
				logger.Warn("%v EOF", out.realFile)
				return out.Close()
			case err != nil && err != io.EOF:
				logger.Error("error:%v", err)
				return err
			default:
				out.buffer <- bytes.TrimRight(line, "\n")
			}
		}
		r.Close()
	}
	return nil
}

func (out *HDFSInput) Read(p []byte) (int, error) {
	n := copy(p, <-out.buffer)
	return n, nil
}

func (out *HDFSInput) Close() error {
	close(out.buffer)
	return out.client.Close()
}

func init() {
	transport.RegistInputer("hdfs", NewHDFSInput())
}
