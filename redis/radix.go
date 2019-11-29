package redis

import (
	"time"

	"github.com/mediocregopher/radix/v3"
)

type radixRedis struct {
	pool *radix.Pool
}

//NewRadixRedis 实例化
func NewRadixRedis(opts ...Option) Redis {
	o := &option{
		network: "tcp",
		addr:    "127.0.0.1:6379",
		pass:    "",
		size:    0,
		db:      0,
		timeout: time.Minute,
	}
	for _, opt := range opts {
		opt(o)
	}
	p, err := radix.NewPool(o.network, o.addr, o.size,
		radix.PoolConnFunc(
			func(network, addr string) (radix.Conn, error) {
				return radix.Dial(network, addr,
					radix.DialTimeout(o.timeout),
					radix.DialAuthPass(o.pass),
					radix.DialSelectDB(o.db),
				)
			},
		),
	)
	if err != nil {
		panic(err)
	}
	return &radixRedis{pool: p}
}

//Do 执行
func (r *radixRedis) Do(res interface{}, cmd string, args ...string) error {
	return r.pool.Do(radix.Cmd(res, cmd, args...))
}
