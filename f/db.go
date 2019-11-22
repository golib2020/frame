package f

import (
	"github.com/go-xorm/xorm"
)

func DB(name ...string) *xorm.Engine {
	return Xorm(name...)
}

//Queue 队列  | 腾讯云CMQ + 云函数
func Queue() interface{} {
	return nil
}


