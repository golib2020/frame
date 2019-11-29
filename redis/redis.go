package redis

type Redis interface {
	Do(res interface{}, cmd string, args ...string) error
}
