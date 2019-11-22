package queue

import (
	"encoding/json"
	"time"
)

type Queue interface {
	Size(topic string) (int, error)
	Push(topic string, job Job) error
	Pop(topic string, job Job) error
	Later(topic string, job Job, delay time.Duration) error
	Delete(topic string, job Job) error
	Release(topic string, job Job, delay time.Duration) error
}

type Job interface {
	IsRetry() bool
	GetBody() string
	Encode() ([]byte, error)
	Decode(data []byte) error
}

type job struct {
	Attempts int    `json:"attempts"`
	Data     string `json:"data"`
}

func NewJob(data string) Job {
	return &job{
		Attempts: 0,
		Data:     data,
	}
}

func (j *job) IsRetry() bool {
	j.Attempts++
	return j.Attempts < 3
}

func (j *job) GetBody() string {
	return j.Data
}

func (j *job) Encode() ([]byte, error) {
	return json.Marshal(j)
}

func (j *job) Decode(data []byte) error {
	return json.Unmarshal(data, j)
}
