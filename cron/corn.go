package cron

type Job interface {
	Run()
}

type Cron interface {
	AddFunc(spec string, cmd func()) (int, error)
	AddJob(spec string, job Job) (int, error)
}
