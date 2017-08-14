package exec

import (
	"github.com/luopengift/golibs/cron"
	"github.com/luopengift/golibs/exec"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport/pipeline"
)

type ExecInput struct {
	*ExecInputConfig
	result chan []byte
}

func NewExecInput() *ExecInput {
	return new(ExecInput)
}

type ExecInputConfig struct {
	Script  string `json:"script"`
	Crontab string `json:"cron"`
}

func (e *ExecInput) Init(config pipeline.Configer) error {
	cfg := &ExecInputConfig{}
	err := config.Parse(cfg)
	if err != nil {
		logger.Error("parse error:%v", err)
		return err
	}
	e.ExecInputConfig = cfg
	e.result = make(chan []byte, 1)
	//task := cron.NewTask("exec", e.ExecInputConfig.Crontab, func() error { return e.Start() })
	//cron.StartTask()
	//cron.AddTask("task", task)
	pipeline.AddCronTask(
		"exec",
		e.ExecInputConfig.Crontab,
		func() error {
			return e.Start()
		},
	)
	return nil
}

func (e *ExecInput) Read(p []byte) (int, error) {
	n := copy(p, <-e.result)
	return n, nil
}

func (e *ExecInput) Start() error {
	result, err := exec.CmdOut("/bin/bash", "-c", e.ExecInputConfig.Script)
	e.result <- result
	return err
}

func (e *ExecInput) Close() error {
	cron.StopTask()
	return nil
}

func init() {
	pipeline.RegistInputer("exec", NewExecInput())
}
