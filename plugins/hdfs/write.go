package hdfs

import (
	"github.com/colinmarc/hdfs"
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/golibs/logger"
	"path"
	"sync"
	"time"
)

type Config struct {
	Namenode string `json:"namenode"`
	Path     string `json:"path"`
	File     string `json:"file"`
}

func NewConfig(namenode, path, file string) *Config {
	cfg := new(Config)
	cfg.Namenode = namenode
	cfg.Path = path
	cfg.File = file
	return cfg
}

type HDFS struct {
	*Config
	*hdfs.Client
	*hdfs.FileWriter
	*sync.Mutex
}

func NewHDFS(namenode, path, file string) *HDFS {
	h := new(HDFS)
	h.Config = NewConfig(namenode, path, file)
	h.Mutex = new(sync.Mutex)
	return h
}

func (h *HDFS) Init() error {
	var err error
	h.Client, err = hdfs.New(h.Namenode)
	if err != nil {
		return err
	}
	h.FileWriter, err = h.PrepareFileFd()
	return err
}

func (h *HDFS) ReConnect() error {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	if err := h.Close(); err != nil {
		return err
	}
	if err := h.Init(); err != nil {
		return err
	}
	return nil
}

func (h *HDFS) PrepareFileFd() (fd *hdfs.FileWriter, err error) {
	hdfs_path := file.TimeRule.Handle(h.Path)
	hdfs_file := file.TimeRule.Handle(h.File)
	hdfs_abs := path.Join(hdfs_path, hdfs_file)
	h.Client.MkdirAll(hdfs_path, 755)

	if _, err = h.Client.Stat(hdfs_abs); err == nil {
		// Append opens an existing file in HDFS and returns an io.WriteCloser for writing to it
		if fd, err = h.Client.Append(hdfs_abs); err == nil {
			return fd, err
		}
	} else {
		// Create opens a new file in HDFS
		if fd, err = h.Client.Create(hdfs_abs); err == nil {
			return fd, err
		}
	}
	return nil, err
}

func (h *HDFS) Write(p []byte) (int, error) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	return h.FileWriter.Write(p)
}

func (h *HDFS) Start() error {
	tc := time.Tick(60 * time.Second)
	for {
		select {
		case <-tc:
			logger.Info("reopen start")
			if err := h.ReConnect(); err != nil {
				logger.Error("reconnect hdfs error:%v", err)
			}
			logger.Info("reopen success, %#v", h)
		}
	}
	return nil
}

func (h *HDFS) Close() (err error) {
	if h.FileWriter != nil {
		if err = h.FileWriter.Close(); err != nil {
			logger.Trace("filewriter close fail:%#v", err)
			return err
		}
	}
	if h.Client != nil {
		if err = h.Client.Close(); err != nil {
			logger.Trace("client close fail:%#v", err)
			return err
		}
	}
	return nil
}
