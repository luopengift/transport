package transport

import (
	"github.com/luopengift/golibs/cron"
)

func startCronTask() error {
	cron.StartTask()
	return nil
}

// 增加定时任务
func AddCronTask(name, spec string, fun func() error) error {
	task := cron.NewTask(name, spec, fun)
	cron.AddTask(name, task)
	return nil
}

// 删除定时任务
func DelCronTask(name string) error {
	cron.DeleteTask(name)
	return nil
}

func stopCronTask() error {
	cron.StopTask()
	return nil
}
