package queue

import (
	"errors"
	"fmt"
	"time"

	"github.com/golib2020/frame/redis"
)

type redisQueue struct {
	db     redis.Redis
	expire time.Duration //60s
}

func NewRedisQueue(db redis.Redis, ex time.Duration) Queue {
	return &redisQueue{
		db:     db,
		expire: ex,
	}
}

func (r *redisQueue) Size(topic string) (int, error) {
	topic = r.getTopic(topic)
	var rvc int
	err := r.db.Do(&rvc, "EVAL",
		LuaScriptSize,
		fmt.Sprintf("%d", 3),
		topic,
		fmt.Sprintf("%s:delayed", topic),
		fmt.Sprintf("%s:reserved", topic),
	)

	return rvc, err
}

func (r *redisQueue) Push(topic string, job Job) error {
	topic = r.getTopic(topic)
	bts, err := job.Encode()
	if err != nil {
		return err
	}
	return r.db.Do(nil, "RPUSH", topic, string(bts))
}

func (r *redisQueue) Pop(topic string, job Job) error {
	topic = r.getTopic(topic)
	r.migrateExpiredJobs(fmt.Sprintf("%s:delayed", topic), topic)
	if r.expire != 0 {
		r.migrateExpiredJobs(fmt.Sprintf("%s:reserved", topic), topic)
	}
	var bts []byte
	err := r.db.Do(&bts, "EVAL",
		LuaScriptPop,
		fmt.Sprintf("%d", 2),
		topic,
		fmt.Sprintf("%s:reserved", topic),
		fmt.Sprintf("%d", time.Now().Add(r.expire).Unix()),
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
	topic = r.getTopic(topic)
	bts, err := job.Encode()
	if err != nil {
		return err
	}
	return r.db.Do(nil, "ZADD",
		fmt.Sprintf("%s:delayed", topic),
		fmt.Sprintf("%d", time.Now().Add(delay).Unix()),
		string(bts),
	)
}

//Delete 执行成功后删除
func (r *redisQueue) Delete(topic string, job Job) error {
	topic = r.getTopic(topic)
	bts, err := job.Encode()
	if err != nil {
		return err
	}
	return r.db.Do(nil, "EVAL",
		LuaScriptDelete,
		fmt.Sprintf("%d", 1),
		fmt.Sprintf("%s:reserved", topic),
		string(bts),
	)
}

//Release 失败了重试的情况
func (r *redisQueue) Release(topic string, job Job, delay time.Duration) error {
	if !job.IsRetry() {
		return errors.New("job不可以重试")
	}
	topic = r.getTopic(topic)
	bts, err := job.Encode()
	if err != nil {
		return err
	}
	return r.db.Do(nil, "EVAL",
		LuaScriptRelease,
		fmt.Sprintf("%d", 2),
		fmt.Sprintf("%s:delayed", topic),
		fmt.Sprintf("%s:reserved", topic),
		string(bts),
		fmt.Sprintf("%d", time.Now().Add(delay).Unix()),
	)
}

//migrateExpiredJobs 合并延时job到正常队列
func (r *redisQueue) migrateExpiredJobs(from, to string) error {
	return r.db.Do(nil, "EVAL", LuaScriptMigrateExpiredJobs,
		fmt.Sprintf("%d", 2),
		from, to,
		fmt.Sprintf("%d", time.Now().Unix()),
	)
}

func (r *redisQueue) getTopic(topic string) string {
	if topic == "" {
		return "default"
	}
	return topic
}
