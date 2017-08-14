package exec

import (
	"github.com/luopengift/golibs/exec"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport/pipeline"
)

type ExecInput struct {
	Script  string `json:"script"`
	Crontab string `json:"cron"`

    result chan []byte
}

func NewExecInput() *ExecInput {
	return new(ExecInput)
}

func (in *ExecInput) Init(config pipeline.Configer) error {
	err := config.Parse(in)
	if err != nil {
		logger.Error("parse error:%v", err)
		return err
	}
	in.result = make(chan []byte, 1)
	pipeline.AddCronTask(
		"exec",
		in.Crontab,
		func() error {
			return in.Start()
		},
	)
	return nil
}

func (in *ExecInput) Read(p []byte) (int, error) {
	n := copy(p, <-in.result)
	return n, nil
}

func (in *ExecInput) Start() error {
	result, err := exec.CmdOut("/bin/bash", "-c", in.Script)
	in.result <- result
	return err
}

func (in *ExecInput) Close() error {
	return nil
}

func init() {
	pipeline.RegistInputer("exec", NewExecInput())
}
