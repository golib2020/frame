package f

import (
	"fmt"
	"github.com/golib2020/frame/email"
)

const (
	emailInstancesPrefix = `email`
	emailDefaultGroup    = `default`
)

//Email 邮件服务
func Email(name ...string) email.Email {
	group := emailDefaultGroup
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	key := fmt.Sprintf("%s.%s", emailInstancesPrefix, group)
	return Instance().GetOrSetFunc(key, func() interface{} {
		return emailInit(key)
	}).(email.Email)
}

func emailInit(key string) email.Email {
	conf := Config().Sub(key)
	m, err := email.NewMail(
		conf.GetString("addr"),
		conf.GetString("user"),
		conf.GetString("pass"),
		conf.GetString("name"),
	)
	if err != nil {
		panic(err)
	}
	return m
}
