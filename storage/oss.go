package storage

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type aliOSS struct {
	Bucket *oss.Bucket
	root   string
	host   string
}

func NewAliOSS(root, host string, opts ...Option) Storage {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	client, err := oss.New(o.endpoint, o.secretId, o.secretKey)
	if err != nil {
		log.Panic("oss.New 失败", err)
	}
	bucket, err := client.Bucket(o.bucketName)
	if err != nil {
		log.Panic("storage.Bucket 失败", err)
	}
	return &aliOSS{
		Bucket: bucket,
		root:   root,
		host:   host,
	}
}

func (a *aliOSS) Url(path string) string {
	return a.host + a.Path(path)
}

func (a *aliOSS) Path(path string) string {
	return strings.TrimLeft(a.root+path, "/")
}

func (a *aliOSS) Exists(path string) (bool, error) {
	return a.Bucket.IsObjectExist(a.Path(path))
}

func (a *aliOSS) Get(path string) (io.Reader, error) {
	readCloser, err := a.Bucket.GetObject(a.Path(path))
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(readCloser); err != nil {
		return nil, err
	}
	return buf, err
}

func (a *aliOSS) Put(path string, reader io.Reader) error {
	return a.Bucket.PutObject(a.Path(path), reader)
}

func (a *aliOSS) Delete(path string) error {
	return a.Bucket.DeleteObject(a.Path(path))
}

func (a *aliOSS) Copy(path1, path2 string) error {
	_, err := a.Bucket.CopyObject(a.Path(path1), a.Path(path2))
	return err
}

func (a *aliOSS) Move(path1, path2 string) error {
	err := a.Copy(a.Path(path1), a.Path(path2))
	if err != nil {
		return err
	}
	return a.Delete(a.Path(path1))
}
