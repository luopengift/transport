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

const (
	VERSION = "0.0.1"
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

func (in *HDFSInput) Init(config transport.Configer) error {
	err := config.Parse(in)
	if err != nil {
		return err
	}

	in.buffer = make(chan []byte, 100)
	in.client, err = hdfs.New(in.NameNode)
	return err
}

func (in *HDFSInput) Start() error {
	for {
		in.realFile = file.TimeRule.Handle(in.File)
		r, err := in.client.Open(in.realFile)
		if err != nil {
			logger.Error("client open file error:%v", err)
			continue
		}
		reader := bufio.NewReader(r)
		for {
			line, err := reader.ReadBytes('\n')
			switch {
			case err != nil && err == io.EOF:
				logger.Warn("%v EOF", in.realFile)
				return in.Close()
			case err != nil && err != io.EOF:
				logger.Error("error:%v", err)
				return err
			default:
				in.buffer <- bytes.TrimRight(line, "\n")
			}
		}
	}
}

func (in *HDFSInput) Read(p []byte) (int, error) {
	n := copy(p, <-in.buffer)
	return n, nil
}

func (in *HDFSInput) Close() error {
	close(in.buffer)
	return in.client.Close()
}

func (in *HDFSInput) Version() string {
	return VERSION
}

func init() {
	transport.RegistInputer("hdfs", NewHDFSInput())
}
