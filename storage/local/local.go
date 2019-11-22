package local

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golib2020/frame/storage"
)

type local struct {
	root string
	host string
}

func NewLocal(root, host string) storage.Storage {
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("filepath.Abs失败:%s\n", err)
	}
	return &local{
		root: fmt.Sprintf("%s%sapp/", dir, root),
		host: host,
	}
}

func (o *local) Url(path string) string {
	return o.host + o.Path(path)
}

func (o *local) Path(path string) string {
	return strings.TrimLeft(o.root+path, "/")
}

func (o *local) Exists(path string) (bool, error) {
	_, err := os.Stat(o.getPath(path))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (o *local) Get(path string) (io.Reader, error) {
	file, err := os.Open(o.getPath(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, err
	}
	return buf, err
}

func (o *local) Put(path string, reader io.Reader) error {
	file, err := os.OpenFile(o.getPath(path), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	io.Copy(file, reader)
	return err
}

func (o *local) Append(path string, reader io.Reader) error {
	file, err := os.OpenFile(o.getPath(path), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	io.Copy(file, reader)
	return err
}

func (o *local) Delete(path string) error {
	return os.Remove(o.getPath(path))
}

func (o *local) Copy(path1, path2 string) error {
	file1, err := os.Open(o.getPath(path1))
	if err != nil {
		return err
	}
	defer file1.Close()
	file2, err := os.OpenFile(o.getPath(path2), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file2.Close()
	io.Copy(file2, file1)
	return err
}

func (o *local) Move(path1, path2 string) error {
	err := o.Copy(o.getPath(path1), o.getPath(path2))
	if err != nil {
		return err
	}
	return o.Delete(o.getPath(path1))
}

func (o *local) getPath(path string) string {
	fileName := o.root + path
	fileDir := filepath.Dir(fileName)
	_, err := os.Stat(fileDir)
	if err != nil {
		err := os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			log.Panicf("创建目录失败:%s\n", err)
		}
	}
	return fileName
}
