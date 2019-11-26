package cache

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	maxTime = time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
)

type localCache struct {
	prefix string
	path   string
}

func NewLocal(prefix, root string) Cache {
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("filepath.Abs失败:%s\n", err)
	}
	path := fmt.Sprintf("%s%s/runtime/cache/", dir, root)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Panicf("os.MkdirAll失败：%s\n", err)
	}
	return &localCache{
		prefix: prefix,
		path:   path,
	}
}

func (l *localCache) Has(key string) bool {
	_, err := os.Stat(l.GetPath(key))
	if err != nil {
		return false
	}
	return true
}

func (l *localCache) Get(key string) (data string, err error) {
	file, err := os.Open(l.GetPath(key))
	if err != nil {
		return "", err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	sec, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return "", err
	}
	t := time.Unix(sec, 0)
	if !t.Equal(maxTime) && time.Now().After(t) {
		l.Del(l.GetPath(key))
		return "", nil
	}

	bts, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(bts), nil
}

func (l *localCache) Set(key string, data string, ex ...time.Duration) error {
	file, err := os.OpenFile(l.GetPath(key), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	var expire time.Time
	if len(ex) > 0 {
		expire = time.Now().Add(ex[0])
	} else {
		expire = maxTime
	}
	writer.WriteString(fmt.Sprintf("%d\n", expire.Unix()))
	writer.WriteString(data)
	return writer.Flush()
}

func (l *localCache) Del(key string) error {
	return os.Remove(l.GetPath(key))
}

func (l *localCache) GetPath(key string) string {
	hs := md5.New()
	hs.Write([]byte(l.prefix + key))
	hashString := fmt.Sprintf("%x", hs.Sum(nil))
	return fmt.Sprintf("%s%s.cache", l.path, hashString)
}
