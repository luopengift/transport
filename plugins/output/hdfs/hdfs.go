package hdfs

import (
	"github.com/colinmarc/hdfs"
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	"path/filepath"
	"time"
)

const (
	VERSION = "0.0.2"
)

type HDFSOutput struct {
	NameNode string `json:"namenode"`
	File     string `json:"file"`
	Batch    int    `json:"batch"`
	RollSize int64  `json:"rollsize"` //bytes

	path   string
	file   string
	client *hdfs.Client
	fd     *hdfs.FileWriter
	buffer chan []byte
}

func NewHDFSOutput() *HDFSOutput {
	return new(HDFSOutput)
}

func (out *HDFSOutput) Init(config transport.Configer) error {
	out.Batch = 1
	out.RollSize = 10 * transport.MB
	err := config.Parse(out)
	if err != nil {
		return err
	}

	out.path = filepath.Dir(out.File)
	out.file = filepath.Base(out.File)
	out.buffer = make(chan []byte, out.Batch*2)
	out.client, err = hdfs.New(out.NameNode)
	return err
}

func (out *HDFSOutput) Size() int64 {
	stat, err := out.client.Stat(file.TimeRule.Handle(out.File))
	if err != nil {
		logger.Error("stat file error:%v", err)
		return 0
	}
	return stat.Size()
}

func (out *HDFSOutput) prepareFd() (*hdfs.FileWriter, error) {
	err := out.client.MkdirAll(file.TimeRule.Handle(out.path), 755)
	if err != nil {
		return nil, err
	}
	out.fd, err = out.client.Append(file.TimeRule.Handle(out.File))
	if err != nil {
		out.fd, err = out.client.Create(file.TimeRule.Handle(out.File))
	}
	return out.fd, err
}

func (out *HDFSOutput) Start() error {
	var err error
	for {
		if out.RollSize <= out.Size() {
			fname := file.TimeRule.Handle(out.File)
			if err = out.client.Rename(fname, fname+"."+time.Now().Format("20060102150405")); err != nil {
				logger.Error("rename file error:%v", err)
				continue
			}
		}
		out.fd, err = out.prepareFd()
		if err != nil {
			logger.Error("prepare fd error:%v", err)
			continue
		}
		for tmp := out.Batch; tmp > 0; tmp-- {
			b, ok := <-out.buffer
			if !ok {
				logger.Error("buffer chan is closed")
				return nil
			}
			if n, err := out.fd.Write(b); err != nil {
				logger.Error("write %s error,len=%v,%v", file.TimeRule.Handle(out.File), n, err)
			}
		}
		out.fd.Close()
	}
}

func (out *HDFSOutput) Write(p []byte) (int, error) {
	out.buffer <- p
	return len(p), nil
}

func (out *HDFSOutput) Close() error {
	out.fd.Close()
	return out.client.Close()
}
func (out *HDFSOutput) Version() string {
	return VERSION
}

func init() {
	transport.RegistOutputer("hdfs", NewHDFSOutput())
}
