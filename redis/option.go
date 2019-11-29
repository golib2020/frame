package redis

import "time"

type option struct {
	network string
	addr    string
	pass    string
	size    int
	db      int
	timeout time.Duration
}

type Option func(*option)

//WithNetwork 网络类型
func WithNetwork(n string) Option {
	return func(o *option) {
		o.network = n
	}
}

//WithAddr 连接地址
func WithAddr(addr string) Option {
	return func(o *option) {
		o.addr = addr
	}
}

//WithPass 密码
func WithPass(pass string) Option {
	return func(o *option) {
		o.pass = pass
	}
}

//WithSize 池子数
func WithSize(size int) Option {
	return func(o *option) {
		o.size = size
	}
}

//WithSelectDB 选择db
func WithSelectDB(db int) Option {
	return func(o *option) {
		o.db = db
	}
}

//WithTimeout 超时
func WithTimeout(d time.Duration) Option {
	return func(o *option) {
		o.timeout = d
	}
}
