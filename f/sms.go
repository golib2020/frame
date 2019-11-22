package f

import (
	"fmt"
	"github.com/golib2020/frame/sms"
)

const (
	smsInstancesPrefix = `sms`
	smsDefaultGroup    = `default`
)

//SMS 短信服务
func SMS(name ...string) sms.Sms {
	group := smsDefaultGroup
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	key := fmt.Sprintf("%s.%s", smsInstancesPrefix, group)
	return Instance().GetOrSetFunc(key, func() interface{} {
		return smsInit(key)
	}).(sms.Sms)
}

func smsInit(key string) sms.Sms {
	conf := Config().Sub(key)
	var s sms.Sms
	switch conf.GetString("driver") {
	case "wise":
		s = sms.NewWise(
			conf.GetString("api"),
			conf.GetString("user"),
			conf.GetString("pass"),
		)
	}
	return s
}
