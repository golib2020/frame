package sms

type Sms interface {
	Send(mobile, message string) error
}
