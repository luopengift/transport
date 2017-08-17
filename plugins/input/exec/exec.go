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
    errchan chan error
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
    in.errchan = make(chan error, 1)
	pipeline.AddCronTask(
		"exec",
		in.Crontab,
		func() error {
			return in.run()
		},
	)
	return nil
}

func (in *ExecInput) Read(p []byte) (int, error) {
    select {
        case err := <-in.errchan:
            return 0,err
        case b := <-in.result:
            return copy(p, b), nil
    }
}

func (in *ExecInput) Start() error {
    return nil
}

func (in *ExecInput) run() error {
	result, err := exec.CmdOut("/bin/bash", "-c", in.Script)
    if err != nil {
        in.errchan <- err
        return err
    }
    in.result <- result
	return nil
}

func (in *ExecInput) Close() error {
	return nil
}

func init() {
	pipeline.RegistInputer("exec", NewExecInput())
}
