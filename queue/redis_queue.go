package queue

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type redisQueue struct {
	db     *redis.Pool
	expire time.Duration //60s
}

func NewRedisQueue(db *redis.Pool, ex time.Duration) Queue {
	return &redisQueue{
		db:     db,
		expire: ex,
	}
}

func (r *redisQueue) Size(topic string) (int, error) {
	conn := r.db.Get()
	defer conn.Close()
	topic = r.getTopic(topic)
	return redis.Int(
		conn.Do("EVAL",
			LuaScriptSize,
			3,
			topic,
			fmt.Sprintf("%s:delayed", topic),
			fmt.Sprintf("%s:reserved", topic),
		),
	)
}

func (r *redisQueue) Push(topic string, job Job) error {
	conn := r.db.Get()
	defer conn.Close()
	topic = r.getTopic(topic)
	bts, err := job.Encode()
	if err != nil {
		return err
	}
	_, err = conn.Do("RPUSH",
		topic,
		bts,
	)
	return err
}

func (r *redisQueue) Pop(topic string, job Job) error {
	conn := r.db.Get()
	defer conn.Close()
	topic = r.getTopic(topic)
	r.migrateExpiredJobs(fmt.Sprintf("%s:delayed", topic), topic)
	if r.expire != 0 {
		r.migrateExpiredJobs(fmt.Sprintf("%s:reserved", topic), topic)
	}
	bts, err := redis.Bytes(
		conn.Do("EVAL",
			LuaScriptPop,
			2,
			topic,
			fmt.Sprintf("%s:reserved", topic),
			time.Now().Add(r.expire).Unix(),
		),
	)
	if err != nil {
		return err
	}
	if bts == nil {
		return nil
	}
	if err := job.Decode(bts); err != nil {
		return err
	}
	return nil
}

func (r *redisQueue) Later(topic string, job Job, delay time.Duration) error {
	conn := r.db.Get()
	defer conn.Close()
	topic = r.getTopic(topic)
	bts, err := job.Encode()
	if err != nil {
		return err
	}
	_, err = conn.Do("ZADD",
		fmt.Sprintf("%s:delayed", topic),
		time.Now().Add(delay).Unix(),
		bts,
	)
	return err
}

//Delete 执行成功后删除
func (r *redisQueue) Delete(topic string, job Job) error {
	conn := r.db.Get()
	defer conn.Close()
	topic = r.getTopic(topic)
	bts, err := job.Encode()
	if err != nil {
		return err
	}
	_, err = conn.Do("EVAL",
		LuaScriptDelete,
		1,
		fmt.Sprintf("%s:reserved", topic),
		fmt.Sprintf("%s", bts),
	)
	return err
}

//Release 失败了重试的情况
func (r *redisQueue) Release(topic string, job Job, delay time.Duration) error {
	conn := r.db.Get()
	defer conn.Close()
	if !job.IsRetry() {
		return errors.New("job不可以重试")
	}
	topic = r.getTopic(topic)
	bts, err := job.Encode()
	if err != nil {
		return err
	}
	_, err = conn.Do("EVAL",
		LuaScriptRelease,
		2,
		fmt.Sprintf("%s:delayed", topic),
		fmt.Sprintf("%s:reserved", topic),
		bts,
		time.Now().Add(delay).Unix(),
	)
	return err
}

//migrateExpiredJobs 合并延时job到正常队列
func (r *redisQueue) migrateExpiredJobs(from, to string) error {
	conn := r.db.Get()
	defer conn.Close()
	_, err := conn.Do("EVAL", LuaScriptMigrateExpiredJobs, 2, from, to, time.Now().Unix())
	return err
}

func (r *redisQueue) getTopic(topic string) string {
	if topic == "" {
		return "default"
	}
	return topic
}
