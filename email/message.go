package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"net/mail"
	"path/filepath"
	"strings"
	"time"
)

type Message interface {
	To(to ...string) *message
	Cc(cc ...string) *message
	Bcc(bcc ...string) *message
	ReplyTo(msg string) *message
	ToList() []string
	Bytes(from mail.Address) []byte
	AttachBuffer(filename string, buf []byte, inline bool) error
	Attach(file string) error
	Inline(file string) error
}

type attachment struct {
	Filename string
	Data     []byte
	Inline   bool
}

type message struct {
	to              []string
	cc              []string
	bcc             []string
	replyTo         string
	subject         string
	body            string
	bodyContentType string
	attachments     map[string]*attachment
}

func NewMessage(subject string, body string) Message {
	return newMessage(subject, body, "text/plain")
}

func NewHTMLMessage(subject string, body string) Message {
	return newMessage(subject, body, "text/html")
}

func newMessage(subject string, body string, bodyContentType string) *message {
	m := &message{subject: subject, body: body, bodyContentType: bodyContentType}
	m.attachments = make(map[string]*attachment)
	return m
}

func (m *message) To(to ...string) *message {
	m.to = to
	return m
}

func (m *message) Cc(cc ...string) *message {
	m.cc = cc
	return m
}

func (m *message) Bcc(bcc ...string) *message {
	m.bcc = bcc
	return m
}

func (m *message) ReplyTo(msg string) *message {
	m.replyTo = msg
	return m
}

func (m *message) AttachBuffer(filename string, buf []byte, inline bool) error {
	m.attachments[filename] = &attachment{
		Filename: filename,
		Data:     buf,
		Inline:   inline,
	}
	return nil
}

func (m *message) Attach(file string) error {
	return m.attach(file, false)
}

func (m *message) Inline(file string) error {
	return m.attach(file, true)
}

func (m *message) attach(file string, inline bool) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	_, filename := filepath.Split(file)

	m.attachments[filename] = &attachment{
		Filename: filename,
		Data:     data,
		Inline:   inline,
	}

	return nil
}

func (m *message) ToList() []string {
	tolist := m.to
	for _, cc := range m.cc {
		tolist = append(tolist, cc)
	}
	for _, bcc := range m.bcc {
		tolist = append(tolist, bcc)
	}
	return tolist
}

// Bytes returns the email data
func (m *message) Bytes(from mail.Address) []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("From: " + from.String() + "\r\n")

	t := time.Now()
	buf.WriteString("Date: " + t.Format(time.RFC1123Z) + "\r\n")

	buf.WriteString("to: " + strings.Join(m.to, ",") + "\r\n")
	if len(m.cc) > 0 {
		buf.WriteString("cc: " + strings.Join(m.cc, ",") + "\r\n")
	}

	//fix  Encode
	var coder = base64.StdEncoding
	var subject = "=?UTF-8?B?" + coder.EncodeToString([]byte(m.subject)) + "?="
	buf.WriteString("subject: " + subject + "\r\n")

	if len(m.replyTo) > 0 {
		buf.WriteString("Reply-to: " + m.replyTo + "\r\n")
	}

	buf.WriteString("MIME-Version: 1.0\r\n")

	boundary := "f46d043c813270fc6b04c2d223da"

	if len(m.attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString("\r\n--" + boundary + "\r\n")
	}

	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n\r\n", m.bodyContentType))
	buf.WriteString(m.body)
	buf.WriteString("\r\n")

	if len(m.attachments) > 0 {
		for _, attachment := range m.attachments {
			buf.WriteString("\r\n\r\n--" + boundary + "\r\n")

			if attachment.Inline {
				buf.WriteString("Content-Type: message/rfc822\r\n")
				buf.WriteString("Content-Disposition: inline; filename=\"" + attachment.Filename + "\"\r\n\r\n")

				buf.Write(attachment.Data)
			} else {
				ext := filepath.Ext(attachment.Filename)
				mistype := mime.TypeByExtension(ext)
				if mistype != "" {
					s := fmt.Sprintf("Content-Type: %s\r\n", mistype)
					buf.WriteString(s)
				} else {
					buf.WriteString("Content-Type: application/octet-stream\r\n")
				}
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")

				buf.WriteString("Content-Disposition: attachment; filename=\"=?UTF-8?B?")
				buf.WriteString(coder.EncodeToString([]byte(attachment.Filename)))
				buf.WriteString("?=\"\r\n\r\n")

				b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
				base64.StdEncoding.Encode(b, attachment.Data)

				// write base64 content in lines of up to 76 chars
				for i, l := 0, len(b); i < l; i++ {
					buf.WriteByte(b[i])
					if (i+1)%76 == 0 {
						buf.WriteString("\r\n")
					}
				}
			}

			buf.WriteString("\r\n--" + boundary)
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}
