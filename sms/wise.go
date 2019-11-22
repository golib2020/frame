package sms

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

type wise struct {
	apiUrl string
	user   string
	pass   string
}

func (w *wise) Send(mobile, message string) error {
	u, err := url.Parse(w.apiUrl)
	if err != nil {
		return fmt.Errorf("url.Parse error %s", err)
	}
	query := u.Query()
	query.Add("userId", w.user)
	query.Add("password", w.pass)
	query.Add("pszMobis", mobile)
	query.Add("pszMsg", message)
	query.Add("iMobiCount", "1")
	query.Add("pszSubPort", "*")
	u.RawQuery = query.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("http.Get error %s", err)
	}
	defer resp.Body.Close()
	decoder := xml.NewDecoder(resp.Body)
	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.CharData:
			if len([]byte(token)) > 15 {
				return nil
			}
		}
	}
	return fmt.Errorf("短信发送失败")
}

func NewWise(api, user, pass string) Sms {
	sms := &wise{
		apiUrl: api,
		user:   user,
		pass:   pass,
	}
	return sms
}
