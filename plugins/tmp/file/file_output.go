package file

import (
	"github.com/luopengift/golibs/file"
	"github.com/luopengift/transport/pipeline"
	"os"
)

type FileOutput struct {
	path  string //配置路径
	cpath string //真实路径
	fd    *os.File
}

func NewFileOutput() *FileOutput {
	return new(FileOutput)
}

func (out *FileOutput) Init(cfg map[string]string) (err error) {
	out.path = cfg["path"]
	out.cpath = file.TimeRule.Handle(out.path)
	out.fd, err = os.OpenFile(out.cpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	return err
}

func (out *FileOutput) Write(p []byte) (int, error) {
	if cpath := file.TimeRule.Handle(out.path); out.cpath != cpath {
		var err error
		out.cpath = cpath
		out.fd, err = os.OpenFile(out.cpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return 0, err
		}
	}
	return out.fd.Write(p)
}

func (out *FileOutput) Start() error {
	return nil
}

func (out *FileOutput) Close() error {
	return out.fd.Close()
}

func init() {
	pipeline.RegistOutputer("file", NewFileOutput())
}
