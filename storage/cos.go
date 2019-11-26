package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type txCOS struct {
	client *cos.Client
	root   string
	host   string
}

func NewCOS(root, host string, opts ...Option) Storage {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	u, _ := url.Parse(fmt.Sprintf(`http://%s.cos.%s.myqcloud.com`, o.bucketName, o.region))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  o.secretId,
			SecretKey: o.secretKey,
		},
	})
	return &txCOS{
		client: c,
		root:   root,
		host:   host,
	}
}

func (t *txCOS) Exists(path string) (bool, error) {
	_, err := t.client.Object.Head(context.Background(), t.Path(path), nil)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (t *txCOS) Get(path string) (io.Reader, error) {
	res, err := t.client.Object.Get(context.Background(), t.Path(path), nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(res.Body); err != nil {
		return nil, err
	}
	return buf, err
}

func (t *txCOS) Put(path string, reader io.Reader) error {
	_, err := t.client.Object.Put(context.Background(), t.Path(path), reader, nil)
	return err
}

func (t *txCOS) Delete(path string) error {
	_, err := t.client.Object.Delete(context.Background(), t.Path(path))
	return err
}

func (t *txCOS) Copy(path1, path2 string) error {
	_, _, err := t.client.Object.Copy(context.Background(), t.Path(path2), t.Path(path1), nil)
	return err
}

func (t *txCOS) Move(path1, path2 string) error {
	err := t.Copy(t.Path(path1), t.Path(path2))
	if err != nil {
		return err
	}
	return t.Delete(t.Path(path1))
}

func (t *txCOS) Path(path string) string {
	return strings.TrimLeft(t.root+path, "/")
}

func (t *txCOS) Url(path string) string {
	return t.host + t.Path(path)
}
