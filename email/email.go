package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

type Email interface {
	SendMessage(msg Message) error
}

type email struct {
	host      string
	port      string
	user      string
	pass      string
	plainAuth smtp.Auth
	from      mail.Address
}

func NewMail(address, user, pass, name string) (Email, error) {
	addr := strings.Split(address, ":")
	if len(addr) < 2 {
		return nil, errors.New("address error")
	}
	host, port := addr[0], addr[1]
	m := &email{
		host:      host,
		port:      port,
		user:      user,
		pass:      pass,
		plainAuth: smtp.PlainAuth("", user, pass, host),
		from:      mail.Address{Name: name, Address: user},
	}
	return m, nil
}

func (m *email) SendMessage(msg Message) error {
	return m.sendMessage(msg.ToList(), msg.Bytes(m.from))
}

func (m *email) sendMessage(to []string, msg []byte) error {
	client, err := m.auth()
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Mail(m.from.Address); err != nil {
		return err
	}
	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(msg); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return nil
}

func (m *email) auth() (*smtp.Client, error) {
	c, err := dialTls(fmt.Sprintf("%s:%s", m.host, m.port))
	if err != nil {
		log.Println("Create smtp client error:", err)
		return nil, err
	}
	if ok, _ := c.Extension("AUTH"); ok {
		if err = c.Auth(m.plainAuth); err != nil {
			log.Println("Error during AUTH", err)
			return nil, err
		}
	}
	return c, nil
}

func dialTls(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}
