package f

import (
	"github.com/golib2020/frame/cron"
)

const (
	cronInstancesPrefix = `cron`
)

//Corn 定时任务
func Cron() cron.Cron {
	return Instance().GetOrSetFunc(cronInstancesPrefix, func() interface{} {
		return cron.NewV3()
	}).(cron.Cron)
}
