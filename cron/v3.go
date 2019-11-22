package cron

import "github.com/robfig/cron/v3"

type v3 struct {
	c *cron.Cron
}

func NewV3() *v3 {
	c := cron.New(cron.WithSeconds())
	c.Start()
	return &v3{c: c}
}

func (r *v3) AddFunc(spec string, cmd func()) (int, error) {
	entryID, err := r.c.AddFunc(spec, cmd)
	return int(entryID), err
}

func (r *v3) AddJob(spec string, job Job) (int, error) {
	entryID, err := r.c.AddJob(spec, job)
	return int(entryID), err
}